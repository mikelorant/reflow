package padding

import (
	"io"

	"github.com/mattn/go-runewidth"

	"github.com/muesli/reflow/ansi"
)

type PaddingFunc func(w io.Writer)

type Writer struct {
	Forward *ansi.Writer
	Padding uint
	PadFunc PaddingFunc

	lineLen int
	ansi    bool
}

func NewWriter(w *ansi.Writer, width uint, paddingFunc PaddingFunc) *Writer {
	return &Writer{
		Forward: w,
		Padding: width,
		PadFunc: paddingFunc,
	}
}

// Write is used to write content to the padding buffer.
func (w *Writer) Write(b []byte) (int, error) {
	for _, c := range string(b) {
		if c == '\x1B' {
			// ANSI escape sequence
			w.ansi = true
		} else if w.ansi {
			if (c >= 0x41 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a) {
				// ANSI sequence terminated
				w.ansi = false
			}
		} else {
			w.lineLen += runewidth.StringWidth(string(c))

			if c == '\n' {
				// end of current line
				if w.Padding > 0 && uint(w.lineLen) < w.Padding {
					for i := 0; i < int(w.Padding)-w.lineLen; i++ {
						w.PadFunc(w.Forward)
					}
				}
				w.Forward.ResetAnsi()
				w.lineLen = 0
			}
		}

		_, err := w.Forward.Write([]byte(string(c)))
		if err != nil {
			return 0, err
		}
	}

	return len(b), nil
}
