package stardict

import (
	"io"
	"sort"
	"time"
)

type Dict struct {
	Info  *DictInfo
	Index DictIndex

	r       io.ReadSeeker
	closers []io.Closer
}

func (d *Dict) Close() error {
	for _, c := range d.closers {
		_ = c.Close()
	}

	return nil
}
func (d *Dict) Search(x string) (*DictIndexEntry, error) {
	idx := d.Index.Search(x)
	if idx != nil {
		return idx, readDictContent(d, idx)
	}
	return nil, nil
}

type DictIndex []DictIndexEntry

func (idx DictIndex) Search(x string) *DictIndexEntry {
	i := sort.Search(len(idx), func(i int) bool {
		return idx[i].Word >= x
	})
	if i < len(idx) && idx[i].Word == x {
		return &idx[i]
	}
	return nil
}

type DictInfo struct {
	Header           string // magic StarDict's dict ifo file
	Version          string // "2.4.2" or "3.0.0"
	WordCount        int
	SynWordCount     int
	IndexFileSize    int
	IndexOffsetBits  int // since 3.0.0 Bits pre offset, 32/64
	BookName         string
	Description      string // <br> for new line
	Date             *time.Time
	SameTypeSequence string
	DictType         string

	Author  string
	Email   string
	Website string
}

type DictIndexEntry struct {
	Word string // a utf-8 string terminated by '\0'.
	// word data's offset in .dict file<br>
	// If the version is "3.0.0" and "idxoffsetbits=64", word_data_offset will be 64-bits unsigned number in network byte order.
	// Otherwise it will be 32-bits.
	Offset int64
	// word data's total size in .dict file<br>
	// word_data_size should be 32-bits unsigned number in network byte order.
	Size     int
	Synonyms []string

	Contents []DictEntry
}

// Lower-case characters signify that a field's size is determined by a
// terminating '\0', while upper-case characters indicate that the data
// begins with a network byte-ordered guint32 that gives the length of
// the following data's size (NOT the whole size which is 4 bytes bigger).
type EntryType rune

const (
	// Word's pure text meaning.
	// The data should be a utf-8 string ending with '\0'.
	TypeNullTerminalText EntryType = 'm'

	// Word's pure text meaning.
	// The data is NOT a utf-8 string, but is instead a string in locale
	// encoding, ending with '\0'. Sometimes using this type will save disk
	// space, but its use is discouraged. This is only a idea.
	TypeLocaleText EntryType = 'l'

	// A utf-8 string which is marked up with the Pango text markup language.
	// For more information about this markup language, See the "Pango
	// Reference Manual."
	// You might have it installed locally at:
	// file:///usr/share/gtk-doc/html/pango/PangoMarkupFormat.html
	// Online:
	// http://library.gnome.org/devel/pango/stable/PangoMarkupFormat.html
	TypePangoText EntryType = 'g'

	// English phonetic string.
	// The data should be a utf-8 string ending with '\0'.
	//
	// Here are some utf-8 phonetic characters:
	//
	// Î¸ÊƒÅ‹Ê§Ã°Ê’Ã¦Ä±ÊŒÊŠÉ’É›É™É‘ÉœÉ”ËŒËˆËË‘á¹ƒá¹‡á¸·
	// Ã¦É‘É’ÊŒÓ™Ñ”Å‹vÎ¸Ã°ÊƒÊ’ÉšËÉ¡ËËŠË‹
	TypeEnglishPhonetic EntryType = 't'

	// A utf-8 string which is marked up with the xdxf language.
	// StarDict have these extension:
	// <rref> can have "type" attribute, it can be "image", "sound", "video"
	// and "attach".
	// <kref> can have "k" attribute.
	//
	// http://xdxf.sourceforge.net
	TypeXdxfMarkup EntryType = 'x'
	// Chinese YinBiao or Japanese KANA.
	// The data should be a utf-8 string ending with '\0'.
	TypeYinBiao EntryType = 'y'
	// KingSoft PowerWord's data. The data is a utf-8 string ending with '\0'.
	// It is in XML format.
	TypeKingsoftXML EntryType = 'k'
	// MediaWiki markup language.
	// http://meta.wikimedia.org/wiki/Help:Editing#The_wiki_markup
	MediawikiMarkupType EntryType = 'w'
	// Html codes.
	TypeHTML EntryType = 'h'
	// WordNet data.
	TypeWordNet EntryType = 'n'
	// Resource file list.
	// The content can be:
	// img:pic/example.jpg	// Image file
	// snd:apple.wav		// Sound file
	// vdo:film.avi		// Video file
	// att:file.bin		// Attachment file
	// More than one line is supported as a list of available files.
	// StarDict will find the files in the Resource Storage.
	// The image will be shown, the sound file will have a play button.
	// You can "save as" the attachment file and so on.
	// The file list must be a utf-8 string ending with '\0'.
	// Use '\n' for separating new lines.
	// Use '/' character as directory separator.
	TypeResource EntryType = 'r'
	// wav file.
	// The data begins with a network byte-ordered guint32 to identify the wav
	// file's size, immediately followed by the file's content.
	// This is only a idea, it is better to use 'r' Resource file list in most
	// case.
	TypeWavFile EntryType = 'W'
	// Picture file.
	// The data begins with a network byte-ordered guint32 to identify the picture
	// file's size, immediately followed by the file's content.
	// This feature is implemented, as stardict-advertisement-plugin needs it.
	// Anyway, it is better to use 'r' Resource file list in most case.
	TypePictureFile EntryType = 'P'
	// this type identifier is reserved for experimental extensions.
	TypeReserved EntryType = 'X'
)

func (v EntryType) String() string {
	return string(v)
}
