package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
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

	var urlSlice []string

	for _, newURL := range strings.Split(string(data), "\n") {
		urlSlice = append(urlSlice, strings.TrimSpace(newURL)) //TrimSpace удаляет символы переноса строки
	}

	fmt.Println(inputFilePath)
	fmt.Println(outputFolderPath)
	fmt.Println(urlSlice[0])

	chanel := make(chan string)

	for _, url := range urlSlice {
		go fetch(url, chanel)
	}

	for range urlSlice {
		fmt.Println(<-chanel)
	}
}

func fetch(url string, chanel chan<- string) {

	responce, err := http.Get(url)
	if err != nil {
		chanel <- fmt.Sprint(err) // Отправка в канал chanel
		return
	}

	b, err := ioutil.ReadAll(responce.Body)
	responce.Body.Close() // Исключение утечки ресурсов
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: чтение %s: %v\n", url, err)
	}
	chanel <- fmt.Sprintf("%s", b)

}
