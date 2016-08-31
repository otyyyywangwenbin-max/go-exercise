package main

import "fmt"

type aT struct {
	A string
}

func testClosure() {
	var c = make(chan *aT, 100)

	c <- &aT{A: "A1"}
	c <- &aT{A: "A2"}
	c <- &aT{A: "A3"}

	// for {
	// 	v := <-c
	// 	go func() {
	// 		fmt.Printf("----> %x", v)
	// 		v.A = v.A + "-B"
	// 		fmt.Println(*v)
	// 	}()
	// }

	for v := range c {
		fmt.Printf("----> %x", v)
		v := v
		go func() {
			fmt.Printf("----> %x", v)
			v.A = v.A + "-B"
			fmt.Println(*v)
		}()
	}
}

func main() {
	testClosure()
}
