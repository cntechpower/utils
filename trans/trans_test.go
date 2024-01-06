package trans

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNvl(t *testing.T) {
	var a int64
	var b string
	a = 1
	b = "test"

	assert.Equal(t, StringNvl(nil), Nvl[string](nil))
	assert.Equal(t, StringNvl(&b), Nvl[string](&b))
	assert.Equal(t, Int64Nvl(nil), Nvl[int64](nil))
	assert.Equal(t, Int64Nvl(&a), Nvl[int64](&a))
}
