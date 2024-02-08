package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNuber int
	message    string
	success    bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func getNonZeroRand(till int) int {
	return rand.Intn(till) + 1
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= NumberOfPizzas {
		delay := getNonZeroRand(5)
		fmt.Printf("Received order #%d!\n", pizzaNumber)

		rnd := getNonZeroRand(12)
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed += 1
		} else {
			pizzasMade += 1
		}

		total++
		fmt.Printf("Making pizza #%d. It will take %d seconds ...\n", pizzaNumber, delay)
		// delay for a bit
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza #%d!\n", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making pizza #%d!\n", pizzaNumber)
		} else {
			msg = fmt.Sprintf("Pizza order #%d is ready!\n", pizzaNumber)
			success = true
		}

		p := PizzaOrder{
			pizzaNuber: pizzaNumber,
			message:    msg,
			success:    success,
		}
		return &p

	}

	return &PizzaOrder{
		pizzaNuber: pizzaNumber,
	}
}

func pizzeria(pizzaMaker *Producer) {
	// PRODUCER
	// keep track of which pizza we are making
	var i = 0

	// run forever or until we receive a quit notification
	// try to make pizzas
	for {
		// decision structure

		currentPizza := makePizza(i)

		if currentPizza != nil {
			i = currentPizza.pizzaNuber
			select {
			// try to make a pizza, sent something to the data channel
			case pizzaMaker.data <- *currentPizza:
			case quitChan := <-pizzaMaker.quit:
				close(pizzaMaker.data)
				close(quitChan)
				return

			}
		}

	}
}

func main() {
	// seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// print out a message
	color.Cyan("The Pizzeria is open for business!")
	color.Cyan("----------------------------------")

	// create a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// run the producer in background (as a go routine)
	go pizzeria(pizzaJob)

	// create and run the consumer
	for i := range pizzaJob.data {
		if i.pizzaNuber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("Order #%d is out for delivery!", i.pizzaNuber)
			} else {
				color.Red(i.message)
				color.Red("The customer is really mad!")
			}
		} else {
			color.Cyan("Done making pizzas..")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("Error closing pizza job channel", err)
			}
		}
	}

	// print out the ending message
	color.Cyan("--------------------")
	color.Cyan("Done for the day !!")
	color.Cyan("We made %d pizzas, but failed to make %d, with %d attempts in total.", pizzasMade, pizzasFailed, total)

	switch {
	case pizzasFailed > 9:
		color.Red("It was an awful day")
	case pizzasFailed >= 6:
		color.Red("It was not a very good day..")
	case pizzasFailed >= 4:
		color.Yellow("It was an okay day..")
	case pizzasFailed >= 2:
		color.Yellow("It was a pretty okay day.")
	default:
		color.Green("It was a great day")
	}

}
