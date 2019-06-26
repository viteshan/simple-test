package pbft

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNet(t *testing.T) {
	n := newNet()

	ch := make(chan interface{}, 10)

	err := n.Register(1, ch)
	if err != nil {
		t.Error("Register error", err)
		t.FailNow()
	}

	err = n.Register(2, ch)
	if err != nil {
		t.Error("Register error", err)
		t.FailNow()
	}

	timeNow := time.Now()
	msg := "hello world"
	errMsg := n.SendTo(1, 2, msg)
	if errMsg != nil {
		t.Error("Send msg error", errMsg)
		t.FailNow()
	}

	select {
	case m := <-ch:
		t.Log("Msg: ", m)
		assert.Equal(t, msg, m)
	case <-time.After(time.Second * 2):
		t.Error("Timeout msg")
		t.FailNow()
	}

	timeEnd := time.Now()
	t.Log("Duration time ", timeEnd.Sub(timeNow))
}
