package pbft

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/golang-collections/collections/queue"

	"github.com/vitelabs/go-vite/log15"
	"github.com/viteshan/go/support/errors"
)

// just from client
type ReqMsg struct {
	operation string
	timestamp int64
	owner     uint32
	hash      string
}

// node to node
type NodeMsg interface {
	getViewId() int64
	getSeq() int64
	getFromIdx() uint32
}

type nodeMsg struct {
	NodeMsg
	viewId  int64
	seq     int64
	fromIdx uint32
}

func (m nodeMsg) getViewId() int64 {
	return m.viewId
}

func (m nodeMsg) getSeq() int64 {
	return m.seq
}

func (m nodeMsg) getFromIdx() uint32 {
	return m.fromIdx
}

type HeartBeatMsg struct {
	nodeMsg
}

// just from master
type PPMsg struct {
	nodeMsg
	requests  *ReqMsg
	hash      string // requests
	signature []byte
}

// 2f + 1
type BFTMsg struct {
	nodeMsg
	tp        int // tp:0  prepare-message   tp:1  commit-message
	hash      string
	signature []byte
}

// f + 1
type ReplyMsg struct {
	viewId    int64
	timestamp int64
	fromIdx   uint32
	hash      string
	signature []byte
}

type SyncMsg struct {
	nodeMsg
	tp   int // tp:0  request sync  tp:1 response sync
	hash string
	logs []*ReqMsg
}

type State struct {
	viewId int64
	seq    int64

	// data
	logs []*ReqMsg
}

func (s State) calcHash() string {
	var result []byte
	for _, v := range s.logs {
		result = append(result, v.hash...)
	}
	return calculateHash(result)
}

func (s State) getLogsBySeq(seq int64) (result []*ReqMsg, hash string) {
	for k, v := range s.logs {
		if int64(k) < seq {
			result = append(result, v)
		}

	}
	var byt []byte
	for _, v := range result {
		byt = append(byt, v.hash...)
	}
	hash = calculateHash(byt)
	return
}

