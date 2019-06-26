package basic

import (
	"testing"
	"time"
)

func TestSelect(t1 *testing.T) {
	i := 0
	j := 0
	t := time.NewTicker(time.Millisecond * 20)
	for {
		select {
		case <-t.C:
			i++
			j++
		default:
			i++
		}

		if j == 100 {
			break
		}
	}

	println(i)
	println(j)
}
