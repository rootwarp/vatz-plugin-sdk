package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	p := NewPlugin()

	p.RegisterFeature(callbackFunc)

	assert.Equal(t, 1, len(p.callbacks))
}

func callbackFunc() error {
	return nil
}
