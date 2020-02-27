package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	rows := []string{
		"Hello Go!",
		"Welcome To Home",
	}
	newFile, err := os.Create("some.txt")
	if err != nil {
		fmt.Println("Ошибка при создании сайта")
		os.Exit(1)
	}
	writer := bufio.NewWriter(newFile)
	defer newFile.Close()

	for _, row := range rows {
		writer.WriteString(row)
		writer.WriteString("\n")
	}
	writer.Flush()
}