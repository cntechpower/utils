package net

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetLocalIps(t *testing.T) {
	ip, err := GetFirstLocalIp()
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", ip)
	t.Logf(ip)
}
