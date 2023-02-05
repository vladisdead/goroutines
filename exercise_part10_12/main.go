package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

/*	Наконец, попробуйте реализовать конкурентную версию утилиты командной
	строки find(1) с помощью мьютекса sync.Mutex.
*/

var (
	m      sync.Mutex
	output string
)

func readDir(dir string) {

	d, err := os.Open(dir)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()
	for {

		fis, err := d.Readdir(10)

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not read dir names in %s: %s", dir, err.Error())
		}

		for _, fi := range fis {
			m.Lock()
			output = fmt.Sprintf("%s/%s", dir, fi.Name())
			addToRead()
			if fi.IsDir() {
				readDir(dir + "/" + fi.Name())
			}

		}

	}

}

var outputArr = make([]string, 0)

func addToRead() {

	outputArr = append(outputArr, output)
	m.Unlock()
}

func read() {
	for _, v := range outputArr {
		fmt.Println(v)
	}
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalf("Usage: %s <dir>", os.Args[0])
	}

	dir := flag.Arg(0)

	var wg sync.WaitGroup

	wg.Add(1)
	go func(dir string) {
		defer wg.Done()
		readDir(dir)

	}(dir)

	wg.Wait()
	read()
}
