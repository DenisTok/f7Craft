package models

import (
	"testing"
)

func TestRandEmojis(t *testing.T) {
	b, err := RandEmojis(10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(b)
	return
}
