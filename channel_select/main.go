package main

import (
	"fmt"
	"time"
)

func server1(ch chan string) {
	for {
		time.Sleep(6 * time.Second)
		ch <- "this is from server1"
	}

}

func server2(ch chan string) {
	for {
		time.Sleep(3 * time.Second)
		ch <- "this is from server2"
	}

}

func main() {
	fmt.Println("Select with channels")
	fmt.Println("--------------------")

	channel1, channel2 := make(chan string), make(chan string)

	go server1(channel1)
	go server2(channel2)

	for {
		select {
		case s1 := <-channel1:
			fmt.Println("Case1", s1)

		case s2 := <-channel1:
			fmt.Println("Case2", s2)

		case s3 := <-channel2:
			fmt.Println("Case3", s3)

		case s4 := <-channel2:
			fmt.Println("Case4", s4)
		}
	}

}
