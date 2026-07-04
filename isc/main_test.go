package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// iscPath is the path to the built isc binary.
var iscPath string

func TestMain(m *testing.M) {
	// Build the binary once
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	iscPath = filepath.Join(dir, "isc-test-binary")
	cmd := exec.Command("go", "build", "-o", iscPath, ".")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic("failed to build isc binary: " + err.Error())
	}
	defer os.Remove(iscPath)

	os.Exit(m.Run())
}

func TestStrip_FileInput(t *testing.T) {
	// Use the simple test fixture
	out, err := runStrip("testdata/simple.is")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should have no comment lines
	if strings.Contains(out, "//") {
		t.Fatal("output contains '//'")
	}
	if strings.Contains(out, "/*") {
		t.Fatal("output contains '/*'")
	}
	// Should preserve code
	if !strings.Contains(out, `state greeting = "Hello, World!"`) {
		t.Fatal("output missing expected code")
	}
}

func TestStrip_StdinInput(t *testing.T) {
	input := `// comment
state x = 1`
	out, err := runStripStdin(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Comment-only line collapsed — no leading blank line
	expected := "state x = 1"
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestStrip_OutputFlag(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "output.is")
	_, err := runStrip("-o", tmpFile, "testdata/simple.is")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}
	if strings.Contains(string(data), "//") {
		t.Fatal("output file contains '//'")
	}
}

func TestStrip_NonexistentFile(t *testing.T) {
	_, err := runStrip("nonexistent.is")
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}
}

func TestStrip_UnterminatedFile(t *testing.T) {
	_, err := runStrip("testdata/unterminated.is")
	if err == nil {
		t.Fatal("expected error for unterminated comment, got nil")
	}
	if !strings.Contains(err.Error(), "unterminated block comment") {
		t.Fatalf("expected 'unterminated block comment' in error, got: %v", err)
	}
}

func TestStrip_EmptyFile(t *testing.T) {
	out, err := runStrip("testdata/empty.is")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Fatalf("expected empty output, got %q", out)
	}
}

func TestStrip_StringsPreserved(t *testing.T) {
	out, err := runStrip("testdata/strings.is")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"http://example.com"`) {
		t.Fatal("string with // was corrupted")
	}
	if !strings.Contains(out, `"This is a /* not a comment */ string"`) {
		t.Fatal("string with /* was corrupted")
	}
	if !strings.Contains(out, `"he said \"hello\""`) {
		t.Fatal("string with escaped quotes was corrupted")
	}
}

func TestStrip_NoCommandError(t *testing.T) {
	cmd := exec.Command(iscPath)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error when no command given")
	}
	if !strings.Contains(string(out), "Usage") {
		t.Fatalf("expected usage message, got: %s", out)
	}
}

func TestStrip_UnknownCommandError(t *testing.T) {
	cmd := exec.Command(iscPath, "unknown")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
	if !strings.Contains(string(out), "unknown") {
		t.Fatalf("expected 'unknown' in error, got: %s", out)
	}
}

// -- helpers --

func runStrip(args ...string) (string, error) {
	cmdArgs := append([]string{"strip"}, args...)
	cmd := exec.Command(iscPath, cmdArgs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.String(), &stripError{msg: stderr.String()}
	}
	return stdout.String(), nil
}

func runStripStdin(input string) (string, error) {
	cmd := exec.Command(iscPath, "strip")
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.String(), &stripError{msg: stderr.String()}
	}
	return stdout.String(), nil
}

type stripError struct {
	msg string
}

func (e *stripError) Error() string {
	return e.msg
}
