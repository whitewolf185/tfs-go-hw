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

func size(val int32) func() {
	return func() {
		sizeVar = int(val)
		fmt.Println("Значение изменилось на size = ", sizeVar)
	}
}

func times(val int32) func() {
	return func() {
		timesVar = int(val)
		fmt.Println("Значение изменилось на times = ", timesVar)
	}
}

func char(val int32) func() {
	return func() {
		charVar = string(val)
		fmt.Println("Значение изменилось на char = ", charVar)
	}
}

func createPic(values ...T) {
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
	fmt.Printf("Здравствтуй пользователь!\nДля работы данной программы можно ввести 3 параметра:\n" +
		"1.size (вводите, пожалуйста, только нечетные числа больше 0)\n" +
		"2.char\n" +
		"3.times\n" +
		"Чтобы инициализировать параметры, вводите их название и число в формате <name> <number || char>. Верно будет последнее введенное значение.\n" +
		"Если вы хотите закончить ввод, напишите end\n")
	// command
	var (
		param string
	)

	charVar = "X"
	timesVar = 1
	sizeVar = 15

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
			var num int32
			_, err := fmt.Scanf("%c", &num)
			if err != nil {
				fmt.Println(err)
			}
			values = append(values, char(num))

		case "times":
			flag = true
			var num int32
			_, err := fmt.Scanf("%d", &num)
			if err != nil {
				fmt.Println(err)
			}
			values = append(values, times(num))
		case "size":
			flag = true
			var num int32
			_, err := fmt.Scanf("%d", &num)
			if err != nil {
				fmt.Println(err)
			}
			values = append(values, size(num))
		}
	}
	if flag {
		createPic(values...)
	} else {
		createPic()
	}
}
