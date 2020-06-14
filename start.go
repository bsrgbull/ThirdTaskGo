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

	inputFilePath := os.Args[1] //Путь к файлу с адресами

	data, err := ioutil.ReadFile(inputFilePath) //data - байтовый срез подстрок из файла

	if err != nil {
		fmt.Println("Ошибка чтения из файла")
		os.Exit(1)
	}

	var urlSlice []string

	//Получаем из data срез URL-строк, переменная urlSlice
	for _, newURL := range strings.Split(string(data), "\n") {
		//TrimSpace удаляет символы переноса строки
		urlSlice = append(urlSlice, strings.TrimSpace(newURL))
	}

	ch := make(chan []byte)

	for _, url := range urlSlice {

		if url == "" { //Пропускаем пустые строки
			continue
		}

		go fetch(url, ch) //получение html-cтраницы
	}

	for _, url := range urlSlice {
		if url == "" { //Пропускаем пустые строки
			continue
		}
		newPage := <-ch
		if newPage == nil {
			continue
		}
		write(url, newPage) //запись на диск
	}

}

func fetch(url string, ch chan []byte) { //Эта функция делает Get-запрос по url
	//Возвращает html страницу в виде []byte по каналу ch
	responce, err := http.Get(url)
	if err != nil {
		fmt.Println("Не удалось связаться с " + url)
		fmt.Println(err)
		ch <- nil //Вернём nil в случае ошибки запроса
	}

	page, err := ioutil.ReadAll(responce.Body) //преобразуем тело запроса в байтовый срез
	responce.Body.Close()                      // Исключение утечки ресурсов

	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: чтение %s: %v\n", url, err)
		ch <- nil
	}

	ch <- page
}

func write(url string, page []byte) {

	re, err := regexp.Compile(`[[:punct:]]`)

	if err != nil {
		// Если произошла ошибка выводим ее в консоль
		fmt.Println(err)
	} //удаляем знаки препинания из url
	//при помощи regexp.Compile
	resultFolder := re.ReplaceAllString(url, "") //и ReplaceAllStrings

	//Создание папки по URL
	err = os.MkdirAll(os.Args[2]+"/"+
		resultFolder, 0644)
	if err != nil {
		// Если произошла ошибка выводим ее в консоль и выходим из функции
		fmt.Println(err)
		return
	}

	outputFolderPath := os.Args[2]                                        //Папка результатов
	pathToNewFile := outputFolderPath + "/" + resultFolder + "/page.html" //конечный адрес на диске

	err = ioutil.WriteFile(pathToNewFile, page, 0644) //запись

	if err != nil {
		// Если произошла ошибка выводим ее в консоль
		fmt.Println(err)
	}
}
