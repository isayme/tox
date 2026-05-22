package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNowInMills(t *testing.T) {
	before := time.Now().UnixMilli()
	mills := NowInMills()
	after := time.Now().UnixMilli()

	assert.GreaterOrEqual(t, mills, before)
	assert.LessOrEqual(t, mills, after)
}
