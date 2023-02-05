package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

/*	Попробуйте реализовать конкурентную версию wc(1), которая бы использовала
	буферизованный канал.

	В эту задачу подходит пункт
	Измените код Go программы workerPool.go таким образом, чтобы реализовать
	функциональность утилиты командной строки wc(1).

	Я сразу решил реализовать эту задачу на воркерах,
	потому как они  используют буферизированый канал
*/

type Worker struct {
	id       int
	fileName string
}

type Data struct {
	job       Worker
	line      int
	words     int
	character int
}

/*
	функция, которая читает запросы из канала worker,
	в котором содержится id воркера и файл, над которым нужно произвести вычисления.
	После завершения, записывается в канал data.
*/
func worker(worker chan Worker, w *sync.WaitGroup, data chan Data) {
	for v := range worker {
		file, err := os.Open(v.fileName)
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
		outPut := Data{
			job:       v,
			line:      lines,
			words:     words,
			character: characters,
		}

		data <- outPut

		/*
			Нужен для того, чтобы наглядно  показать, как работаю воркеры.
			Чем меньше воркеров, тем дольше выполнение задачи
		*/
		time.Sleep(time.Second)
	}
	w.Done()
}

// правильно создать все запросы, используя
// тип Worker, а затем записать их в канал workerChan для обработки
// workerChan читается функцией worker()
func createWorker(workerChan chan Worker, n int, fileName []string) {
	for i := 0; i < n; i++ {
		w := Worker{id: i, fileName: fileName[i]}
		workerChan <- w
	}
	close(workerChan)
}

// формирует нужно количество горутин worker для обработки всех запросов
func makeWorkerPool(data chan Data, workerChan chan Worker, n int) {
	var w sync.WaitGroup
	for i := 0; i < n; i++ {
		w.Add(1)
		go worker(workerChan, &w, data)
	}
	w.Wait()
	close(data)
}

func main() {
	var fileNames []string

	workerCount := 2

	if len(os.Args) <= 1 {
		log.Fatal("no input files")
	}

	fileNames = make([]string, 0)
	fileNames = append(fileNames, os.Args[1:]...)

	fileCount := len(fileNames)

	log.Println("Files count", fileCount)
	log.Println("Worker count", workerCount)

	workerChan := make(chan Worker, workerCount)
	dataChan := make(chan Data, fileCount)
	finished := make(chan struct{})

	go createWorker(workerChan, fileCount, fileNames)

	go func() {
		for d := range dataChan {
			fmt.Printf("Worker ID %d %8d%8d%8d %s\n", d.job.id, d.line, d.words, d.character, d.job.fileName)
		}
		finished <- struct{}{}
	}()

	makeWorkerPool(dataChan, workerChan, workerCount)

	<-finished
}
