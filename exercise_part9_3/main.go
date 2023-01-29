package main

import (
	"fmt"
)

/*

Удалите изз программы simple.go оператор time.Sleep(1 * time.Second)
Посмотрите, что произойдет. Почему так происходит?

*/

func function() {
	for i := 0; i < 10; i++ {
		fmt.Print(i)
	}
}

func main() {
	go function()

	go func() {
		for i := 10; i < 20; i++ {
			fmt.Print(i, " ")
		}
	}()

	fmt.Println()
}

/*
	При удалении time.Sleep(1 * time.Second) горутинам не хватает времени, чтобы  сделать свою  работу
*/
