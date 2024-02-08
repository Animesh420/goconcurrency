package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

// - if there are no customers, the barber falls asleep in the chair
// - a customer must wake the barber if he is asleep
// - if a customer arrives while the barber is working, the customer leaves if all chairs are occupied and sits in an empty chair if it is available
// - when the barber finishes a haircut, he inspects the waiting room to see if there are any waiting customers and falls asleep  if thhere are none
// - shop can stop accepting new clients at closing time, but the barbers cannot leave until the waiting room is empty
// - after the shop is closed and there are no clients left in the waiting area, the barber goes home

// variables
var seatingCapacity = 10

var arrivalRate = 100 // in milliseconds
var cutDuration = 1000 * time.Millisecond

var timeOpen = 10 * time.Second

func main() {

	// seed our random number generator
	rand.Seed(time.Now().UnixNano())

	// print welcome message
	color.Yellow("The Sleeping Barber Problem")
	color.Yellow("---------------------------")

	// create channels if we need any
	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	// create the barbershop
	shop := BarberShop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		ClientsChan:     clientChan,
		BarbersDoneChan: doneChan,
		Open:            true,
	}

	color.Green("The shop is open for the day!!")

	// add barbers
	shop.addBarber("Frank")
	shop.addBarber("Gerard")
	shop.addBarber("Milton")
	shop.addBarber("Susan")
	shop.addBarber("Kelly")
	shop.addBarber("Pal")
	shop.addBarber("John")

	// start the barber shop as a goroutine
	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func() {
		<-time.After(timeOpen)
		shopClosing <- true
		shop.closeShopForDay()
		closed <- true
	}()

	// add clients
	i := 1
	go func() {
		for {
			// get a random number with average arrival rate
			randomMilliseconds := rand.Int() % (2 * arrivalRate)
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(randomMilliseconds)):
				shop.addClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}
	}()

	// block until the barber shop is closed
	// time.Sleep(5 * time.Second)

	<-closed
}
