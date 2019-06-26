package basic

import "testing"

func TestName(t *testing.T) {
	a3 := A3{&A2{test: 3}}

	a1 := A1{a2: a3.A2}

	a3.test = 4

	println(a1.a2.test)
}

type A1 struct {
	a2 *A2
}

type A2 struct {
	test int
}

type A3 struct {
	*A2
}
