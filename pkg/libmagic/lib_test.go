package libmagic_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wenerme/tools/pkg/libmagic"
)

func TestFile(t *testing.T) {
	fmt.Println("version", libmagic.Version())

	mgc := libmagic.Open(libmagic.MAGIC_NONE)
	defer mgc.Close()
	fmt.Println(mgc.GetFlags())
	assert.NoError(t, mgc.Load(""))
	fmt.Printf("file: %s - error %#v errno %v\n", mgc.File(os.Args[0]), mgc.Error(), mgc.Errno())
	mgc.SetFlags(libmagic.MAGIC_MIME | libmagic.MAGIC_MIME_ENCODING)
	fmt.Printf("file: %s - error %#v errno %v\n", mgc.File(os.Args[0]), mgc.Error(), mgc.Errno())
}
