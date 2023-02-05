package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

/*	Измените код Go файла workerPool.go таким образом, чтобы сохранять результаты в файле.
	При работе с файлом примените мьютекс и критический раздел
	или управляющую горутину, которая будет записывать данные на диск


	что произойдет с программой workerPool.go, если значение глобальной переменной size станет равным 1? Почему?

	Если изменить size на 1. То ClientID будут выводиться не по порядку
	Это происходит потому что число одновременно выполняющихся горутин = 1

	Измените код Go программы workerPool.go таким образом, чтобы размер буферизованных каналов clients и data можно было задавать в виде аргументов
	командной строки.

*/

type Client struct {
	id      int
	integer int
}

type Data struct {
	job    Client
	square int
}

var (
	size = 1

	clients chan Client
	data    chan Data
)

func worker(w *sync.WaitGroup) {
	for c := range clients {
		square := c.integer * c.integer
		output := Data{c, square}
		data <- output
		time.Sleep(time.Second)
	}
	w.Done()
}

func makeWP(n int) {
	var w sync.WaitGroup
	for i := 0; i < n; i++ {
		w.Add(1)
		go worker(&w)
	}
	w.Wait()
	close(data)
}

func create(n int) {
	for i := 0; i < n; i++ {
		c := Client{i, i}
		clients <- c
	}
	close(clients)
}

func writeResultToFile(file *os.File, result string) {
	m.Lock()
	_, err := file.WriteString(result + "\n")
	if err != nil {
		return
	}
	m.Unlock()
}

var (
	m    sync.Mutex
	file *os.File
)

func main() {

	countClients := flag.Int("clients", 10, "Number of clients")
	countData := flag.Int("data", 10, "Number of data")

	flag.Parse()

	clients = make(chan Client, *countClients)
	data = make(chan Data, *countData)

	fmt.Println("Capacity of clients:", cap(clients))
	fmt.Println("Capacity of data:", cap(data))

	if len(os.Args) < 3 {
		fmt.Println("Need #jobs and  #workers")
		os.Exit(1)
	}

	var nJobs int
	var nWorkers int

	if len(os.Args) == 3 {
		Jobs, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
			return
		}

		Workers, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		nJobs = Jobs
		nWorkers = Workers
	} else if len(os.Args) == 7 {
		Jobs, err := strconv.Atoi(os.Args[5])
		if err != nil {
			log.Fatal(err)
			return
		}

		Workers, err := strconv.Atoi(os.Args[6])
		if err != nil {
			fmt.Println(err)
			return
		}
		nJobs = Jobs
		nWorkers = Workers
	} else {
		log.Fatal("usage -data 10 -clients 10 10 10")
	}

	file, err := os.Create("result.txt")
	if err != nil {
		log.Fatal(err)
		return
	}

	go create(nJobs)
	finished := make(chan interface{})
	go func() {
		for d := range data {
			writeResultToFile(file, fmt.Sprintf("Client ID: %d\tint: %d\tsquare: %d", d.job.id, d.job.integer, d.square))
			fmt.Printf("Client ID: %d\tint: ", d.job.id)
			fmt.Printf("%d\tsquare: %d\n", d.job.integer, d.square)
		}
		finished <- true
	}()
	makeWP(nWorkers)
	fmt.Printf(":  %v\n", <-finished)
	err = file.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
}
