// Parses JSON files of unknown length/object size
// and converts the resulting JSON files to Go struct code
package main

import (
	"JSONtoGo"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// file, _ := os.Open("./testlog2.json")
	fileName := getFileName()

	file, _ := os.Open(fileName)
	defer file.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter struct name: ")
	structName, _ := reader.ReadString('\n')
	structName = strings.ReplaceAll(structName, " ", "")
	structName = strings.TrimSuffix(structName, "\r\n")

	_ = JSONtoGo.CreateStruct(file, structName)

	fmt.Print("Press Enter key to quit")
	fmt.Scanln()
}

func getFileName() string {
	out, err := exec.Command("/python/python", "fileopendialog.py").Output()

	if err != nil {
		return ""
	}

	fileName := string(out[:])
	if filepath.Ext(fileName) != ".json" {
		fmt.Println("File not a JSON type. Please try again.")
		fileName = getFileName()
	}

	return fileName
}
