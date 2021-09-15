package main

import (
	"fmt"
)

const (
	PARAMETERS = 3
)

func createPic(size int32, char string, times int32) {
	//init
	if size == 0 {
		size = 15
	}
	if char == "" {
		char = "X"
	}
	if times == 0 {
		times = 1
	}

	for i := 0; i < int(times); i++ {
		for i := 0; i < int(size); i++ {
			for j := 0; j < int(size); j++ {
				if i == 0 || i == int(size-1) {
					fmt.Printf("%s", char)
				} else if i == j || i == int(size)-j-1 {
					fmt.Printf("%s", char)
				} else {
					fmt.Printf("%s", " ")
				}
			}
			fmt.Println("")
		}
		fmt.Println("")
		fmt.Println("")
	}

}

func enterVar_str() string {
	var res string
	_, err := fmt.Scanf("%s", &res)
	if err != nil {
		fmt.Println(err)
	}

	return res
}

func enterVar_int() int32 {
	var res int32
	_, err := fmt.Scanf("%d", &res)
	if err != nil {
		fmt.Println(err)
	}

	return res
}

func main() {
	fmt.Printf("Здравствтуй пользователь!\nДля работы данной программы можно ввести %d параметра:\n1.size (вводите, пожалуйста, только нечетные числа больше 0)\n2.char\n3.times\n"+
		"Чтобы инициализировать параметры, вводите их название и число в формате <name> <number || char>\nЕсли вы хотите закончить ввод, напишите end\n", PARAMETERS)
	//command
	var (
		param string
	)

	//options
	var (
		size  int32
		char  string
		times int32
	)

	for {
		param = enterVar_str()

		if param == "end" {
			break
		}

		switch param {

		case "size":
			size = enterVar_int()
			if size > int32(0) && size%2 == 1 {
				fmt.Println("Поле size введено. Size = ", size)
			} else {
				size = 0
				fmt.Println("Попробуйте еще раз")
			}

		case "char":
			char = string(enterVar_str()[0])
			fmt.Println("Поле char введено. Сhar = ", char)

		case "times":
			times = enterVar_int()
			if times > 0 {
				fmt.Println("Поле char введено. Times = ", times)
			} else {
				times = 0
				fmt.Println("Попробуйте еще раз")
			}

		default:
			fmt.Println("Попробуйте еще раз")
		}

		//fmt.Println(param)
	}

	createPic(size, char, times)

}
