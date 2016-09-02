package main

import (
	"fmt"
	"time"
	"unsafe"
)

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

	// for v := range c {
	// 	fmt.Printf("----> %x", v)
	// 	v := v
	// 	go func() {
	// 		fmt.Printf("----> %x", v)
	// 		v.A = v.A + "-B"
	// 		fmt.Println(*v)
	// 	}()
	// }

	test(c)
	test2(c)
}

func test(c chan *aT) {

	fmt.Println(unsafe.Sizeof(c))
	fmt.Printf("----> %#v\n", c)
	fmt.Println("xxxxxxxx", <-c)
}

func test2(c <-chan *aT) {
	//c <- &aT{}
	fmt.Printf("----> %#v\n", c)
	fmt.Println("xxxxxxxx", <-c)
}

func testSelectBlocked() {
	run := func(_ <-chan struct{}) {
		fmt.Println("------------->>>")
		select {} //blocked
	}

	run(nil)
	fmt.Println("|||||||||||")
}

func testLoopChan() {
	c := make(chan int)
	run := func(c <-chan int) {
		for {
			v := <-c
			fmt.Println("value: ", v)
		}
	}

	go run(c)
	time.Sleep(10 * time.Second)
	c <- 1
	time.Sleep(1 * time.Second)
}

func main() {
	/*
		fmt.Println(runtime.NumCPU())
		fmt.Println(runtime.GOMAXPROCS())
		fmt.Println(runtime.GOMAXPROCS(1))
		fmt.Println(runtime.GOMAXPROCS(2))
	*/
	//testClosure()

	/*
		go testSelectBlocked()
		fmt.Println("begin")
		time.Sleep(5 * time.Second)
		fmt.Println("end")
	*/
	testLoopChan()

}
