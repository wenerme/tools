package apk

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestFetchJoin(t *testing.T) {
	assert.Equal(t, fetchJoin("/", "/a"), "/a")
	assert.Equal(t, fetchJoin("http://a/", "/a"), "http://a/a")
}
