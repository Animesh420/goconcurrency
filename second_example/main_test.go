package main

import "testing"

func Test_updateMessage(t *testing.T) {
	msg = "Hello world!"
	wg.Add(2)
	go updateMessage("Goodby world")
	go updateMessage("Goodbye Alucard")
	wg.Wait()

	if msg != "Goodby world" {
		t.Error("incorrect value in msg")
	}
}
