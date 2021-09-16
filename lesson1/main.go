package main

import (
	"fmt"
)

const (
	PARAMETERS = 3
)

var values map[string]int32

func createPic() {
	char := "X"
	times := 1
	size := 15
	if len(values) > 0 {
		value, ok := values["char"]
		if ok {
			char = string(value)
		}

		value, ok = values["size"]
		if ok {
			size = int(value)
		}

		value, ok = values["times"]
		if ok {
			times = int(value)
		}

		if size > 0 && size%2 == 1 {
			fmt.Println("size = ", size)
		} else {
			size = 15
			fmt.Println("Значение сброшено до начального. Size = 15")
		}

		fmt.Println("Сhar = ", char)

		if times > 0 {
			fmt.Println("Times = ", times)
		} else {
			times = 1
			fmt.Println("Значение сброшено до начального. Times = 1")
		}
	}
	var result string
	for k := 0; k < times; k++ {
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				switch {
				case i == 0 || i == size-1:
					result += char
				case i == j || i == size-j-1:
					result += char
				default:
					result += " "
				}
			}
			fmt.Println(result)
			result = ""
		}
		fmt.Println()
		fmt.Println()
	}
}

func main() {
	values = make(map[string]int32)
	fmt.Printf("Здравствтуй пользователь!\nДля работы данной программы можно ввести %d параметра:\n1.size (вводите, пожалуйста, только нечетные числа больше 0)\n2.char\n3.times\n"+
		"Чтобы инициализировать параметры, вводите их название и число в формате <name> <number || char>. Верно будет последнее введенное значение.\nЕсли вы хотите закончить ввод, напишите end\n", PARAMETERS)
	// command
	var (
		param string
		char  rune
		num   int32
	)

	for {
		_, err := fmt.Scanf("%s", &param)
		if err != nil {
			fmt.Println(err)
		}

		if param == "end" {
			break
		}

		switch param {
		case "char":
			_, err := fmt.Scanf("%c\n", &char)
			if err != nil {
				fmt.Println(err)
			}

			values["char"] = char

		default:
			_, err := fmt.Scanf("%d", &num)
			if err != nil {
				fmt.Println(err)
			}
			values[param] = num
		}
	}

	createPic()
}
