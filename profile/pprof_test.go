package profile

import (
	"testing"

	"github.com/pkg/profile"
)

func TestBlock(t *testing.T) {
	defer profile.Start(profile.BlockProfile).Stop()
}
