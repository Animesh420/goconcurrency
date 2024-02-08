package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_printSomething(t *testing.T) {
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	var wg sync.WaitGroup

	wg.Add(1)

	go printSomething(&wg, "epsilon")
	wg.Wait()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdout

	if !strings.Contains(output, "epsilons") {
		t.Errorf("Expected to find epsilon, but it is not there found %s", output)
	}
}
