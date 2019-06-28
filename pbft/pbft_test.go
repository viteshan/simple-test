package pbft

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/golang-collections/collections/queue"
)

func TestPbft(t *testing.T) {

	nn := newNet()

	clusterSize := int32(4)
	f := int(1)
	var nodes []*Node
	for i := int32(0); i < clusterSize; i++ {
		var logs []*ReqMsg
		blacklist, err := NewBlacklist()
		if err != nil {
			t.Fatal(err)
		}
		n := &Node{
			clusterSize: clusterSize,
			f:           f,
			idx:         uint32(i),
			state: &State{
				viewId: 0,
				seq:    0,
				logs:   logs,
			},
			ch:      make(chan interface{}, 100),
			peers:   make(map[uint32]*Node),
			clients: make(map[uint32]*Cli),
			cur:     &CurState{},
			net:     nn,
			mBuf: &msgBuffer{
				syncBuf:    queue.New(),
				reqBuf:     queue.New(),
				ppBuf:      queue.New(),
				prepareBuf: queue.New(),
				commitBuf:  queue.New(),
				replyBuf:   queue.New(),
			},
			syncCounter: &stateCounterImpl{
				okCnt:     2*f + 1,
				done:      make(map[string]map[uint32]struct{}),
				timeoutCh: nil,
			},
			timeout: blacklist,
		}

		err = nn.Register(n.idx, n.ch)
		if err != nil {
			t.Error("Register error ", err)
			t.FailNow()
		}

		nodes = append(nodes, n)
	}
	c := &Cli{
		clusterSize:  clusterSize,
		f:            f,
		idx:          100,
		viewId:       0,
		ch:           make(chan interface{}),
		peers:        make(map[uint32]*Node),
		waitingReply: nil,
		net:          nn,
	}
	errC := nn.Register(c.idx, c.ch)
	if errC != nil {
		t.Error("Register error ", errC)
		t.FailNow()
	}

	for _, n := range nodes {
		for _, nn := range nodes {
			if n.idx == nn.idx {
				continue
			}
			n.peers[nn.idx] = nn
		}
		n.clients[c.idx] = c
		c.peers[n.idx] = n
	}
	go c.loopRead()
	for _, n := range nodes {
		go n.loopRead()
		go n.loopBuf()
		go n.loopHeartBeat()
	}

	nodes[1].down = true
	err := c.SendRequest("hello")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	for _, v := range nodes {
		t.Log(v.idx, v.cur.cur, v.state.seq, len(v.state.logs))
	}

	nodes[1].down = false

	err = c.SendRequest("hello world2")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ok := true
	s := int64(-1)

	for i := 0; i < 10; i++ {
		ok = true
		for _, v := range nodes {
			if s < 0 {
				s = v.state.seq
			}
			if s != v.state.seq {
				ok = false
			}
			t.Log(v.idx, v.cur.cur, v.state.seq, len(v.state.logs))
		}
		if ok {
			break
		}
		time.Sleep(time.Second)
	}

	assert.True(t, ok)
}
