package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func updateMessage(s string) {
	defer wg.Done()
	msg = s
}

func printMessage() {
	fmt.Println(msg)

}

func main() {

	msg = "hello world"

	wg.Add(1)
	go updateMessage("hello universe")
	wg.Wait()
	printMessage()

	wg.Add(1)
	go updateMessage("hello cosmos")
	wg.Wait()
	printMessage()

	wg.Add(1)
	go updateMessage("hello world")
	wg.Wait()
	printMessage()
}
