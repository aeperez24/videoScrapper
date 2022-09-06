package cuevana

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveElement(t *testing.T) {
	client, proxy := getHttpClientWithProxy([]string{"proxy"})
	assert.Equal(t, "proxy", proxy)
	assert.NotNil(t, client)
}
