package pbft

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/vitelabs/go-vite/log15"
	"github.com/viteshan/go/support/errors"
)

type Net interface {
	SendTo(fromId uint32, toId uint32, msg interface{}) error

	Register(id uint32, wCh chan<- interface{}) error
	UnRegister(id uint32) error

	SetStatus(idx uint32, flag bool) // true:up false:down
}

type net struct {
	channels   map[string]chan *msgWrapper
	channelsMu sync.Mutex

	nodes map[uint32]chan<- interface{}
	status map[uint32]bool
}

type msgWrapper struct {
	msg interface{}
	to  uint32

	msgTime time.Time
}

func newNet() *net {
	return &net{
		channels: make(map[string]chan *msgWrapper),
		nodes:    make(map[uint32]chan<- interface{}),
		status: make(map[uint32]bool),
	}
}

func (n *net) getChannel(key string, toId uint32) chan *msgWrapper {
	n.channelsMu.Lock()
	defer n.channelsMu.Unlock()

	ch, ok := n.channels[key]
	if !ok {
		ch = n.newChannel(toId)
		n.channels[key] = ch
	}
	return ch
}

func (n *net) SetStatus(idx uint32, flag bool)  {
	n.status[idx] = flag
}

func (n *net) SendTo(fromId uint32, toId uint32, msg interface{}) error {

	_, ok := n.status[toId]
	if ok && !n.status[toId] {
		return errors.New("Node is down")
	}


	if !n.exist(fromId) {
		return errors.New("from is not exist")
	}
	if !n.exist(toId) {
		return errors.New("to is not exist")
	}

	key := fmt.Sprintf("%d-%d", fromId, toId)

	ch := n.getChannel(key, toId)
	select {
	case ch <- &msgWrapper{
		msg:     msg,
		to:      toId,
		msgTime: time.Now(),
	}:
		return nil
	default:
		return errors.New("full channel")
	}
}

func (n *net) exist(id uint32) bool {
	_, ok := n.nodes[id]
	return ok
}

func (n *net) newChannel(toId uint32) chan *msgWrapper {
	ch := make(chan *msgWrapper, 100)
	go func() {
		for {
			select {
			case msg := <-ch:
				if n.exist(msg.to) {
					diff := time.Now().Sub(msg.msgTime)
					delay := n.randomDelay()
					if delay > diff {
						time.Sleep(delay - diff)
					}
					select {
					case n.nodes[msg.to] <- msg.msg:
					default:
						log15.Error(fmt.Sprintf("node channel[%d] full, loss msg:%v", msg.to, msg.msg))
					}
				}
			case <-time.After(5 * time.Second):
				if !n.exist(toId) {
					fmt.Println("send channel exit", toId)
					return
				}
			}
		}
	}()
	return ch
}

func (n *net) randomDelay() time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ra := r.Intn(500)
	return time.Duration(ra) * time.Millisecond
}

func (n *net) Register(id uint32, wCh chan<- interface{}) error {
	_, ok := n.nodes[id]
	if ok {
		return errors.New("exists")
	}
	n.nodes[id] = wCh
	return nil
}

func (n *net) UnRegister(id uint32) error {
	_, ok := n.nodes[id]
	if !ok {
		return errors.New("not exists")
	}
	delete(n.nodes, id)
	return nil
}
