package main

import (
	"encoding/json"
	"fmt"
	"os"
)

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

	stat, err := f.InFile.Stat()
	if err != nil {
		panic(err)
	}

	buf := make([]byte, stat.Size())

	_, err = f.InFile.Read(buf)
	ErrorCheck(err)

	var information []Inf

	err = json.Unmarshal(buf, &information)
	ErrorCheck(err)

	parsedStruct := CreateMap(information)

	var sliceOfInformation []trueBill

	for _, bill := range parsedStruct {
		var res trueBill
		res.Valid = bill.Valid
		res.Balance = bill.Balance
		res.Company = bill.Company
		for _, i := range bill.Invalid {
			res.Invalid = append(res.Invalid, i)
		}
		sliceOfInformation = append(sliceOfInformation, res)
	}

	f.CreateFile(sliceOfInformation)

	DEBUGprinter(parsedStruct)

	errF = f.InFile.Close()
	if errF != nil {
		fmt.Println("Ошибка закрытия файла:")
		panic(errF)
	}
}
