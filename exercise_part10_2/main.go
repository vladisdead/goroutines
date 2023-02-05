package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

/*	Попробуйте реализовать конкурентную версию wc(1),
	которая бы использовала общую память.
*/

// эти глобальные переменные и есть общая память
var (
	m          sync.Mutex
	lines      int
	characters int
	words      int
)

// функция которая обрабатывает файл.
func wc(filename string) {
	m.Lock()
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	lines, words, characters = 0, 0, 0
	for scanner.Scan() {
		lines++

		line := scanner.Text()
		characters += len(line)

		splitLines := strings.Split(line, " ")
		words += len(splitLines)
	}

	m.Unlock()
}

// функция, которая пишет в консоль
func writeToConsole(filename string) {
	m.Lock()
	alines := lines
	awords := words
	acharacters := characters
	m.Unlock()
	fmt.Printf("%8d%8d%8d %s\n", alines, awords, acharacters, filename)
}

// wc() и writeToConsole(). Все дейсвтия в них происходят  между Lock() и Unlock(),
// которые используют один mutex.Когда мьютекс блокирован, это означает, что никто другой не может заблокировать этот мьютекс, пока он не будет
// освобожден с помощью функции sync.Unlock().
func main() {
	var fileNames []string

	if len(os.Args) <= 1 {
		log.Fatal("no input files")
	}

	fileNames = make([]string, 0)
	fileNames = append(fileNames, os.Args[1:]...)

	fileCount := len(fileNames)

	var waitGroup sync.WaitGroup
	log.Println("Files count", fileCount)

	for i := 0; i < fileCount; i++ {
		waitGroup.Add(1)
		go func(filename string) {
			defer waitGroup.Done()
			wc(filename)
			writeToConsole(filename)
		}(fileNames[i])
	}
	waitGroup.Wait()
}
