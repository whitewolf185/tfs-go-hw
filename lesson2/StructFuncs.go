package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

func CreateMap(information []Inf) map[string]Bill {
	result := make(map[string]Bill)

	for _, rhs := range information {
		// если я не могу использовать ID, то и не смогу операцию записать в невалидную. Поэтому я ее не учитываю
		if !ParseID(rhs.ID) {
			DEBUGprinter("Ошибка ID", rhs)
			continue
		}
		if rhs.Company == "" {
			DEBUGprinter("Ошибка company", rhs) // без этого параментра не понятно к чему присваивать
			continue
		}
		if rhs.CreatedAt == "" {
			DEBUGprinter("Ошибка created_at", rhs) // без данного параметра не понятно как сортировать
			continue
		}

		mapVal, ok := result[rhs.Company]
		if !ok {
			mapVal.Invalid = make(map[int64]interface{})
		}
		err := ParseVal(&mapVal, rhs)
		mapVal.Company = rhs.Company

		// добавление валидной операции
		if err == nil {
			mapVal.Valid++
		} else { // добавление невалидной операции
			myTime, err := time.Parse(time.RFC3339, rhs.CreatedAt)
			if err != nil {
				panic(err)
			}

			mapVal.Invalid[myTime.Unix()] = rhs.ID
		}

		result[rhs.Company] = mapVal
	}
	return result
}

func ParseID(iD interface{}) bool {
	switch v := iD.(type) {
	case string:
		return true
	case nil:
		return false
	case int:
		return true
	case float64:
		if v == float64(int(v)) {
			return true
		}
		return false
	case float32:
		if v == float32(int(v)) {
			return true
		}
		return false
	default:
		return false
	}
}

func ParseVal(mapVal *Bill, rhs Inf) error {
	var err = errors.New("non valid operation")

	// добавление значения
	switch rhs.Type {
	case "income":
		err = AddVal(mapVal, rhs)
	case "+":
		err = AddVal(mapVal, rhs)
	case "-":
		err = SubVal(mapVal, rhs)
	case "outcome":
		err = SubVal(mapVal, rhs)
	}
	return err
}

func AddVal(mapVal *Bill, rhs Inf) error {
	val, err := CheckVal(rhs.Value)
	if err != nil {
		DEBUGprinter("надо добавить")
		return err
	}
	mapVal.Balance += val
	return nil
}

func SubVal(mapVal *Bill, rhs Inf) error {
	val, err := CheckVal(rhs.Value)
	if err != nil {
		DEBUGprinter("надо добавить")
		return err
	}
	mapVal.Balance -= val
	return nil
}

func CheckVal(val interface{}) (int, error) {
	switch v := val.(type) {
	case string:
		return strconv.Atoi(v)
	case int:
		return v, nil
	case float32:
		if v == float32(int(v)) {
			return int(v), nil
		}
		return 0, errors.New("non valid operation")
	case float64:
		if v == float64(int(v)) {
			return int(v), nil
		}
		return 0, errors.New("non valid operation")
	default:
		return 0, errors.New("non valid operation")
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
				s, ok := val.(string)
				if ok {
					tmpInf.Type = s
				}
			case "value":
				tmpInf.Value = val
			case "id":
				tmpInf.ID = val
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
	inf.ID = tmpInf.ID
	inf.CreatedAt = tmpInf.CreatedAt

	return nil
}

type anotherInf struct {
	// name        type                  json tags
	Company   string                 `json:"company"`
	Type      string                 `json:"type"`
	Operation map[string]interface{} `json:"operation,omitempty"`
	Value     interface{}            `json:"value"`
	ID        interface{}            `json:"id"`
	CreatedAt string                 `json:"created_at"`
}

type Inf struct {
	// name        type       json tags
	Company   string      `json:"company"`
	Type      string      `json:"type"`  // non valid when not income, outcome, +, - AND nil
	Value     interface{} `json:"value"` // int float string | non valid when not int AND nil
	ID        interface{} `json:"id"`    // int string | non valid when not type AND nil
	CreatedAt string      `json:"created_at"`
}

type Bill struct {
	Company string                `json:"company"`
	Valid   int                   `json:"valid_operations_count"`
	Balance int                   `json:"balance"`
	Invalid map[int64]interface{} `json:"invalid_operations,omitempty"`
}

type trueBill struct {
	Company string        `json:"company"`
	Valid   int           `json:"valid_operations_count"`
	Balance int           `json:"balance"`
	Invalid []interface{} `json:"invalid_operations,omitempty"`
}
