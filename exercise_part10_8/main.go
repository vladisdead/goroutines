package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

/*
 	Попробуйте написать конкурентную версию утилиты командной строки find(1),
	в которой бы применялась управляющая горутина.
*/

var writeValue = make(chan string)

func set(dir string) {
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
			writeValue <- fmt.Sprintf("%s/%s\n", dir, fi.Name())
			if fi.IsDir() {
				set(dir + "/" + fi.Name())
			}
		}
	}
}

func monitor() {
	for {
		select {
		case newValue := <-writeValue:
			fmt.Print(newValue)
		}
	}
}

func main() {
	var w sync.WaitGroup

	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalf("Usage: %s <dir>", os.Args[0])
	}

	dir := flag.Arg(0)

	go monitor()

	w.Add(1)
	go func() {
		defer w.Done()
		set(dir)
	}()

	w.Wait()
}
