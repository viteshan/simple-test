package pbft

import (
	"bytes"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/vitelabs/go-vite/log15"
	"github.com/viteshan/go/support/errors"
)

// just from client
type ReqMsg struct {
	operation string
	timestamp int64
	owner     uint32
}

// just from master
type PPMsg struct {
	requests  *ReqMsg
	viewId    int64
	seq       int64
	hash      []byte // requests
	fromIdx   uint32
	signature []byte
}

// 2f + 1
type BFTMsg struct {
	tp        int // tp:0  prepare-message   tp:1  commit-message
	viewId    int64
	seq       int64
	hash      []byte
	fromIdx   uint32
	signature []byte
}

// f + 1
type ReplyMsg struct {
	viewId    int64
	timestamp int64
	fromIdx   uint32
	hash      []byte
	signature []byte
}

type State struct {
	viewId int64
	seq    int64

	// data
	logs []*ReqMsg
}

// 0. waiting request
// 1. waiting 2f+1 prepare
// 2. waiting 2f+1 commit
type CurState struct {
	cur          int32
	waitingState *waitingStateImpl
}

func (cur *CurState) switchCur(from int32, to int32) bool {
	return atomic.CompareAndSwapInt32(&cur.cur, from, to)
}

type waitingStateImpl struct {
	okCnt     int
	ppMsg     *PPMsg
	done      map[uint32]struct{}
	timeoutCh chan struct{}
}

func (s *waitingStateImpl) Done(idx uint32) error {
	s.done[idx] = struct{}{}
	return nil
}

func (s *waitingStateImpl) Destroy() error {
	close(s.timeoutCh)
	return nil
}

func (s waitingStateImpl) Ok() bool {
	return len(s.done) >= s.okCnt
}

type WaitingState interface {
	Ok() bool
	Done(idx uint32) error
	Destroy() error
}

type Node struct {
	clusterSize int32
	f           int

	idx   uint32
	state *State
	ch    chan interface{}

	peers   map[uint32]*Node
	clients map[uint32]*Cli

	cur *CurState
}

func (n *Node) GetCh() (chan<- interface{}, error) {
	return n.ch, nil
}

func (n *Node) loopRead() error {
	for v := range n.ch {
		err := n.onReceive(v)
		if err != nil {
			log15.Error(err.Error())
			return err
		}
	}
	return nil
}

func (n *Node) onReceive(msg interface{}) (err error) {

	switch m := msg.(type) {
	case *ReqMsg:
		fmt.Printf("[%d]receive request msg:%s\n", n.idx, msg)
		err = n.onRequestMsg(m)
	case *PPMsg:
		fmt.Printf("[%d]receive pp msg:%s\n", n.idx, msg)
		err = n.onPrePrepareMsg(m)
	case *BFTMsg:
		if m.tp == 0 {
			fmt.Printf("[%d]receive prepare msg:%s\n", n.idx, msg)
			err = n.onPrepareMsg(m)
		} else {
			fmt.Printf("[%d]receive commit msg:%s\n", n.idx, msg)
			err = n.onCommitMsg(m)
		}
	}
	return
}
func (n *Node) onRequestMsg(m *ReqMsg) error {
	if !n.isPrimary() {
		return nil
	}

	if n.cur.cur != 0 {
		return nil
	}
	pp := &PPMsg{
		requests:  m,
		viewId:    n.state.viewId,
		seq:       n.state.seq + 1,
		hash:      []byte(m.operation),
		fromIdx:   n.idx,
		signature: nil,
	}
	err := n.switchToPrepare(pp)
	if err != nil {
		log15.Error(err.Error())
		return nil
	}
	n.broadcast(pp)
	return nil
}

func (n *Node) prepareTimeout() error {
	flag := n.cur.switchCur(1, 0)
	if flag {
		n.cur.waitingState = nil
	}
	return nil
}

func (n *Node) commitTimeout() error {
	flag := n.cur.switchCur(2, 0)
	if flag {
		n.cur.waitingState = nil
	}
	return nil
}

func (n *Node) switchToPrepare(msg *PPMsg) error {
	flag := n.cur.switchCur(0, 1)
	if !flag {
		return errors.Errorf("switch fail.")
	}

	n.cur.waitingState = &waitingStateImpl{
		okCnt:     2*n.f + 1,
		ppMsg:     msg,
		done:      make(map[uint32]struct{}),
		timeoutCh: make(chan struct{}),
	}
	n.cur.waitingState.Done(n.idx)
	go func(closed chan struct{}) {
		select {
		case <-time.After(10 * time.Second):
			n.prepareTimeout()
		case <-closed:
		}
	}(n.cur.waitingState.timeoutCh)
	return nil
}

func (n *Node) apply(msg *PPMsg) error {
	flag := n.cur.switchCur(2, 0)
	if !flag {
		return errors.Errorf("switch fail.")
	}

	n.cur.waitingState.Destroy()
	n.cur.waitingState = nil
	n.state.seq = msg.seq
	n.state.logs = append(n.state.logs, msg.requests)
	return nil
}

