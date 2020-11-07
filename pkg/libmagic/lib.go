package libmagic

/*
   #cgo LDFLAGS: -lmagic
   #include <magic.h>
   #include <stdlib.h>

static magic_t f(void* x) {
	return (magic_t)x;
}
*/
import "C"
import (
	"errors"
	"os"
	"path"
	"unsafe"
)

// nolint: golint
const (
	MAGIC_NONE              = C.MAGIC_NONE
	MAGIC_DEBUG             = C.MAGIC_DEBUG
	MAGIC_SYMLINK           = C.MAGIC_SYMLINK
	MAGIC_COMPRESS          = C.MAGIC_COMPRESS
	MAGIC_DEVICES           = C.MAGIC_DEVICES
	MAGIC_MIME_TYPE         = C.MAGIC_MIME_TYPE
	MAGIC_CONTINUE          = C.MAGIC_CONTINUE
	MAGIC_CHECK             = C.MAGIC_CHECK
	MAGIC_PRESERVE_ATIME    = C.MAGIC_PRESERVE_ATIME
	MAGIC_RAW               = C.MAGIC_RAW
	MAGIC_ERROR             = C.MAGIC_ERROR
	MAGIC_MIME_ENCODING     = C.MAGIC_MIME_ENCODING
	MAGIC_MIME              = C.MAGIC_MIME
	MAGIC_APPLE             = C.MAGIC_APPLE
	MAGIC_EXTENSION         = C.MAGIC_EXTENSION
	MAGIC_COMPRESS_TRANSP   = C.MAGIC_COMPRESS_TRANSP
	MAGIC_NO_CHECK_COMPRESS = C.MAGIC_NO_CHECK_COMPRESS
	MAGIC_NO_CHECK_TAR      = C.MAGIC_NO_CHECK_TAR
	MAGIC_NO_CHECK_SOFT     = C.MAGIC_NO_CHECK_SOFT
	MAGIC_NO_CHECK_APPTYPE  = C.MAGIC_NO_CHECK_APPTYPE
	MAGIC_NO_CHECK_ELF      = C.MAGIC_NO_CHECK_ELF
	MAGIC_NO_CHECK_TEXT     = C.MAGIC_NO_CHECK_TEXT
	MAGIC_NO_CHECK_CDF      = C.MAGIC_NO_CHECK_CDF
	MAGIC_NO_CHECK_TOKENS   = C.MAGIC_NO_CHECK_TOKENS
	MAGIC_NO_CHECK_ENCODING = C.MAGIC_NO_CHECK_ENCODING
	MAGIC_NO_CHECK_ASCII    = C.MAGIC_NO_CHECK_ASCII
	MAGIC_NO_CHECK_FORTRAN  = C.MAGIC_NO_CHECK_FORTRAN
	MAGIC_NO_CHECK_TROFF    = C.MAGIC_NO_CHECK_TROFF
)

// nolint: golint
const (
	MAGIC_NO_CHECK_BUILTIN = MAGIC_NO_CHECK_COMPRESS |
		MAGIC_NO_CHECK_TAR |
		MAGIC_NO_CHECK_APPTYPE |
		MAGIC_NO_CHECK_ELF |
		MAGIC_NO_CHECK_TEXT |
		MAGIC_NO_CHECK_CDF |
		MAGIC_NO_CHECK_TOKENS |
		MAGIC_NO_CHECK_ENCODING
)

type Magic uintptr

var locations = []string{
	"/usr/share/misc/magic.mgc",
	"/usr/share/file/magic.mgc",
	"/usr/share/magic/magic.mgc",
}

/* Find the real magic file location */
func GetDefaultDir() string {
	var f string
	found := false
	for _, f = range locations {
		fi, err := os.Lstat(f)
		if err == nil && fi.Mode()&os.ModeSymlink != os.ModeSymlink {
			found = true
			break
		}
	}
	if found {
		return path.Dir(f)
	}
	return ""
}

func Open(flags int) Magic {
	cookie := Magic(unsafe.Pointer(C.magic_open(C.int(flags))))
	return cookie
}

func (m Magic) Close() error {
	C.magic_close(magic(m))
	return nil
}

func (m Magic) Error() string {
	s := C.magic_error(magic(m))
	return C.GoString(s)
}
func magic(x Magic) C.magic_t {
	return C.magic_t(unsafe.Pointer(x))
}

func (m Magic) Errno() int {
	return (int)(C.magic_errno(magic(m)))
}

func (m Magic) File(filename string) string {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	return C.GoString(C.magic_file(magic(m), cfilename))
}

func (m Magic) Buffer(b []byte) string {
	length := C.size_t(len(b))
	return C.GoString(C.magic_buffer(magic(m), unsafe.Pointer(&b[0]), length))
}

func (m Magic) SetFlags(flags int) int {
	return (int)(C.magic_setflags(magic(m), C.int(flags)))
}

func (m Magic) GetFlags() int {
	return (int)(C.magic_getflags(magic(m)))
}

func (m Magic) Check(filename string) error {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	return m.err((int)(C.magic_check(magic(m), cfilename)))
}

func (m Magic) Compile(filename string) error {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	return m.err((int)(C.magic_compile(magic(m), cfilename)))
}

func (m Magic) Load(filename string) error {
	if filename == "" {
		return m.err((int)(C.magic_load(magic(m), nil)))
	}
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	return m.err((int)(C.magic_load(magic(m), cfilename)))
}
func (m Magic) err(errno int) error {
	if errno == 0 {
		return nil
	}
	err := m.Error()
	if err == "" {
		return nil
	}
	return errors.New(err)
}

func Version() int {
	return (int)(C.magic_version())
}
