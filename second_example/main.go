package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func updateMessage(s string) {
	// race condition when updated by multiple
	defer wg.Done()
	msg = s
}

func updateMessageMutex(s string, m *sync.Mutex) {
	defer wg.Done()

	m.Lock()
	msg = s
	m.Unlock()
}

func main() {
	msg = "Hello world"
	var mutex sync.Mutex

	// Example of race condition when both goroutines try to update the same variable
	// wg.Add(2)
	// go updateMessage("Hello universe")
	// go updateMessage("Hello cosmos")
	// wg.Wait()

	// fmt.Println(msg)

	wg.Add(2)
	go updateMessageMutex("Hello universe", &mutex)
	go updateMessageMutex("Hello cosmos", &mutex)
	wg.Wait()

	fmt.Println(msg)

}
