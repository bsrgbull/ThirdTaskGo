package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Ошибка. Неправильно заданы аргументы. " +
			"Необходимо ввести 2 аргумента - путь к файлу с адресами и к папке результатов")
		os.Exit(1)
	}

	inputFilePath := os.Args[1]

	data, err := ioutil.ReadFile(inputFilePath)

	if err != nil {
		fmt.Println("Ошибка чтения из файла")
		os.Exit(1)
	}

	var urlSlice []string

	for _, newURL := range strings.Split(string(data), "\n") {
		//TrimSpace удаляет символы переноса строки
		urlSlice = append(urlSlice, strings.TrimSpace(newURL))
	}

	chanel := make(chan string)

	for _, url := range urlSlice {
		//		fmt.Println(url)
		go fetch(url, chanel)
	}

	for range urlSlice {
		fmt.Println(<-chanel)
	}
}

func fetch(url string, chanel chan<- string) {
	//	fmt.Println("555")
	responce, err := http.Get(url)
	if err != nil {
		chanel <- fmt.Sprint(err) // Отправка в канал chanel
		return
	}
	//	fmt.Println("654")
	page, err := ioutil.ReadAll(responce.Body)
	responce.Body.Close() // Исключение утечки ресурсов

	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: чтение %s: %v\n", url, err)
	}
	//	fmt.Println("946")

	re, err := regexp.Compile(`[[:punct:]]`)

	if err != nil {
		// Если произошла ошибка выводим ее в консоль
		fmt.Println(err)
	}
	str := re.ReplaceAllString(url, "")
	fmt.Println(str)

	fmt.Println(os.Args[2] + "/" + str)
	err = os.MkdirAll(os.Args[2]+"/"+str, 0644)
	fmt.Println(err != nil)
	if err != nil {
		// Если произошла ошибка выводим ее в консоль
		fmt.Println(err)
	}

	outputFolderPath := os.Args[2]
	pathToNewFile := outputFolderPath + "/" + str + "/page.html"
	fmt.Println(pathToNewFile)

	err = ioutil.WriteFile(pathToNewFile, page, 0644)
	fmt.Println(err != nil)

	if err != nil {
		// Если произошла ошибка выводим ее в консоль
		fmt.Println(err)
	}

	chanel <- fmt.Sprintf("%s", err)
}
