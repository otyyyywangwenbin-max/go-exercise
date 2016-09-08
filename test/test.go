package main

import (
	"fmt"
	"reflect"
	"time"
	"unsafe"

	_ "github.com/otyyyywangwenbin/go-exercise/test/pkg1" // for invoke pkg1/sub1.init()
	"github.com/otyyyywangwenbin/go-exercise/test/pkg2"   // invoke all pkg2.init()
	"golang.org/x/tools/go/loader"
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

func testDynamicLoader() {
	fmt.Println("testDynamicLoader")
	var conf loader.Config
	conf.Import("github.com/otyyyywangwenbin/go-exercise/test/dynamicpkg1")
	if prog, err := conf.Load(); err != nil {
		fmt.Println(err)
	} else {
		//fmt.Println(prog)
		for _, pkgInfo := range prog.InitialPackages() {
			// reflect.
			// vc := reflect.New(pkgInfo.Pkg.Name)
			fmt.Println("-------", pkgInfo.Pkg.Name())
			for key, val := range pkgInfo.Types {
				fmt.Println("-------> key: ", key, ", value: ", val)
			}
			fmt.Println("============")
			fmt.Println(reflect.TypeOf(testLoopChan))
			for key, val := range pkgInfo.Defs {
				fmt.Println("-------> key: ", key, ", value: ", reflect.TypeOf(val))

				// if key.Name == "init" {
				// 	if obj, ok := val.(*types.Func); ok {
				// 		reflect.ValueOf(*obj).Call(make([]reflect.Value, 0))
				// 		fmt.Println(obj.Pos())
				// 		fmt.Println("xxxxxxxx:", obj)
				// 	}
				// }
			}
		}

	}
}

func main() {
	fmt.Println("begin main")
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

	//testLoopChan()

	pkg2.Method1()
	pkg2.Method2()

	testDynamicLoader()

}
