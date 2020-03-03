// Use `go run foo.go` to run your program

package main

import (
    . "fmt"
    "runtime"
)

// Control signals
const (
	GetNumber = iota
	Exit
)

func number_server(add_number <-chan int, control <-chan int, number chan<- int) {
	var i = 0

	// This for-select pattern is one you will become familiar with if you're using go "correctly".
	for {
		select {
			// TODO: receive different messages and handle them correctly
			// You will at least need to update the number and handle control signals.
    case msg1 := <-add_number: //Hvis kanalen add_number får inn noe, sendes dette til msg1
      i += msg1 //Plusser msg1 inn på i, msg1 er enten 1 eller -1
    case msg2 :=<- control: //Sender hva vi får på control inn på msg2, kan enten få inn GetNumber eller Exit
      switch msg2 {
      case GetNumber: //Hvis msg2 er lik GetNumber, sendes i inn på number
            number <- i
          case Exit: //Ferdig
            break
      }
    		}
	}
}

func incrementing(add_number chan<-int, finished chan<- bool) {
	for j := 0; j<1000000; j++ {
		add_number <- 1
	}
	finished <- true
	//TODO: signal that the goroutine is finished
}

func decrementing(add_number chan<- int, finished chan<- bool) {
	for j := 0; j<1000000; j++ {
		add_number <- -1
	}
	finished <- true
	//TODO: signal that the goroutine is finished
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// TODO: Construct the required channels
	// Think about wether the receptions of the number should be unbuffered, or buffered with a fixed queue size.

	control := make(chan int)
	number := make(chan int)
	add_number := make(chan int)
	finished := make(chan bool)

	// TODO: Spawn the required goroutines
	go incrementing(add_number, finished)
	go decrementing(add_number, finished)
	go number_server(add_number,control,number)

	// TODO: block on finished from both "worker" goroutines

	<-finished
	<-finished
	//når vi er her, er incrementing og descrementing ferdig


	control<-GetNumber
	Println("The magic number is:", <- number)
	control<-Exit
}
