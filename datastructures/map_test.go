package datastructures

import (
	"fmt"
	"strconv"
	"testing"
)

type S1 interface {
	s()
}

type S2 struct {
}

func (self *S2) s() {

}

func TestMap(t *testing.T) {
	m := make(map[string]*S2)
	fmt.Println(get(m))
	fmt.Println(get(m) == nil)
	fmt.Println(get2(m))
	fmt.Println(get2(m) == nil)
	fmt.Println(get3(m) == nil)
	fmt.Println(m["2"] == nil)
}

func get(m map[string]*S2) S1 {
	s2 := m["2"]
	return s2
}

func get2(m map[string]*S2) S1 {
	s2, ok := m["2"]
	if !ok {
		return nil
	} else {
		return s2
	}
}

func get3(m map[string]*S2) *S2 {
	s2 := m["2"]
	return s2
}

func TestAppend(t *testing.T) {
	m := make(map[string][]string)

	for i := 0; i < 10; i++ {
		s := m[strconv.Itoa(i)]
		for j := 0; j < 10; j++ {
			s = append(s, strconv.Itoa(j))
		}
		m[strconv.Itoa(i)] = s
	}

	fmt.Printf("%v", m)
}
