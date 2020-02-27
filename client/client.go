package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
)

var cmd = flag.Bool("list", false, "getList")
var download = flag.Bool("download", false, "downloadFile")
var upload = flag.Bool("upload", false, "uploadFile")

func main() {
	log.Println("Try connect to server 127.0.0.1:4545")
	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("client connected successful")
	flag.Parse()
	if *cmd {
		list(conn)
	}
	if *download {
		downloadFromServer(conn, os.Args[2])
	}
	if *upload {
		uploadToServer(conn, os.Args[2])
	}
}

func uploadToServer(conn net.Conn, fileDir string) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing connection %e", err)
		}
	}()
	newWriter := bufio.NewWriter(conn)
	_, err := newWriter.Write([]byte("upload\n"))
	if err != nil {
		fmt.Println("can't send command")
		return
	}
	err = newWriter.Flush()
	if err != nil {
		fmt.Println("не удалось отправить буффер")
		return
	}
	file, err2 := os.Open("client/" + fileDir)
	fmt.Println(fileDir)
	if err2 != nil {
		fmt.Println("Не удалось найти папку или открыть")
		return
	}
	bytes, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		fmt.Println("err", err2)
		return
	}
	writer := bufio.NewWriter(conn)
	base := filepath.Base(fileDir) // ->
	_, err3 := writer.Write([]byte(base + "\n"))
	if err3 != nil {
		fmt.Println("Не удалось отправить нащвание файла")
		return
	}
	err2 = writer.Flush()
	if err2 != nil {
		fmt.Println("Не удалось отправить")
		return
	}
	_, err2 = writer.Write(bytes)
	if err2 != nil {
		fmt.Println("не удалось отправить данные")
		return
	}
	err2 = writer.Flush()
	if err2 != nil {
		fmt.Println("Не удалось отправить")
		return
	}
	fmt.Println("Файл Успешно отправлен")
}

func downloadFromServer(conn net.Conn, fileName string) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing connection %e", err)
		}
	}()
	writer := bufio.NewWriter(conn)
	_, err := writer.Write([]byte("download " + fileName + "\n"))
	if err != nil {
		fmt.Println("Все пошло не так")
		return
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("Не удалось очистить буффер и отправить данные")
		return
	}
	reader := bufio.NewReader(conn)
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("Файл не был получен")
		return
	}
	file, err := os.Create("downloads/" + fileName)
	if err != nil {
		fmt.Println("Не удалось создать файл")
		return
	}
	_, err = file.Write(bytes)
	if err != nil {
		fmt.Println("Не удалось записать файл")
	}
	fmt.Println("Файл успешно скачен")
}

func list(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing connection %e", err)
		}
	}()
	writer := bufio.NewWriter(conn)
	_, err := writer.Write([]byte("list\n"))
	if err != nil {
		fmt.Println("не удалось отправить ")
		return
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("Не удалось очистить буффер и передать значение")
		return
	}
	reader := bufio.NewReader(conn)
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("Bla Bla Bla")
		return
	}
	fmt.Println(string(bytes))
}
