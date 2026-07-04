package stripper

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// strip is a convenience wrapper for testing.
func strip(input string) (string, error) {
	s := New("")
	var buf bytes.Buffer
	err := s.Strip(strings.NewReader(input), &buf)
	return buf.String(), err
}

func TestStrip_EmptyInput(t *testing.T) {
	out, err := strip("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Fatalf("expected empty output, got %q", out)
	}
}

func TestStrip_NoComments(t *testing.T) {
	input := `state greeting = "Hello, World!"
state count = 42`
	out, err := strip(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != input {
		t.Fatalf("expected:\n%q\ngot:\n%q", input, out)
	}
}

func TestStrip_SingleLineComment(t *testing.T) {
	input := `// this is a comment
state x = 1`
	// Comment-only line collapsed — no blank line left behind
	expected := "state x = 1"
	out, err := strip(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestStrip_SingleLineCommentNoNewline(t *testing.T) {
	// Comment at end of file without trailing newline
	input := `state x = 1
// trailing comment without newline`
	expected := "state x = 1\n"
	out, err := strip(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestStrip_MultiLineComment(t *testing.T) {
	input := `state x = 1
/* block comment */
state y = 2`
	// Block comment-only line collapsed
	expected := "state x = 1\nstate y = 2"
	out, err := strip(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestStrip_MultiLineCommentSpanning(t *testing.T) {
	input := `state x = 1
/* block
   comment
   spanning */
state y = 2`
	// Multi-line block comment — entire span collapsed
	expected := "state x = 1\nstate y = 2"
	out, err := strip(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestStrip_InlineBlockComment(t *testing.T) {
	input := `state x = /* inline */ 42`
	expected := "state x =  42"
	out, err := strip(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestStrip_StringWithDoubleSlash(t *testing.T) {
	input := `state url = "http://example.com"`
	// The // inside the string must NOT be treated as a comment
	out, err := strip(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != input {
		t.Fatalf("expected %q, got %q", input, out)
	}
}

func TestStrip_StringWithBlockCommentDelimiters(t *testing.T) {
	input := `state desc = "/* not a comment */ string"`
	out, err := strip(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != input {
		t.Fatalf("expected %q, got %q", input, out)
	}
}

func TestStrip_StringWithCommentSlash(t *testing.T) {
	input := `state msg = "hello // world"`
	out, err := strip(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != input {
		t.Fatalf("expected %q, got %q", input, out)
	}
}

func TestStrip_StringWithEscapedQuote(t *testing.T) {
	input := `state msg = "he said \"hello\""`
	out, err := strip(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != input {
		t.Fatalf("expected %q, got %q", input, out)
	}
}

func TestStrip_UnterminatedBlockComment(t *testing.T) {
	input := `state x = 1
/* this never closes
state y = 2`
	_, err := strip(input)
	if err == nil {
		t.Fatal("expected error for unterminated block comment, got nil")
	}
	// Verify it's the right error type
	var ucerr *UnterminatedCommentError
	if !as(err, &ucerr) {
		t.Fatalf("expected *UnterminatedCommentError, got %T: %v", err, err)
	}
	if ucerr.Line != 2 {
		t.Fatalf("expected line 2, got %d", ucerr.Line)
	}
	if ucerr.Col != 1 {
		t.Fatalf("expected col 1, got %d", ucerr.Col)
	}
}

func TestStrip_UnterminatedBlockCommentAtEOFAfterStar(t *testing.T) {
	input := `state x = 1
/* unterminated *`
	_, err := strip(input)
	if err == nil {
		t.Fatal("expected error for unterminated block comment, got nil")
	}
	var ucerr *UnterminatedCommentError
	if !as(err, &ucerr) {
		t.Fatalf("expected *UnterminatedCommentError, got %T: %v", err, err)
	}
}

func TestStrip_NestedBlockComment(t *testing.T) {
	// InterSpec (like Go/C) does NOT nest block comments.
	// The first */ closes the comment. The trailing */ is emitted as code.
	input := `/* outer /* inner */ still alive */`
	expected := " still alive */"
	out, err := strip(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestStrip_BlockCommentInsideLineComment(t *testing.T) {
	// /* is not special inside a line comment (//)
	input := `state x = 1 // line with /* not a block comment */`
	// No trailing newline — comment runs to EOF, only the space before // is emitted
	expected := "state x = 1 "
	out, err := strip(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestStrip_OnlyComments(t *testing.T) {
	input := "// just a comment\n/* another */\n"
	// All lines are comment-only — entire output collapses to empty
	expected := ""
	out, err := strip(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestStrip_Unicode(t *testing.T) {
	input := "// komentář s diakritikou\nstate x = 1\n/* 日本語 */\nstate y = 2"
	// Line 1 (line comment) and line 3 (block comment) are collapsed
	expected := "state x = 1\nstate y = 2"
	out, err := strip(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestStrip_MultipleCommentsMixed(t *testing.T) {
	input := `// header comment
page Main() {
    // inline comment
    Text("Hello") // trailing
    /* block
       comment */
    state x = /* inline */ 42
}`
	// Traced manually with comment-only line collapsing:
	//   page Main() {\n         ← line 1 (comment) collapsed, line 2 emitted
	//       Text("Hello") \n    ← line 3 (comment) collapsed, line 4 has content + trailing comment
	//       state x =  42\n     ← line 5 (block comment) collapsed, line 6 has inline comment
	//   }                       ← line 7, no trailing newline
	expected := "page Main() {\n    Text(\"Hello\") \n    state x =  42\n}"
	out, err := strip(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestStrip_LoneForwardSlash(t *testing.T) {
	// A lone / (not followed by / or *) should be emitted as-is
	input := `state x = 1 / 2`
	out, err := strip(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != input {
		t.Fatalf("expected %q, got %q", input, out)
	}
}

func TestStrip_ErrorFileField(t *testing.T) {
	s := New("testfile.is")
	input := `/* broken`
	var buf bytes.Buffer
	err := s.Strip(strings.NewReader(input), &buf)
	if err == nil {
		t.Fatal("expected error")
	}
	var ucerr *UnterminatedCommentError
	if !as(err, &ucerr) {
		t.Fatalf("expected *UnterminatedCommentError, got %T", err)
	}
	if ucerr.File != "testfile.is" {
		t.Fatalf("expected file 'testfile.is', got %q", ucerr.File)
	}
}

func TestStrip_ReadError(t *testing.T) {
	// Simulate a read error
	r := &errorReader{err: io.ErrUnexpectedEOF}
	s := New("")
	var buf bytes.Buffer
	err := s.Strip(r, &buf)
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("expected ErrUnexpectedEOF, got %v", err)
	}
}

func TestStrip_WriteError(t *testing.T) {
	// Stripper buffers internally and only flushes at EOF.
	// A writer that returns an error triggers that error on flush.
	w := &errorWriter{}
	s := New("")
	err := s.Strip(strings.NewReader("hello world"), w)
	if err == nil {
		t.Fatal("expected write error")
	}
}

// -- helpers --

// as is like errors.As but returns a bool.
func as(err error, target interface{}) bool {
	if err == nil {
		return false
	}
	// Simple type assertion for *UnterminatedCommentError
	// We can't use errors.As in older Go versions easily, so let's just
	// do a type assertion chain.
	switch e := err.(type) {
	case *UnterminatedCommentError:
		if t, ok := target.(**UnterminatedCommentError); ok {
			*t = e
			return true
		}
	}
	return false
}

// errorReader returns a fixed error on Read.
type errorReader struct {
	err error
}

func (r *errorReader) Read(p []byte) (int, error) {
	return 0, r.err
}

// errorWriter always fails writes.
type errorWriter struct{}

func (w *errorWriter) Write(p []byte) (int, error) {
	return 0, io.ErrShortWrite
}
