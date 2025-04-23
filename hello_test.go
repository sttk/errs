package errs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sttk/errs"
)

func TestHello(t *testing.T) {
	assert.Equal(t, errs.Hello(), "hello")
}
