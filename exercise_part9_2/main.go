package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

/*Создайте конвейер для вычисления суммы квадратов всех натуральных чисел
в заданном диапазоне.*/

func WriteToChannel(n1, n2 float64, out chan<- float64) {
	for i := n1; i <= n2; i++ {
		out <- i
	}

	close(out)
}

func CalculateSqrt(out chan<- float64, in <-chan float64) {
	for i := range in {
		out <- math.Pow(i, 2)
	}
	close(out)
}

func SumSqrt(out <-chan float64) {
	sum := 0.0

	for square := range out {
		sum = sum + square
	}

	fmt.Println("Sum of square ", sum)
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

	if n1 <= 0 || n2 <= 0 {
		fmt.Println("Enter natural numbers")
		return
	}

	A := make(chan float64)
	B := make(chan float64)

	fmt.Printf("From %d to %d\n", n1, n2)

	go WriteToChannel(float64(n1), float64(n2), A)
	go CalculateSqrt(B, A)

	SumSqrt(B)
}
