package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

//---------PathReader class--------------

type PathReader struct {
	fileENV string
	File    *os.File
}

func (p *PathReader) Constructor() {
	p.fileENV = "FILE"
}

func (p PathReader) readFromUser() (string, bool) {
	var fileUser string
	fmt.Println("Введите файл:")
	_, err := fmt.Scanf("%s", &fileUser)
	if err == nil {
		return fileUser, true
	} else {
		panic(err)
		return "", false
	}
}

func (p PathReader) readingArgs() (string, bool) {
	var fileFlag = flag.String("file", "", "path to file")
	flag.Parse()

	if fileFlag == nil {
		er := fmt.Errorf("Ошибка указателя fileFlag. func readingArgs()\n")
		panic(er)
		return "", false
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
	p.File, errF = os.Open(filePath)
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

		p.File, errF = os.Open(filePath)
	}

	return nil
}

//----------END pathReader----------

func DEBUGprinter(val ...interface{}) {
	var printer []interface{}
	printer = append(printer, "[DEBUG] ")
	printer = append(printer, val...)
	_, err := fmt.Fprintln(os.Stdout, printer...)
	if err != nil {
		panic(err)
	}
}

func ErrorCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	var f PathReader
	f.Constructor()
	errF := f.ReadPathFile()
	ErrorCheck(errF)

	stat, err := f.File.Stat()
	if err != nil {
		panic(err)
	}

	buf := make([]byte, stat.Size())

	_, err = f.File.Read(buf)
	ErrorCheck(err)

	var information []Inf

	err = json.Unmarshal(buf, &information)
	ErrorCheck(err)

	//DEBUGprinter(string(buf))
	DEBUGprinter(information)

	errF = f.File.Close()
	if errF != nil {
		fmt.Println("Ошибка закрытия файла:")
		panic(errF)
	}

}

func (inf *Inf) UnmarshalJSON(data []byte) error {
	var tmpInf anotherInf
	err := json.Unmarshal(data, &tmpInf)
	if err != nil {
		return err
	}

	if len(tmpInf.Operation) != 0 {
		for name, val := range tmpInf.Operation {
			switch name {
			case "type":
				tmpInf.Type = val
			case "value":
				tmpInf.Value = val
			case "id":
				tmpInf.Id = val
			case "created_at":
				s, ok := val.(string)
				if ok {
					tmpInf.CreatedAt = s
				}
			}
		}
	}

	inf.Company = tmpInf.Company
	inf.Type = tmpInf.Type
	inf.Value = tmpInf.Value
	inf.Id = tmpInf.Id
	inf.CreatedAt = tmpInf.CreatedAt

	return nil
}

type anotherInf struct {
	//name        type                     json tags
	Company   string                 `json:"company"`
	Type      interface{}            `json:"type"`
	Operation map[string]interface{} `json:"operation,omitempty"`
	Value     interface{}            `json:"value"`
	Id        interface{}            `json:"id"`
	CreatedAt string                 `json:"created_at"`
}

type Inf struct {
	//name        type                     json tags
	Company   string      `json:"company"`
	Type      interface{} `json:"type"`
	Value     interface{} `json:"value"`
	Id        interface{} `json:"id"`
	CreatedAt string      `json:"created_at"`
}
