package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

/*Создайте конвейер, который бы читал текстовые файлы, вычислял количество
вхождений заданной фразы в каждом текстовом файле и подсчитывал общее
количество вхождений этой фразы во всех файлах.*/

func readFile(fileName string) string {
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	return string(file)
}

func WriteFileToChannel(fileNames []string, fileBody chan<- string) {
	for i := 0; i < len(fileNames); i++ {
		fileBody <- readFile(fileNames[i])

	}
	close(fileBody)
}

func ReadFromChannel(phrase string, phrasesCount chan<- int, fileBody <-chan string) {
	for x := range fileBody {
		n := strings.Count(x, phrase)
		fmt.Println(x)
		fmt.Println("Count", n)
		fmt.Println("------------------------")

		phrasesCount <- n
	}
	close(phrasesCount)
}

func CountPhrase(phrasesCount <-chan int) {
	count := 0

	for phrase := range phrasesCount {
		count = count + phrase
	}

	fmt.Printf("The sum of phrases %d\n", count)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Need one or more filename")
		return
	}

	A := make(chan string)
	B := make(chan int)

	var filesNames []string

	phrase := os.Args[1]

	for i := 2; i < len(os.Args); i++ {
		filesNames = append(filesNames, os.Args[i])
	}

	fmt.Println("Filenames -", filesNames)
	fmt.Println("Phrase -", phrase)

	go WriteFileToChannel(filesNames, A)
	go ReadFromChannel(phrase, B, A)

	CountPhrase(B)
}
