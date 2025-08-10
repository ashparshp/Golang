package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func test_printMessage() {
	fmt.Println(msg)
}

func Test_updateMessage(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	var wg sync.WaitGroup
	wg.Add(1)

	msg := "Test Message"

	go updateMessage(msg, &wg)
	test_printMessage()
	wg.Wait()

	w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, "Hello, World!") {
		t.Errorf("Expected output to contain 'Hello, World!', got '%s'", output)
	}
}
