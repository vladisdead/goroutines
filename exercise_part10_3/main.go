package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

/*	попробуйте реализовать конкурентную версию wc(1), в которой бы
	использовалась управляющая горутина
*/

type FileData struct {
	lines      int
	words      int
	characters int
	filename   string
}

//var readData = make(chan string)  // канал для чтение результата после открытия
var writeData = make(chan string) // канал для записи нового результата после открытия новго файла

// вычисляет нужные нам данные и устанавливает результат в канал writeData
func wc(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	lines, words, characters := 0, 0, 0
	for scanner.Scan() {
		lines++

		line := scanner.Text()
		characters += len(line)

		splitLines := strings.Split(line, " ")
		words += len(splitLines)
	}
	writeData <- fmt.Sprintf("%8d%8d%8d %s", lines, words, characters, filename)
}

/*	в этой функции заключена логика управляющей горутины
	когда поступает запрос на чтение, функция read() пытается выполнить операцию чтения из канала readData, который управляется функцией monitor().
	Результатом операции является текущее значение, которое хранится в переменной
	output. И наоборот, когда мы хотим изменить сохраненное значение, то вызываем
	функцию wc(). Она записывает данные в канал writeData, который также обрабатывается оператором select. В результате никто не может обратиться к общей
	переменной в обход функции monitor(). */

func monitor() {
	var output string
	for {
		select {
		case newOutput := <-writeData:
			output = newOutput
			fmt.Println(output)
		}
	}
}

func main() {
	var fileNames []string

	if len(os.Args) <= 1 {
		log.Fatal("no input files")
	}

	fileNames = make([]string, 0)
	fileNames = append(fileNames, os.Args[1:]...)

	fileCount := len(fileNames)
	go monitor()

	var w sync.WaitGroup

	for i := 0; i < fileCount; i++ {
		w.Add(1)
		go func(i int) {
			defer w.Done()
			wc(fileNames[i])
		}(i)
	}
	w.Wait()
}
