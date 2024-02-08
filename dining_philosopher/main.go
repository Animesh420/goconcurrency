package main

import (
	"fmt"
	"sync"
	"time"
)

// Dining philosopher problem
// In a table there sit some philosophers and each eat a spagethi
// Each plate has 1 chop stick next to eat and each philosopher can eat using two chopsticks, so no two neighbors can eat
// Implementing Djikstra's solution

// struct to store info about a philosopher
type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

// list of philosophers
var philosophers = []Philosopher{
	{name: "Plato", leftFork: 4, rightFork: 0},
	{name: "Socrates", leftFork: 0, rightFork: 1},
	{name: "Aristotle", leftFork: 1, rightFork: 2},
	{name: "Pascal", leftFork: 2, rightFork: 3},
	{name: "Locke", leftFork: 3, rightFork: 4},
}

var orderMutex sync.Mutex
var orderFinished = []string{}

// define some variables
var hungry = 3                  // how many times does a person eat
var eatTime = 1 * time.Second   // time duration to eat for 1 philosopher
var thinkTime = 3 * time.Second // time duration to think for 1 philosopher
var sleepTime = 1 * time.Second // time duration to log to console for 1 philosopher

func dine() {
	// eatTime = 0 * time.Second
	// thinkTime = 0 * time.Second
	// sleepTime = 0 * time.Second

	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	// forks is a map of all 5 forks
	var forks = make(map[int]*sync.Mutex)
	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	// start the meal
	for i := 0; i < len(philosophers); i++ {
		// fire of a go routine for the current philosopher
		go diningProblem(philosophers[i], wg, forks, seated)
	}

	wg.Wait()

}

func diningProblem(philosopher Philosopher, wg *sync.WaitGroup, forks map[int]*sync.Mutex, seated *sync.WaitGroup) {
	defer wg.Done()

	// seat the philosopher at the table
	fmt.Printf("%s is seated at the table.\n", philosopher.name)
	seated.Done()

	seated.Wait()

	// eat three times
	for i := hungry; i > 0; i-- {

		if philosopher.leftFork > philosopher.rightFork {
			// get a lock on both forks (logical race condition)
			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s takes the right fork %d.\n", philosopher.name, philosopher.rightFork)
			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s takes the left fork %d.\n", philosopher.name, philosopher.leftFork)
		} else {
			// get a lock on both forks
			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s takes the left fork %d.\n", philosopher.name, philosopher.leftFork)
			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s takes the right fork %d.\n", philosopher.name, philosopher.rightFork)

		}

		fmt.Printf("\t%s has both forks and is eating.\n", philosopher.name)
		time.Sleep(eatTime)

		fmt.Printf("\t%s is thinking.\n", philosopher.name)
		time.Sleep(thinkTime)

		fmt.Printf("\t%s leaves the left fork %d.\n", philosopher.name, philosopher.leftFork)
		forks[philosopher.leftFork].Unlock()

		fmt.Printf("\t%s leaves the right fork %d.\n", philosopher.name, philosopher.rightFork)
		forks[philosopher.rightFork].Unlock()

		fmt.Printf("\t%s has put down the forks.\n", philosopher.name)

	}
	orderMutex.Lock()
	orderFinished = append(orderFinished, philosopher.name)
	fmt.Println(philosopher.name, " is satisifed")
	fmt.Println(philosopher.name, " has left the table")
	orderMutex.Unlock()

}

func main() {
	// print out a welcome message
	fmt.Println("Dining philosophers problem BEGINNING")
	fmt.Println("The table is empty")
	fmt.Println("-------------------------------------")
	// Sleep a little
	time.Sleep(time.Duration(sleepTime.Milliseconds()))
	// start the meal
	dine()

	// print out the finished message
	fmt.Println("Dining philosophers problem ENDING")
	fmt.Println("The table is empty")
	fmt.Println("-------------------------------------")

	for i, name := range orderFinished {
		fmt.Printf("Ate at position %d -> philosopher -> %s\n", i, name)
	}
}
