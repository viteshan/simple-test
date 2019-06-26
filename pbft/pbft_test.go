package pbft

import "testing"

func TestPbft(t *testing.T) {

	clusterSize := int32(4)
	f := int(1)
	var nodes []*Node
	for i := int32(0); i < clusterSize; i++ {
		var logs []*ReqMsg
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
	}

	err := c.SendRequest("hello")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	for _, v := range nodes {
		t.Log(v.idx, v.cur.cur, v.state.seq, len(v.state.logs))
	}
}
