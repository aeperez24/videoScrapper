package cuevana

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveElement(t *testing.T) {
	client, proxy := getHttpClientWithProxy([]string{"proxy"})
	assert.Equal(t, "proxy", proxy)
	assert.IsType(t, &http.Client{}, client)
}
