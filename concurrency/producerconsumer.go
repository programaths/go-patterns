package main


import (
	"fmt"
)

var done = make(chan bool)
var msgs = make(chan int)

// produce sends integers to the msgs channel
// and closes the channel upon completion.
func produce() {
	defer close(msgs)
	for i := 0; i < 10; i++ {
		msgs <- i
	}
}

// consume will pull msgs off of the channel
// for as long as msgs are being procuded.
// When the msgs channel is closed by produce(),
// the for statement will terminate, and consume
// will signal on the done channel.
func consume() {
	for msg := range msgs {
		fmt.Println(msg)
	}
	done <- true
}

// After creating the produce and consume 
// goroutines, the main function will
// wait until the done channel is signalled
// (by the consume function)
func main() {
	go produce()
	go consume()
	<-done
}
