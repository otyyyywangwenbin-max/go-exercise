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
	//c <- &aT{}
	fmt.Println(unsafe.Sizeof(c))
	fmt.Printf("----> %#v\n", c)
	fmt.Println("xxxxxxxx", <-c)
}

func test2(c <-chan *aT) {
	//c <- &aT{}
	fmt.Printf("----> %#v\n", c)
	fmt.Println("xxxxxxxx", <-c)
}

func testForever() {
	run := func() {
		fmt.Println("------------->>>")
		select {}
	}

	run()
	fmt.Println("|||||||||||")
}

func main() {
	/*
		fmt.Println(runtime.NumCPU())
		fmt.Println(runtime.GOMAXPROCS())
		fmt.Println(runtime.GOMAXPROCS(1))
		fmt.Println(runtime.GOMAXPROCS(2))
	*/
	//testClosure()
	go testForever()
	fmt.Println("begin")
	time.Sleep(5 * time.Second)
	fmt.Println("end")

}
