package basic

import (
	"fmt"
	"testing"
)

func TestSwitch(t *testing.T) {
	switch 1 {
	case 1:
		fmt.Println("case 1.")
	case 2:
		fmt.Println("case 2.")
	default:
		fmt.Println("case default.")
	}
}
