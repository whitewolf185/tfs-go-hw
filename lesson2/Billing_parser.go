package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
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

func KeysSort(invalid *map[int64]interface{}) []int64 {
	keys := make([]int64, len(*invalid))
	i := 0
	for k := range *invalid {
		keys[i] = k
		i++
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	return keys
}

func main() {
	f := MakePathReader()
	errF := f.ReadFileData()
	ErrorCheck(errF)

	var information []ParsedStruct

	errF = json.Unmarshal(f.Data, &information)
	for i, val := range information {
		OperationParser(&val, val.Operation)
		information[i] = val
	}
	ErrorCheck(errF)

	parsedStruct := CreateMap(information)

	var sliceOfInformation []trueBill

	for _, bill := range parsedStruct {
		var res trueBill
		res.Valid = bill.Valid
		res.Balance = bill.Balance
		res.Company = bill.Company

		keys := KeysSort(&bill.Invalid)

		for _, i := range keys {
			res.Invalid = append(res.Invalid, bill.Invalid[i])
		}
		sliceOfInformation = append(sliceOfInformation, res)
	}

	f.CreateFile(sliceOfInformation)

	DEBUGprinter(parsedStruct)
}
