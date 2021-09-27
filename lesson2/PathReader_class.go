package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

// ---------PathReader class--------------

type PathReader struct {
	fileENV    string
	Data       []byte
	outputPath string
}

func MakePathReader() PathReader {
	var p PathReader
	p.fileENV = "FILE"
	p.outputPath = "D:/GO/tfs-go-hw/lesson2/output.json"
	return p
}

func (p PathReader) readFromUser() []byte {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	return data
}

func (p PathReader) readArgs() (string, bool) {
	var fileFlag = flag.String("file", "", "path to file")
	flag.Parse()

	if fileFlag == nil {
		er := fmt.Errorf("%s", "Ошибка указателя fileFlag. func readArgs()\n")
		panic(er)
	}

	if *fileFlag != "" {
		return *fileFlag, true
	} else if file, ok := os.LookupEnv(p.fileENV); ok {
		return file, true
	} else {
		return "", false
	}
}

func (p *PathReader) ReadFileData() error {
	filePath, ok := p.readArgs()
	if !ok {
		p.readFromUser()
		return nil
	}

	DEBUGprinter("filePath:", filePath)

	var errF error
	p.Data, errF = ioutil.ReadFile(filePath)
	ErrorCheck(errF)
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
