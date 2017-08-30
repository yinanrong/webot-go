package main

import (
	"fmt"
	"sync"
)

var l sync.Mutex
var a string

func f(i int) {
	defer l.Unlock()
	a = fmt.Sprintf("hello, world:%d", i)
}

func main() {
	for i := 0; i <= 3; i++ {
		l.Lock()
		go f(i)

		println(a)
	}

}
