package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

/* Problem Statement:

The Dining Philosophers problem is well known in computer science circles.
Five philosophers, numbered from 0 through 4, live in a house where the
table is laid for them; each philosopher has their own place at the table.
Their only difficulty – besides those of philosophy – is that the dish
served is a very difficult kind of spaghetti which has to be eaten with
two forks. There are two forks next to each plate, so that presents no
difficulty. As a consequence, however, this means that no two neighbours
may be eating simultaneously, since there are five philosophers and five forks.

*/

// Philosopher struct holds information about philosopher and the forks that philosopher is holding.
type Philosopher struct {
	name      string
	leftFork  int
	rightFork int
}

// All Philosophers in the system
var philosophers = []Philosopher{
	{name: "P0", leftFork: 4, rightFork: 0},
	{name: "P1", leftFork: 0, rightFork: 1},
	{name: "P2", leftFork: 1, rightFork: 2},
	{name: "P3", leftFork: 2, rightFork: 3},
	{name: "P4", leftFork: 3, rightFork: 4},
}

const (
	timesEachPhilosopherEats = 3
	singleEatTime            = 1 * time.Second
	thinkTime                = 3 * time.Second
)

// To store the order in which philosophers completed eating.
var orderMutex sync.Mutex
var orderFinished []string

func main() {
	fmt.Println("\nDinning Philosopher Problem")
	fmt.Println("===========================")
	fmt.Println("The table is empty")
	startDine()
	fmt.Println("===========================")
	fmt.Println("The table is empty")
	fmt.Printf("Order finished: %s.\n", strings.Join(orderFinished, ", "))
}

func startDine() {
	philosophersEatingWG := &sync.WaitGroup{}
	philosophersEatingWG.Add(len(philosophers))

	philosophersSeated := &sync.WaitGroup{}
	philosophersSeated.Add(len(philosophers))

	var forks = make(map[int]*sync.Mutex)
	for forkCount := 0; forkCount < len(philosophers); forkCount++ {
		forks[forkCount] = &sync.Mutex{}
	}

	for philosopherCount := 0; philosopherCount < len(philosophers); philosopherCount++ {
		go diningProblem(philosophers[philosopherCount], philosophersEatingWG, forks, philosophersSeated)
	}

	philosophersEatingWG.Wait()
}

func diningProblem(philosopher Philosopher, philosophersEatingWG *sync.WaitGroup, forks map[int]*sync.Mutex,
	philosophersSeated *sync.WaitGroup) {
	defer philosophersEatingWG.Done()

	fmt.Printf("%s is seated at the table.\n", philosopher.name)
	// marking the philosopher as seated
	philosophersSeated.Done()

	// Waiting for all other philosopher to seat
	philosophersSeated.Wait()

	for timesToEat := timesEachPhilosopherEats; timesToEat > 0; timesToEat-- {

		if philosopher.leftFork > philosopher.rightFork {
			// This condition is for last philosopher who tries to pick the rightFork first,
			// but since that fork must be already picked by first philosopher, it will wait.
			// Otherwise, there are chances of a deadlock, every philosopher may have one fork
			// and wait for other forks to be released.
			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s takes the rightFork.\n", philosopher.name)

			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s takes the leftFork.\n", philosopher.name)
		} else {
			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s takes the leftFork.\n", philosopher.name)

			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s takes the rightFork.\n", philosopher.name)
		}
		fmt.Printf("\t%s has both forks and is eating.\n", philosopher.name)
		time.Sleep(singleEatTime)

		fmt.Printf("\t%s is thinking.\n", philosopher.name)
		time.Sleep(thinkTime)

		forks[philosopher.leftFork].Unlock()
		forks[philosopher.rightFork].Unlock()

		fmt.Printf("\t%s put down the forks.\n", philosopher.name)
	}
	fmt.Println(philosopher.name, "is satisfied.")
	fmt.Println(philosopher.name, "left the table.")
	
	// This condition is to acquire lock and add the philosopher name to a slice,
	// which can be used at the end of the program to print the order in which
	// philosophers completed eating.
	orderMutex.Lock()
	orderFinished = append(orderFinished, philosopher.name)
	orderMutex.Unlock()
}