func (n *Node) switchToCommit(msg *PPMsg) error {
	flag := n.cur.switchCur(1, 2)
	if !flag {
		return errors.Errorf("switch fail.")
	}
	old := n.cur.waitingState
	old.Destroy()

	n.cur.waitingState = &waitingStateImpl{
		okCnt:     2*n.f + 1,
		ppMsg:     msg,
		done:      make(map[uint32]struct{}),
		timeoutCh: make(chan struct{}),
	}
	n.cur.waitingState.Done(n.idx)
	go func(closed chan struct{}) {
		select {
		case <-time.After(10 * time.Second):
			n.commitTimeout()
		case <-closed:
		}
	}(n.cur.waitingState.timeoutCh)
	return nil
}

func (n *Node) broadcast(msg interface{}) error {
	for _, v := range n.peers {
		ch, _ := v.GetCh()
		ch <- msg
	}
	return nil
}
func (n *Node) broadcastToClients(msg interface{}) error {
	for _, v := range n.clients {
		ch, _ := v.GetCh()
		ch <- msg
	}
	return nil
}

func (n *Node) isPrimary() bool {
	return n.state.viewId%int64(n.clusterSize) == int64(n.idx)
}

func (n *Node) onPrePrepareMsg(m *PPMsg) error {
	if n.isPrimary() {
		return nil
	}
	if n.cur.cur != 0 {
		return nil
	}
	if m.viewId != n.state.viewId {
		return nil
	}
	// todo verify sig
	err := n.switchToPrepare(m)
	if err != nil {
		log15.Error(err.Error())
		return nil
	}
	msg := &BFTMsg{
		tp:        0,
		viewId:    m.viewId,
		seq:       m.seq,
		hash:      m.hash,
		fromIdx:   n.idx,
		signature: nil,
	}
	n.broadcast(msg)
	return nil
}

func (n *Node) onPrepareMsg(m *BFTMsg) error {
	if n.cur.cur != 1 {
		return nil
	}

	pp := n.cur.waitingState.ppMsg
	if bytes.Equal(pp.hash, m.hash) &&
		pp.seq == m.seq &&
		pp.viewId == m.viewId {
		n.cur.waitingState.Done(m.fromIdx)
		if n.cur.waitingState.Ok() {
			err := n.switchToCommit(pp)
			if err != nil {
				log15.Error(err.Error())
				return nil
			}
			msg := &BFTMsg{
				tp:        1,
				viewId:    m.viewId,
				seq:       m.seq,
				hash:      m.hash,
				fromIdx:   n.idx,
				signature: nil,
			}
			n.broadcast(msg)
		}
		return nil
	} else {
		return nil
	}
}

func (n *Node) onCommitMsg(m *BFTMsg) error {
	if n.cur.cur != 2 {
		return nil
	}

	pp := n.cur.waitingState.ppMsg
	if bytes.Equal(pp.hash, m.hash) &&
		pp.seq == m.seq &&
		pp.viewId == m.viewId {
		n.cur.waitingState.Done(m.fromIdx)
		if n.cur.waitingState.Ok() {
			err := n.apply(pp)
			if err != nil {
				return nil
			}
			msg := &ReplyMsg{
				viewId:    n.state.viewId,
				timestamp: 0,
				fromIdx:   n.idx,
				hash:      m.hash,
				signature: nil,
			}
			n.broadcastToClients(msg)
		}
		return nil
	} else {
		return nil
	}
}

type Cli struct {
	clusterSize int32
	f           int
	idx         uint32
	viewId      int64
	ch          chan interface{}
	peers       map[uint32]*Node

	waitingReply *waitingStateImpl
}

func (c Cli) calcPrimary() uint32 {
	return uint32(c.viewId % int64(c.clusterSize))

}

func (c *Cli) SendRequest(op string) error {
	if c.waitingReply != nil {
		return errors.Errorf("waiting reply")
	}
	msg := &ReqMsg{
		operation: op,
		timestamp: time.Now().Unix(),
		owner:     c.idx,
	}

	c.waitingReply = &waitingStateImpl{
		okCnt: c.f + 1,
		ppMsg: &PPMsg{
			hash: []byte(op),
		},
		done:      make(map[uint32]struct{}),
		timeoutCh: make(chan struct{}),
	}

	waitingCh := make(chan error)
	go func(closed chan struct{}) {
		select {
		case <-time.After(10 * time.Second):
			c.requestTimeout()
			waitingCh <- errors.New("request timeout")
		case <-closed:
			close(waitingCh)
		}
	}(c.waitingReply.timeoutCh)
	primary := c.calcPrimary()
	node := c.peers[primary]

	ch, _ := node.GetCh()
	ch <- msg
	err := <-waitingCh
	c.waitingReply = nil
	return err
}

func (c *Cli) GetCh() (chan<- interface{}, error) {
	return c.ch, nil
}

func (c *Cli) requestTimeout() error {
	log15.Error("request timeout")
	return nil
}

func (c *Cli) loopRead() error {
	for v := range c.ch {
		err := c.onReceive(v)
		if err != nil {
			log15.Error(err.Error())
			return err
		}
	}
	return nil
}

func (c *Cli) onReceive(msg interface{}) (err error) {
	switch m := msg.(type) {
	case *ReplyMsg:
		fmt.Printf("[%d]receive reply msg:%s\n", c.idx, msg)
		err = c.onReplyMsg(m)
	}
	return
}

func (c *Cli) onReplyMsg(m *ReplyMsg) error {
	impl := c.waitingReply
	if impl == nil {
		return nil
	}
	if bytes.Equal(impl.ppMsg.hash, m.hash) {
		impl.Done(m.fromIdx)
		if impl.Ok() {
			impl.Destroy()
		}
	}
	return nil
}
