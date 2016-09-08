package dynamicpkg1

import (
	"fmt"

	_ "github.com/otyyyywangwenbin/go-exercise/test/dynamicpkg1/sub1"
)

func init() {
	fmt.Println("----------------dynamicpkg1")
}