func calculateHash(bs []byte) string {
	h := sha256.New()
	h.Write(bs)
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// 0. waiting request
// 1. waiting 2f+1 prepare
// 2. waiting 2f+1 commit
type CurState struct {
	cur          int32
	ppMsg        *PPMsg
	waitingState *stateCounterImpl
}

func (cur *CurState) reset() {
	cur.ppMsg = nil
	cur.waitingState = nil
}

func (cur *CurState) switchCur(from int32, to int32) bool {
	return atomic.CompareAndSwapInt32(&cur.cur, from, to)
}

type stateCounterImpl struct {
	okCnt     int
	done      map[string]map[uint32]struct{}
	timeoutCh chan struct{}
}

func (s *stateCounterImpl) Done(key string, idx uint32) error {
	cnt, ok := s.done[key]
	if !ok {
		cnt = make(map[uint32]struct{})
		s.done[key] = cnt
	}
	cnt[idx] = struct{}{}
	return nil
}

func (s *stateCounterImpl) Destroy() error {
	close(s.timeoutCh)
	return nil
}

func (s stateCounterImpl) Ok(key string) bool {
	cnt := s.done[key]
	return len(cnt) >= s.okCnt
}

type StateCounter interface {
	Ok(key string) bool
	Done(key string, idx uint32) error
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

	net Net

	mBuf *msgBuffer

	down bool

	cur *CurState

	// key: sync:{seq}-{hash}  ->  counter sync response
	// key: seq:{seq}   -> counter sync request trigger
	syncCounter *stateCounterImpl
	timeout     Blacklist
}

func (n *Node) GetCh() (chan<- interface{}, error) {
	return n.ch, nil
}

func (n *Node) loopHeartBeat() error {
	for {
		time.Sleep(2 * time.Second)
		s := *n.state
		msg := &HeartBeatMsg{
			nodeMsg: nodeMsg{
				viewId:  s.viewId,
				seq:     s.seq,
				fromIdx: n.idx,
			},
		}
		n.broadcast(msg)
	}
}

func (n *Node) loopRead() error {
	for v := range n.ch {
		if n.down {
			continue
		}
		err := n.onReceive(v)
		if err != nil {
			log15.Error(err.Error())
			return err
		}
	}
	return nil
}

func (n *Node) loopBuf() error {
	for {
		if n.cur.cur == 0 {
			if n.isPrimary() {
				n.waitingRequestForPrimary()
			} else {
				n.waitingRequest()
			}
		} else if n.cur.cur == 1 {
			n.waitingPrepare()
		} else if n.cur.cur == 2 {
			n.waitingCommit()
		} else {
			panic("unknown status")
		}
		n.waitingSync()
		time.Sleep(20 * time.Millisecond)
	}
}

func (n *Node) waitingRequestForPrimary() (err error) {
	for n.mBuf.reqBuf.Len() > 0 {
		m := n.mBuf.reqBuf.Dequeue().(*ReqMsg)
		n.onRequestMsg(m)
	}
	return nil
}

func (n *Node) waitingRequest() (err error) {
	for n.mBuf.ppBuf.Len() > 0 {
		m := n.mBuf.ppBuf.Dequeue().(*PPMsg)
		n.onPrePrepareMsg(m)
	}
	return nil
}

func (n *Node) waitingPrepare() (err error) {
	for n.mBuf.prepareBuf.Len() > 0 {
		m := n.mBuf.prepareBuf.Dequeue().(*BFTMsg)
		n.onPrepareMsg(m)
	}
	return nil
}

func (n *Node) waitingCommit() (err error) {
	for n.mBuf.commitBuf.Len() > 0 {
		m := n.mBuf.commitBuf.Dequeue().(*BFTMsg)
		n.onCommitMsg(m)
	}
	return nil
}

func (n *Node) intercept(msg interface{}) error {
	if nnMsg, ok := msg.(NodeMsg); ok {
		if _, ok := msg.(*SyncMsg); ok {
			return nil
		}

		if _, ok := msg.(*HeartBeatMsg); ok {
			fmt.Printf("[%d]receive heartbeat msg:%v\n", n.idx, msg)
			if nnMsg.getSeq() > n.state.seq && n.cur.cur == 0 {
				key := fmt.Sprintf("seq:%d", nnMsg.getSeq())
				if n.timeout.Exists(key) {
					return nil
				}
				n.syncCounter.Done(key, nnMsg.getFromIdx())
				if n.syncCounter.Ok(key) {
					n.timeout.AddAddTimeout(key, time.Second*10)
					n.syncTo(nnMsg.getSeq())
				}
			}
		} else {
			if nnMsg.getSeq() > n.state.seq+1 && n.cur.cur == 0 {
				key := fmt.Sprintf("seq:%d", nnMsg.getSeq()-1)
				if n.timeout.Exists(key) {
					return nil
				}
				n.syncCounter.Done(key, nnMsg.getFromIdx())
				if n.syncCounter.Ok(key) {
					n.timeout.AddAddTimeout(key, time.Second*10)
					n.syncTo(nnMsg.getSeq() - 1)
				}
			}
		}
	}
	return nil
}

func (n *Node) onReceive(msg interface{}) (err error) {

	// todo
	n.intercept(msg)

	switch m := msg.(type) {
	case *ReqMsg:
		fmt.Printf("[%d]receive request msg:%s\n", n.idx, msg)
		n.mBuf.reqBuf.Enqueue(m)
		//fmt.Printf("[%d]request msg loss, %v", n.idx, m)
	case *PPMsg:
		fmt.Printf("[%d]receive pp msg:%s\n", n.idx, msg)
		n.mBuf.ppBuf.Enqueue(m)
		//fmt.Printf("[%d]pp msg loss, %v", n.idx, m)
	case *BFTMsg:
		if m.tp == 0 {
			fmt.Printf("[%d]receive prepare msg:%s\n", n.idx, msg)
			n.mBuf.prepareBuf.Enqueue(m)
			//fmt.Printf("[%d]prepare msg loss, %v", n.idx, m)
		} else {
			fmt.Printf("[%d]receive commit msg:%s\n", n.idx, msg)
			n.mBuf.commitBuf.Enqueue(m)
			//fmt.Printf("[%d]commit msg loss, %v", n.idx, m)
		}
	case *SyncMsg:
		if m.tp == 0 {
			fmt.Printf("[%d]receive sync request msg:%s\n", n.idx, msg)
			n.onRequestSyncMsg(m)
		} else {
			fmt.Printf("[%d]receive sync response msg:%s\n", n.idx, msg)
			n.mBuf.syncBuf.Enqueue(m)
		}
	case *HeartBeatMsg:
	}
	return
}
func (n *Node) onRequestMsg(m *ReqMsg) error {
	fmt.Printf("[%d]on request msg:%v\n", n.idx, m)
	if !n.isPrimary() {
		return nil
	}

	if n.cur.cur != 0 {
		return nil
	}
	pp := &PPMsg{
		nodeMsg: nodeMsg{
			viewId:  n.state.viewId,
			seq:     n.state.seq + 1,
			fromIdx: n.idx,
		},
		requests:  m,
		hash:      calculateHash([]byte(m.operation)),
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
		n.cur.reset()
	}
	return nil
}

func (n *Node) commitTimeout() error {
	flag := n.cur.switchCur(2, 0)
	if flag {
		n.cur.reset()
	}
	return nil
}

func (n *Node) syncTimeout() error {
	flag := n.cur.switchCur(3, 0)
	if flag {
		n.cur.reset()
	}
	return nil
}

func (n *Node) switchToPrepare(msg *PPMsg) error {
	flag := n.cur.switchCur(0, 1)
	if !flag {
		return errors.Errorf("switch fail.")
	}

	n.cur.waitingState = &stateCounterImpl{
		okCnt: 2*n.f + 1,

		done:      make(map[string]map[uint32]struct{}),
		timeoutCh: make(chan struct{}),
	}
	n.cur.ppMsg = msg
	n.cur.waitingState.Done("", n.idx)
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
	n.cur.reset()
	n.state.seq = msg.seq
	n.state.logs = append(n.state.logs, msg.requests)
	return nil
}

func (n *Node) syncSuccess(msg *SyncMsg) error {
	if n.cur.cur != 0 {
		flag := n.cur.switchCur(n.cur.cur, 0)
		if !flag {
			return errors.Errorf("switch fail.")
		}
	}

	if n.cur.waitingState != nil {
		n.cur.waitingState.Destroy()
	}
	n.cur.reset()
	n.state.seq = msg.seq
	var logs []*ReqMsg
	for _, v := range msg.logs {
		logs = append(logs, v)
	}
	n.state.logs = logs
	fmt.Printf("sync success, seq:%d.\n", msg.seq)
	return nil
}

func (n *Node) syncTo(targetSeq int64) error {
	msg := &SyncMsg{
		nodeMsg: nodeMsg{
			viewId:  n.state.viewId,
			seq:     targetSeq,
			fromIdx: n.idx,
		},
		tp:   0,
		hash: "",
		logs: nil,
	}

	n.broadcast(msg)
	return nil
}

func (n *Node) switchToCommit(msg *PPMsg) error {
	flag := n.cur.switchCur(1, 2)
	if !flag {
		return errors.Errorf("switch fail.")
	}
	old := n.cur.waitingState
	old.Destroy()

	n.cur.waitingState = &stateCounterImpl{
		okCnt:     2*n.f + 1,
		done:      make(map[string]map[uint32]struct{}),
		timeoutCh: make(chan struct{}),
	}
	n.cur.ppMsg = msg
	n.cur.waitingState.Done("", n.idx)
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
		err := n.net.SendTo(n.idx, v.idx, msg)
		if err != nil {
			log15.Error(fmt.Sprintf("Send msg error %s - %d", err.Error(), n.idx))
			continue
		}
		//ch, _ := v.GetCh()
		//ch <- msg
	}
	return nil
}
func (n *Node) broadcastToClients(msg interface{}) error {
	for _, v := range n.clients {
		err := n.net.SendTo(n.idx, v.idx, msg)
		if err != nil {
			log15.Error(fmt.Sprintf("Send msg error %s - %d", err.Error(), n.idx))
			continue
		}
		//ch, _ := v.GetCh()
		//ch <- msg
	}
	return nil
}

