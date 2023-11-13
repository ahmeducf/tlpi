package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func createPipe(t *testing.T) (*os.File, *os.File, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	t.Cleanup(func() {
		r.Close()
		w.Close()
	})

	return r, w, err
}

func TestWriteToStdout(t *testing.T) {
	originalStdin := os.Stdin
	originalStdout := os.Stdout

	rStdout, wStdout, err := createPipe(t)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	os.Stdout = wStdout

	rStdin, wStdin, err := createPipe(t)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	os.Stdin = rStdin

	wStdin.Write([]byte("Hello, World!\n"))
	wStdin.Close()

	rootCmd.Execute()
	wStdout.Close()

	os.Stdout = originalStdout
	os.Stdin = originalStdin

	var buf bytes.Buffer
	io.Copy(&buf, rStdout)

	actual := buf.String()
	expected := "Hello, World!\n"

	if actual != expected {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestWriteToFile(t *testing.T) {
	originalStdin := os.Stdin
	originalStdout := os.Stdout
	_, wStdout, err := createPipe(t)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	os.Stdout = wStdout

	rStdin, wStdin, err := createPipe(t)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	os.Stdin = rStdin

	wStdin.Write([]byte("Hello, World!\n"))
	wStdin.Close()

	rootCmd.SetArgs([]string{"/tmp/test.txt"})
	rootCmd.Execute()

	os.Stdin = originalStdin
	os.Stdout = originalStdout

	f, err := os.Open("/tmp/test.txt")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer func() {
		f.Close()
		os.Remove("/tmp/test.txt")
	}()

	var buf2 bytes.Buffer
	io.Copy(&buf2, f)

	actual := buf2.String()
	expected := "Hello, World!\n"

	if actual != expected {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestWriteToMultipleFiles(t *testing.T) {
	originalStdin := os.Stdin
	originalStdout := os.Stdout
	_, wStdout, err := createPipe(t)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	os.Stdout = wStdout

	rStdin, wStdin, err := createPipe(t)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	os.Stdin = rStdin

	wStdin.Write([]byte("Hello, World!\n"))
	wStdin.Close()

	rootCmd.SetArgs([]string{"/tmp/test.txt", "/tmp/test2.txt"})
	rootCmd.Execute()

	os.Stdin = originalStdin
	os.Stdout = originalStdout

	f, err := os.Open("/tmp/test.txt")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer func() {
		f.Close()
		os.Remove("/tmp/test.txt")
	}()

	var buf bytes.Buffer
	io.Copy(&buf, f)

	actual := buf.String()
	expected := "Hello, World!\n"

	if actual != expected {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}

	f2, err := os.Open("/tmp/test2.txt")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer func() {
		f2.Close()
		os.Remove("/tmp/test2.txt")
	}()

	var buf2 bytes.Buffer
	io.Copy(&buf2, f2)

	actual2 := buf2.String()
	expected2 := "Hello, World!\n"

	if actual2 != expected2 {
		t.Errorf("expected: %v, actual: %v", expected2, actual2)
	}
}

func TestWriteToStdoutAndFile(t *testing.T) {
	originalStdin := os.Stdin
	originalStdout := os.Stdout

	rStdout, wStdout, err := createPipe(t)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	os.Stdout = wStdout

	rStdin, wStdin, err := createPipe(t)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	os.Stdin = rStdin

	wStdin.Write([]byte("Hello, World!\n"))
	wStdin.Close()

	rootCmd.SetArgs([]string{"/tmp/test.txt"})
	rootCmd.Execute()
	wStdout.Close()

	os.Stdout = originalStdout
	os.Stdin = originalStdin

	f, err := os.Open("/tmp/test.txt")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer func() {
		f.Close()
		os.Remove("/tmp/test.txt")
	}()

	var buf bytes.Buffer
	io.Copy(&buf, f)

	actual := buf.String()
	expected := "Hello, World!\n"

	if actual != expected {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}

	var buf2 bytes.Buffer
	io.Copy(&buf2, rStdout)

	actual2 := buf2.String()
	expected2 := "Hello, World!\n"

	if actual2 != expected2 {
		t.Errorf("expected: %v, actual: %v", expected2, actual2)
	}
}

func TestAppendToFile(t *testing.T) {
	f, err := os.Create("/tmp/test.txt")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	n, err := f.Write([]byte("Hello\n"))
	if len([]byte("Hello\n")) != n || err != nil {
		t.Fatalf("error: %v", err)
	}

	f.Close()

	originalStdin := os.Stdin
	originalStdout := os.Stdout
	_, wStdout, err := createPipe(t)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	os.Stdout = wStdout

	rStdin, wStdin, err := createPipe(t)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	os.Stdin = rStdin

	wStdin.Write([]byte("World\n"))
	wStdin.Close()

	rootCmd.SetArgs([]string{"/tmp/test.txt", "-a"})
	rootCmd.Execute()

	os.Stdin = originalStdin
	os.Stdout = originalStdout

	f, err = os.Open("/tmp/test.txt")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, f)

	actual := buf.String()
	expected := "Hello\nWorld\n"

	if actual != expected {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func teeIteration(b *testing.B, args []string) {
	originalStdin := os.Stdin
	rStdout, wStdout, err := os.Pipe()
	if err != nil {
		b.Fatalf("error: %v", err)
	}
	defer func() {
		rStdout.Close()
		wStdout.Close()
		os.Stdout = originalStdin
	}()
	os.Stdout = wStdout

	rStdin, wStdin, err := os.Pipe()
	if err != nil {
		b.Fatalf("error: %v", err)
	}
	os.Stdin = rStdin

	wStdin.Write([]byte("Hello, World!\n"))
	wStdin.Close()

	for i := 0; i < b.N; i++ {
		rootCmd.SetArgs(args)
		rootCmd.Execute()
	}

	rStdin.Close()
	os.Stdin = originalStdin
}

func generateFileNames(n int) []string {
	var args []string
	for i := 0; i < n; i++ {
		args = append(args, "/tmp/test"+fmt.Sprint(i)+".txt")
	}
	return args
}

func BenchmarkTee(b *testing.B) {
	for _, v := range []int{1, 10, 100, 1000} {
		args := generateFileNames(v)
		b.Run(fmt.Sprintf("Tee-%d", v), func(b *testing.B) {
			teeIteration(b, args)
		})

		err := os.RemoveAll("/tmp/test*.txt")
		if err != nil {
			b.Fatalf("error: %v", err)
		}
	}
}
