package main

import (
	"fmt"
	"time"
)

func listenToChan(ch chan int) {
	for {
		// print a got data message
		i := <-ch
		time.Sleep(1 * time.Second)
		fmt.Println("Got", i, "from a channel")

	}
}

func main() {

	ch := make(chan int, 5)

	go listenToChan(ch)

	s := time.Now()

	for i := 0; i <= 20; i++ {
		fmt.Println("Sending ", i, " to channel")
		ch <- i
		fmt.Println("Sent ", i, " to channel")
	}

	fmt.Println("Done between ", time.Now(), s)
	close(ch)
}
