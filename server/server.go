package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:9999")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	readString, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("не удалось, прочитать")
		return
	}
	readString = strings.TrimSpace(readString)
	fmt.Println(readString)
	split := strings.Split(readString, " ")
	if split[0] == "list" {
		dir, err := ioutil.ReadDir("./")
		writer := bufio.NewWriter(conn)
		for _, value := range dir {
			if !value.IsDir() {
				_, err = writer.Write([]byte(value.Name() + "\n"))
				if err != nil {
					fmt.Println("Не получилось отправть")
					return
				}
			}
		}
		err = writer.Flush()
		if err != nil {
			fmt.Println("dada", err)
		}
	}
	if split[0] == "download" {
		dir, err := ioutil.ReadDir("./")
		if err != nil {
			fmt.Println("не получилось прочитать директорию")
			return
		}
		for _, value := range dir {
			if !value.IsDir() && value.Name() == split[1]{
				file, err := os.Open(split[1])
				if err != nil {
					fmt.Println("Не удалось открыть файл")
					return
				}
				defer file.Close()
				bytes, err := ioutil.ReadAll(file)
				if err != nil {
					fmt.Println("Не удалось прочитать всё с файла")
					return
				}
				newWriter := bufio.NewWriter(conn)
				_, err = newWriter.Write(bytes)
				if err != nil {
					fmt.Println("Не удалось отправить содержимое файла")
					return
				}
				err = newWriter.Flush()
				if err != nil {
					fmt.Println("Не удалось чето то")
					return
				}
				return
			}
		}
	}
	if split[0] == "upload" {
		newReader := bufio.NewReader(conn)
		s, err := newReader.ReadString('\n')
		if err != nil {
			fmt.Println("Не удалось получить наименование файла")
			return
		}
		bytes, err := ioutil.ReadAll(newReader)
		if err != nil {
			fmt.Println("не удалось прочитать всё")
			return
		}
		s = strings.TrimSpace(s)
		file, err := os.Create("downloads/" + s)
		if err != nil {
			fmt.Println("Не удалось создать файл")
			return
		}
		_, err = file.Write(bytes)
		if err != nil {
			fmt.Println("Файл не удалось образовать")
			return
		}
		fmt.Println("Already is ok")
	}
}