package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Ошибка. Неправильно заданы аргументы. " +
			"Необходимо ввести 2 аргумента - путь к файлу с адресами и к папке результатов")
		os.Exit(1)
	}

	inputFilePath := os.Args[1]
	outputFolderPath := os.Args[2]

	data, err := ioutil.ReadFile(inputFilePath)

	if err != nil {
		fmt.Println("Ошибка чтения из файла")
		os.Exit(1)
	}

	var url []string

	for _, newURL := range strings.Split(string(data), "\n") {
		url = append(url, newURL)
	}

	fmt.Println(inputFilePath)
	fmt.Println(outputFolderPath)
	fmt.Println(url[2])

}
