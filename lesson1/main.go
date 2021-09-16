package main

import (
	"fmt"
	"strconv"
)

const (
	PARAMETERS = 3
)

func createPic(values []string) {
	char := "X"
	times := 1
	size := 15

	if len(values) > 0 {
		for _, value := range values {
			var (
				param string
				i     int
			)

			for i = 0; value[i] != '='; i++ {
				param += string(value[i])
			}
			i++

			switch param {
			case "size":
				var res string
				for ; i < len(value); i++ {
					res += string(value[i])
				}
				val, err := strconv.Atoi(res)
				if err == nil {
					size = val
				} else {
					fmt.Println("Что-то не так с цифрами")
				}

				if size > 0 && size%2 == 1 {
					fmt.Println("size = ", size)
				} else {
					size = 15
					fmt.Println("Значение сброшено до начального. Size = 15")
				}

			case "char":
				char = string(value[i])
				fmt.Println("Поле char введено. Сhar = ", char)

			case "times":
				var res string
				for ; i < len(value); i++ {
					res += string(value[i])
				}
				val, err := strconv.Atoi(res)
				if err == nil {
					times = val
				} else {
					fmt.Println("Что-то не так с цифрами")
				}

				if times > 0 {
					fmt.Println("Поле times введено. Times = ", times)
				} else {
					times = 1
					fmt.Println("Значение сброшено до начального. Times = 1")
				}

			default:
				fmt.Println("Значения сброшены до начальных")
				times = 1
				size = 15
				char = "X"
			}
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
	fmt.Printf("Здравствтуй пользователь!\nДля работы данной программы можно ввести %d параметра:\n1.size (вводите, пожалуйста, только нечетные числа больше 0)\n2.char\n3.times\n"+
		"Чтобы инициализировать параметры, вводите их название и число в формате <name>=<number || char>. Верно будет последнее введенное значение.\nЕсли вы хотите закончить ввод, напишите end\n", PARAMETERS)
	// command
	var (
		param string
	)

	// options

	var args []string

	for {
		_, err := fmt.Scanf("%s", &param)
		if err != nil {
			fmt.Println(err)
		}

		if param == "end" {
			break
		}

		args = append(args, param)
	}

	createPic(args)
}
