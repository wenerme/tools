package apk

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestDependency(t *testing.T) {
	for _, test := range []struct {
		raw string
		dep Dependency
	}{
		{"xxd", Dependency{Name: "xxd"}},
		{"cmd:eview", Dependency{Type: "cmd", Name: "eview"}},
		{"pc:cogl-1.0=1.22.2", Dependency{Type: "pc", Name: "cogl-1.0", Version: "1.22.2"}},
	} {
		assert.Equal(t, ParseDependency(test.raw), test.dep)
		assert.Equal(t, test.raw, test.dep.String())
	}
}
