package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

/*	Измените код simpleContext.go таким образом, чтобы анонимная функция,
	используемая в функциях f1(), f2() и f3(), стала отдельной функцией. Какая
	основная проблема возникает при таком изменении кода?

	Я немного не понял, какая  проблема  должна была появиться.:)


	 Измените код Go программы simpleContext.go таким образом, чтобы функции
	f1(), f2() и f3() использовали созданную извне переменную Context вместо
	определения собственной такой переменной.

*/

func anonFunction(cancel context.CancelFunc) {
	time.Sleep(4 * time.Second)
	cancel()

}

func f1(ctx context.Context, t int) {
	//c1 := context.Background()
	c1, cancel := context.WithCancel(ctx)
	defer cancel()

	go anonFunction(cancel)

	/*go func() {
		time.Sleep(4 * time.Second)
		cancel()
	}()*/

	select {
	case <-c1.Done():
		fmt.Println("f1():", c1.Err())
		return
	case r := <-time.After(time.Duration(t) * time.Second):
		fmt.Println("f1():", r)

	}
	return
}

func f2(ctx context.Context, t int) {
	//c2 := context.Background()
	c2, cancel := context.WithTimeout(ctx, time.Duration(t)*time.Second)
	defer cancel()

	go anonFunction(cancel)

	/*go func() {
		time.Sleep(4 * time.Second)
		cancel()
	}()*/

	select {
	case <-c2.Done():
		fmt.Println("f2():", c2.Err())
		return
	case r := <-time.After(time.Duration(t) * time.Second):
		fmt.Println("f2():", r)
	}
	return
}

func f3(ctx context.Context, t int) {
	//c3 := context.Background()
	deadline := time.Now().Add(time.Duration(2*t) * time.Second)
	c3, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()

	go anonFunction(cancel)

	/*go func() {
		time.Sleep(4 * time.Second)
		cancel()
	}()*/

	select {
	case <-c3.Done():
		fmt.Println("f3():", c3.Err())
		return
	case r := <-time.After(time.Duration(t) * time.Second):
		fmt.Println("f3():", r)
	}
	return
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Need a delay!")
		return
	}

	delay, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx := context.Background()

	fmt.Println("Delay:", delay)
	f1(ctx, delay)
	f2(ctx, delay)
	f3(ctx, delay)
}