func (n *Node) isPrimary() bool {
	return n.state.viewId%int64(n.clusterSize) == int64(n.idx)
}

func (n *Node) onPrePrepareMsg(m *PPMsg) error {
	fmt.Printf("[%d]on pp msg:%v\n", n.idx, m)
	if n.isPrimary() {
		return nil
	}
	if n.state.seq+1 != m.seq {
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
	n.cur.waitingState.Done("", m.fromIdx)
	msg := &BFTMsg{
		nodeMsg: nodeMsg{
			viewId:  m.viewId,
			seq:     m.seq,
			fromIdx: n.idx,
		},
		tp:        0,
		hash:      m.hash,
		signature: nil,
	}
	n.broadcast(msg)
	return nil
}

func (n *Node) onPrepareMsg(m *BFTMsg) error {
	fmt.Printf("[%d]on prepare msg:%v\n", n.idx, m)
	if n.cur.cur != 1 {
		// put to buffer
		return nil
	}

	pp := n.cur.ppMsg
	if pp.hash == m.hash &&
		pp.seq == m.seq &&
		pp.viewId == m.viewId {
		n.cur.waitingState.Done("", m.fromIdx)
		if n.cur.waitingState.Ok("") {
			err := n.switchToCommit(pp)
			if err != nil {
				log15.Error(err.Error())
				return nil
			}
			msg := &BFTMsg{
				nodeMsg: nodeMsg{
					viewId:  m.viewId,
					seq:     m.seq,
					fromIdx: n.idx,
				},
				tp:        1,
				hash:      m.hash,
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
	fmt.Printf("[%d]on commit msg:%v\n", n.idx, m)
	if n.cur.cur != 2 {
		return nil
	}

	pp := n.cur.ppMsg
	if pp.hash == m.hash &&
		pp.seq == m.seq &&
		pp.viewId == m.viewId {
		n.cur.waitingState.Done("", m.fromIdx)
		if n.cur.waitingState.Ok("") {
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

func (n *Node) onRequestSyncMsg(m *SyncMsg) error {
	fmt.Printf("[%d]on sync request msg:%v\n", n.idx, m)
	if m.viewId != n.state.viewId {
		return nil
	}
	if m.seq > n.state.seq {
		return nil
	}
	logs, hash := n.state.getLogsBySeq(m.seq)

	result := &SyncMsg{
		nodeMsg: nodeMsg{
			viewId:  n.state.viewId,
			seq:     m.seq,
			fromIdx: n.idx,
		},
		tp:   1,
		hash: hash,
		logs: logs,
	}

	err := n.net.SendTo(n.idx, m.fromIdx, result)
	if err != nil {
		log15.Error(fmt.Sprintf("Send msg error %s - %d", err.Error(), n.idx))
	}
	return nil
}

func (n *Node) onResponseSyncMsg(m *SyncMsg) error {
	fmt.Printf("[%d]on sync response msg:%v\n", n.idx, m)
	if m.seq <= n.state.seq {
		return nil
	}
	key := fmt.Sprintf("sync:%d-%s", m.seq, m.hash)
	n.syncCounter.Done(key, m.fromIdx)
	if n.syncCounter.Ok(key) {
		n.syncSuccess(m)
	}
	return nil
}

func (n *Node) waitingSync() error {
	for n.mBuf.syncBuf.Len() > 0 {
		m := n.mBuf.syncBuf.Dequeue().(*SyncMsg)
		n.onResponseSyncMsg(m)
	}
	return nil
}

type Cli struct {
	clusterSize int32
	f           int
	idx         uint32
	viewId      int64
	ch          chan interface{}
	peers       map[uint32]*Node

	net Net

	waitingReply *stateCounterImpl
	cur          *ReqMsg
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
		hash:      calculateHash([]byte(op)),
	}
	c.cur = msg

	c.waitingReply = &stateCounterImpl{
		okCnt:     c.f + 1,
		done:      make(map[string]map[uint32]struct{}),
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

	errMsg := c.net.SendTo(c.idx, node.idx, msg)
	if errMsg != nil {
		return errMsg
	}
	//ch, _ := node.GetCh()
	//ch <- msg
	err := <-waitingCh
	c.waitingReply = nil
	c.cur = nil
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
	if impl == nil || c.cur == nil {
		return nil
	}

	if c.cur.hash == m.hash {
		impl.Done("", m.fromIdx)
		if impl.Ok("") {
			impl.Destroy()
		}
	}
	q := queue.New()
	q.Peek()
	return nil
}

type msgBuffer struct {
	// viewId + hash
	syncBuf    *queue.Queue
	reqBuf     *queue.Queue
	ppBuf      *queue.Queue
	prepareBuf *queue.Queue
	commitBuf  *queue.Queue
	replyBuf   *queue.Queue
}
