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
	f := MakePathReader()
	errF := f.ReadFileData()
	ErrorCheck(errF)

	var information []Inf

	errF = json.Unmarshal(f.Data, &information)
	ErrorCheck(errF)

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
}
