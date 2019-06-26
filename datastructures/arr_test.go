package datastructures

import (
	"strconv"
	"testing"
)

func TestArr(t *testing.T) {
	var hashes []string

	for i := 0; i < 10; i++ {
		hashes = append(hashes, strconv.Itoa(i))
	}

	splits := split(hashes, 10)

	for _, a := range splits {
		for _, b := range a {
			print(b)
			print(",")
		}
		println()
	}

}

func split(buf []string, lim int) [][]string {
	var chunk []string
	chunks := make([][]string, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
