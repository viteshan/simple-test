package dag

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/viteshan/go/support/errors"
)

type pool struct {
	chains
}

type block struct {
	height uint64
	prev   string
	hash   string
	addr   string
	refers []string
}

type chains struct {
	cs sync.Map
}

type chain struct {
	blocks map[uint64]*block
	last   *block
}

func (self *chain) add(b *block) error {
	last := self.last
	if last == nil && len(self.blocks) != 0 {
		return errors.New("block must be empty")
	}

	if last != nil {
		if last.hash != b.prev {
			return errors.New("hash and prev")
		}
	}

	self.blocks[b.height] = b
	self.last = b
	return nil
}

type bucket struct {
	bs   []*block
	last int
}

func (self *bucket) add(b *block) error {
	if self.last == -1 && len(self.bs) != 0 {
		return errors.New("bucket must be empty")
	}

	if self.last != -1 {
		lastB := self.bs[self.last]
		if lastB == nil {
			return errors.New("lastB must be exists")
		}
		if lastB.hash != b.prev {
			return errors.New("prev and hash")
		}
	}

	self.bs = append(self.bs, b)
	self.last = self.last + 1
	return nil
}
func (self *bucket) print() {
	for _, v := range self.bs {
		fmt.Print(strconv.FormatUint(v.height, 10) + ",")
	}
	fmt.Println()
}

type level struct {
	bs map[string]*bucket
}

func newLevel() *level {
	return &level{bs: make(map[string]*bucket)}
}

func (self *level) add(b *block) error {
	bu, ok := self.bs[b.addr]
	if !ok {
		self.bs[b.addr] = newBucket()
		bu = self.bs[b.addr]
	}
	return bu.add(b)
}
func (self *level) print() {
	for k, v := range self.bs {
		fmt.Println("--------Bucket[" + k + "]----------")
		v.print()
	}

}
func newBucket() *bucket {
	return &bucket{last: -1}
}

type addrLevel struct {
	addr  string
	level int
}

type queue struct {
	all map[string]*addrLevel
	ls  []*level
}

func newTask() *queue {
	tmpLs := make([]*level, 10)
	for i := 0; i < 10; i++ {
		tmpLs[i] = newLevel()
	}
	return &queue{all: make(map[string]*addrLevel), ls: tmpLs}
}

func (self *queue) add(b *block) error {
	max := 0
	for _, r := range b.refers {
		tmp, ok := self.all[r]
		if !ok {
			continue
		}
		lNum := tmp.level
		if lNum >= max {
			if tmp.addr == b.addr {
				max = lNum
			} else {
				max = lNum + 1
			}
		}
	}
	if max > 9 {
		return errors.New("arrived to max")
	}
	self.all[b.hash] = &addrLevel{b.addr, max}

	return self.ls[max].add(b)
}
func (self *queue) print() {
	for i, v := range self.ls {
		fmt.Println("----------Level" + strconv.Itoa(i) + "------------")
		v.print()
	}
}
