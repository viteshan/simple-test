package main

import (
	"fmt"
	"math/big"
	"testing"
)

func TestNum(t *testing.T) {

	a := big.NewInt(5)
	b := big.NewInt(2)
	c := big.NewInt(26)
	d := big.NewInt(0).Exp(a, b, c)

	fmt.Println(a.String(), b.String(), c.String(), d.String())
}
