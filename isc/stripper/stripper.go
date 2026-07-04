// Package stripper removes comments from InterSpec (.is) source files.
//
// It implements a state machine that tracks whether the scanner is
// inside code, a string literal, a single-line comment (//), or a
// multi-line block comment (/* */). String literals are preserved
// verbatim — comment delimiters inside "..." are not treated as comments.
//
// Comment-only lines (lines that contain only a comment and optional
// leading whitespace) are collapsed — their trailing newline is suppressed
// so they don't leave blank lines in the output. Genuinely blank lines
// (no comment, no content at all) are preserved.
package stripper

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// State represents the current lexer state.
type State int

const (
	stateCode         State = iota // normal code
	stateString                    // inside "..." string literal
	stateStringEscape              // inside string, after backslash
	stateLineComment               // after //, until \n
	stateBlockComment              // after /*, until */
)

// UnterminatedCommentError is returned when a /* block comment reaches EOF.
type UnterminatedCommentError struct {
	File string
	Line int
	Col  int
}

func (e *UnterminatedCommentError) Error() string {
	if e.File != "" {
		return fmt.Sprintf("%s:%d:%d: unterminated block comment", e.File, e.Line, e.Col)
	}
	return fmt.Sprintf("line %d, col %d: unterminated block comment", e.Line, e.Col)
}

// Stripper removes comments from InterSpec source code.
type Stripper struct {
	file string // optional file path, used in error messages
}

// New creates a new Stripper. If file is non-empty, it's used in error messages.
func New(file string) *Stripper {
	return &Stripper{file: file}
}

// Strip reads from r, strips all comments, and writes the result to w.
// It returns an error if the input contains an unterminated block comment
// or if a read/write error occurs.
func (s *Stripper) Strip(r io.Reader, w io.Writer) error {
	br := bufio.NewReader(r)

	state := stateCode
	line := 1
	col := 0

	// Tracks start position of an unterminated block comment for error reporting.
	var blockStartLine, blockStartCol int

	// Output buffer
	var buf strings.Builder
	buf.Grow(4096)

	// Pending whitespace buffer — whitespace at the start of a line or before
	// a comment is buffered here and only flushed when we know it's needed.
	var pendingWs strings.Builder

	// Line tracking flags
	lineHasContent := false // true if any non-whitespace code was emitted this line
	lineHadComment := false // true if we entered a comment on this line

	// flushWs writes pending whitespace to the main buffer.
	flushWs := func() {
		buf.WriteString(pendingWs.String())
		pendingWs.Reset()
	}

	// discardWs drops pending whitespace (used when entering a comment on an empty line).
	discardWs := func() {
		pendingWs.Reset()
	}

	// emitCode writes a non-whitespace code character to the buffer.
	emitCode := func(b byte) {
		flushWs()
		buf.WriteByte(b)
		lineHasContent = true
	}

	// emitNewline handles a newline in code state.
	emitNewline := func() {
		if lineHasContent {
			// Normal line with content: flush whitespace and emit \n
			flushWs()
			buf.WriteByte('\n')
		} else if lineHadComment {
			// Comment-only line: suppress the newline, discard any pending whitespace
			discardWs()
		} else {
			// Genuinely blank line: preserve it
			flushWs()
			buf.WriteByte('\n')
		}
		lineHasContent = false
		lineHadComment = false
	}

	// flush writes the buffer to the output writer.
	flush := func() error {
		_, err := w.Write([]byte(buf.String()))
		buf.Reset()
		return err
	}

	for {
		b, err := br.ReadByte()
		if err != nil {
			if err == io.EOF {
				if state == stateBlockComment {
					return &UnterminatedCommentError{
						File: s.file,
						Line: blockStartLine,
						Col:  blockStartCol,
					}
				}
				// Flush remaining buffer
				flushWs()
				return flush()
			}
			return err
		}

		// Track position
		col++
		if b == '\n' {
			line++
			col = 0
		}

		switch state {
		case stateCode:
			switch {
			case b == '"':
				flushWs()
				buf.WriteByte(b)
				lineHasContent = true
				state = stateString

			case b == '/':
				// Peek next byte to check for // or /*
				next, err := br.ReadByte()
				if err == io.EOF {
					// Lone '/' at end of file — emit it
					emitCode(b)
					if err2 := flush(); err2 != nil {
						return err2
					}
					return nil
				}
				if err != nil {
					return err
				}

				switch next {
				case '/':
					// Single-line comment
					if lineHasContent {
						flushWs()
					} else {
						discardWs()
					}
					lineHadComment = true
					state = stateLineComment

				case '*':
					// Block comment
					if lineHasContent {
						flushWs()
					} else {
						discardWs()
					}
					lineHadComment = true
					state = stateBlockComment
					blockStartLine = line
					blockStartCol = col
					col++ // for the * we already consumed

				default:
					// Not a comment — emit '/' and put back the peeked byte
					emitCode(b)
					if err := br.UnreadByte(); err != nil {
						return err
					}
				}

			case b == ' ' || b == '\t':
				// Buffer whitespace; don't emit yet
				pendingWs.WriteByte(b)

			case b == '\n':
				emitNewline()

			default:
				// Regular code character
				emitCode(b)
			}

		case stateString:
			switch {
			case b == '\\':
				buf.WriteByte(b)
				state = stateStringEscape
			case b == '"':
				buf.WriteByte(b)
				state = stateCode
			default:
				buf.WriteByte(b)
			}

		case stateStringEscape:
			buf.WriteByte(b)
			state = stateString

		case stateLineComment:
			if b == '\n' {
				if lineHasContent {
					// Line had code before the comment — preserve newline
					buf.WriteByte('\n')
				}
				// Comment-only line: suppress the newline
				lineHasContent = false
				lineHadComment = false
				state = stateCode
			}
			// Everything else is consumed without emitting

		case stateBlockComment:
			if b == '*' {
				next, err := br.ReadByte()
				if err == io.EOF {
					return &UnterminatedCommentError{
						File: s.file,
						Line: blockStartLine,
						Col:  blockStartCol,
					}
				}
				if err != nil {
					return err
				}
				if next == '/' {
					col++ // for '/'
					state = stateCode
					// Don't reset line flags here — we might have content
					// after the block comment on the same line.
				} else {
					if err := br.UnreadByte(); err != nil {
						return err
					}
				}
			}
			// Newlines inside block comments are tracked for position but
			// consumed silently. They don't reset line flags because the
			// block comment hasn't ended yet.
		}
	}
}
