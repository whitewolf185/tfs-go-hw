package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

// ---------PathReader class--------------

type PathReader struct {
	fileENV    string
	InFile     *os.File
	outputPath string
}

func (p *PathReader) Constructor() {
	p.fileENV = "FILE"
	p.outputPath = "D:/Documents/tfs-go-hw/lesson2/output.json"
}

func (p PathReader) readFromUser() (string, bool) {
	var fileUser string
	fmt.Println("Введите файл:")
	_, err := fmt.Scanf("%s", &fileUser)
	if err == nil {
		return fileUser, true
	}
	panic(err)
}

func (p PathReader) readingArgs() (string, bool) {
	var fileFlag = flag.String("file", "", "path to file")
	flag.Parse()

	if fileFlag == nil {
		er := fmt.Errorf("%s", "Ошибка указателя fileFlag. func readingArgs()\n")
		panic(er)
	}

	if *fileFlag != "" {
		return *fileFlag, true
	} else if file, ok := os.LookupEnv(p.fileENV); ok {
		return file, true
	} else {
		return p.readFromUser()
	}
}

func (p *PathReader) ReadPathFile() error {
	filePath, ok := p.readingArgs()
	if !ok {
		var err = errors.New("something wrong in path reading")
		return err
	}

	DEBUGprinter("filePath:", filePath)

	var errF error
	p.InFile, errF = os.Open(filePath)
	for errF != nil {
		if errors.Is(errF, os.ErrNotExist) {
			fmt.Println("Файл не существует. Хотите попробовать еще раз? y/n")
			ans := ""

			_, e := fmt.Scanf("%s", &ans)
			if e != nil {
				panic(e)
			}
			if ans == "n" {
				os.Exit(0)
			}
		} else {
			panic(errF)
		}
		filePath, ok = p.readFromUser()
		if !ok {
			var err = errors.New("something wrong in path reading")
			return err
		}

		p.InFile, errF = os.Open(filePath)
	}

	return nil
}

func (p PathReader) CreateFile(parsedStruct []trueBill) {
	data, err := json.MarshalIndent(parsedStruct, "", "\t")
	if err != nil {
		panic(err)
	}

	tmp := string(data)
	DEBUGprinter(tmp)

	var f *os.File
	f, err = os.Create(p.outputPath)
	ErrorCheck(err)

	_, err = f.WriteString(tmp)
	ErrorCheck(err)

	err = f.Close()
	ErrorCheck(err)
}

// ----------END pathReader----------
