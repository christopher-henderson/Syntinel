package main

import (
	"fmt"
	"reflect"
)

func main() {
	// channel := make(chan int, 5)
	// queue := structures.NewQueue()
	// channel <- 5
	// queue.Push(channel)
	// fmt.Println(reflect.TypeOf(channel))
	// c := queue.Pop().(chan int)
	// fmt.Println(<-c)
	// c <- 6
	// fmt.Println(<-c)

	// tick := time.Tick(100 * time.Millisecond)
	// boom := time.After(500 * time.Millisecond)
	// for {
	// 	select {
	// 	case <-tick:
	// 		fmt.Println("tick.")
	// 	case <-boom:
	// 		fmt.Println("BOOM!")
	// 		// return
	// 	default:
	// 		fmt.Println("    .")
	// 		time.Sleep(50 * time.Millisecond)
	// 	}
	// }
	fmt.Println(reflect.TypeOf(lol))
}

func lol() {

}
