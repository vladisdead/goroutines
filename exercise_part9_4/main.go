package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

/*Измените код Go в файле pipeline.go таким образом, чтобы получить конвейер
из пяти функций и соответствующего количества каналов.*/

var CLOSEA = false

var DATA = make(map[int]bool)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func first(min, max int, out chan<- int) {
	for {
		if CLOSEA {
			close(out)
			return
		}
		out <- random(min, max)
	}
}

func second(out chan<- int, in <-chan int) {
	for x := range in {
		out <- x
	}
	close(out)
}

func third(out chan<- int, in <-chan int) {
	for x := range in {
		out <- x
	}
	close(out)
}

func fourth(out chan<- int, in <-chan int) {
	for x := range in {
		fmt.Print(x, " ")
		_, ok := DATA[x]
		if ok {
			CLOSEA = true
		} else {
			DATA[x] = true
			out <- x
		}
	}
	fmt.Println()
	close(out)
}

func fifth(in <-chan int) {
	var sum int
	sum = 0
	for x2 := range in {
		sum = sum + x2
	}
	fmt.Printf("The sum of  the random numbers is %d\n", sum)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Need two  integer parameters!")
		return
	}

	n1, _ := strconv.Atoi(os.Args[1])
	n2, _ := strconv.Atoi(os.Args[2])

	if n1 > n2 {
		fmt.Printf("%d should be smaller  than %d\n", n1, n2)
		return
	}

	rand.Seed(time.Now().UnixNano())

	A := make(chan int)
	B := make(chan int)
	C := make(chan int)
	D := make(chan int)
	go first(n1, n2, A)
	go second(B, A)
	go third(C, B)
	go fourth(D, C)
	fifth(D)
}
