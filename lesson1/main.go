package main

import (
	"fmt"
)

type T func()

var (
	charVar  string
	timesVar int
	sizeVar  int
)

func size() {
	fmt.Println("Напишите параметр size")
	_, err := fmt.Scanf("%d", &sizeVar)
	if err != nil {
		fmt.Println(err)
	}
}

func times() {
	fmt.Println("Напишите параметр times")
	_, err := fmt.Scanf("%d", &timesVar)
	if err != nil {
		fmt.Println(err)
	}
}

func char() {
	fmt.Println("Напишите параметр char")
	_, err := fmt.Scanf("%s", &charVar)
	if err != nil {
		fmt.Println(err)
	}
}

func createPic(values ...T) {
	charVar = "X"
	timesVar = 1
	sizeVar = 15

	if len(values) > 0 {
		for _, value := range values {
			value()
		}

		if sizeVar > 0 && sizeVar%2 == 1 {
			fmt.Println("size = ", sizeVar)
		} else {
			sizeVar = 15
			fmt.Println("Значение сброшено до начального. Size = 15")
		}

		fmt.Println("Сhar = ", charVar)

		if timesVar > 0 {
			fmt.Println("Times = ", timesVar)
		} else {
			timesVar = 1
			fmt.Println("Значение сброшено до начального. Times = 1")
		}
	}
	var result string
	for k := 0; k < timesVar; k++ {
		for i := 0; i < sizeVar; i++ {
			for j := 0; j < sizeVar; j++ {
				switch {
				case i == 0 || i == sizeVar-1:
					result += charVar
				case i == j || i == sizeVar-j-1:
					result += charVar
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
	fmt.Printf("Здравствтуй пользователь!\nВведите данные, которые вы хотите изменить. Можно изменить параметры:\n" +
		"1.size по умолнанию = 15\n" +
		"2.char по умолчанию = X\n" +
		"3.times по умолчанию = 1\n" +
		"Если вы хотите закончить ввод, напишите end\n")
	// command
	var (
		param string
	)

	flag := false
	var values []T

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
			flag = true
			values = append(values, char)

		case "times":
			flag = true
			values = append(values, times)
		case "size":
			values = append(values, size)
			flag = true
		}
	}
	if flag {
		createPic(values...)
	} else {
		createPic()
	}
}
