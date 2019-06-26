package basic

import (
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	datetime := time.Now().Format("15:04:05")
	println(datetime)
}
