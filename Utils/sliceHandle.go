package Utils

import (
	"math"
	"time"
)

func MaxAndMin(interval float64, values ...float64) (max float64, min float64) {
	if len(values) == 0 {
		return
	}
	max = values[0]
	min = values[0]
	for _, val := range values {
		if val > max {
			max = val
		}
		if val <= min || min == 0 {
			min = val
		}
	}
	if interval > 5 {
		max = float64(int(math.Floor(max))/10)*10 + 10
		min = float64(int(math.Floor(min))/10) * 10
	} else {
		max = float64(int(math.Floor(max))/5)*5 + 5
		min = float64(int(math.Floor(min))/5) * 5
	}
	return
}

func NegativeMaxAndMin(interval float64, values ...float64) (max float64, min float64) {
	if len(values) == 0 {
		return
	}
	max = values[0]
	min = values[0]
	for _, val := range values {
		if val > max {
			max = val
		}
		if val <= min {
			min = val
		}
	}
	if interval > 5 {
		max = float64(int(math.Ceil(max))/10) * 10
		min = float64(int(math.Ceil(min))/10)*10 - 10
	} else {
		max = float64(int(math.Ceil(max))/5) * 5
		min = float64(int(math.Ceil(min))/5)*5 - 5
	}
	return
}

func Remove(slice []interface{}, elem interface{}) []interface{} {
	if len(slice) == 0 {
		return slice
	}
	for i, v := range slice {
		if v == elem {
			slice = append(slice[:i], slice[i+1:]...)
			return Remove(slice, elem)
		}
	}
	return slice
}

func RemoveZero(slice []float64) []float64 {
	if len(slice) == 0 {
		return slice
	}
	for i, v := range slice {
		if IfZero(v) {
			slice = append(slice[:i], slice[i+1:]...)
			return RemoveZero(slice)
		}
	}
	return slice
}

func IfZero(arg interface{}) bool {
	if arg == nil {
		return true
	}
	switch v := arg.(type) {
	case int, int32, int16, int64:
		if v == 0 {
			return true
		}
	case float32:
		r := float64(v)
		return math.Abs(r-0) < 0.0000001
	case float64:
		return math.Abs(v-0) < 0.0000001
	case string:
		if v == "" || v == "%%" || v == "%" {
			return true
		}
	case *string, *int, *int64, *int32, *int16, *int8, *float32, *float64, *time.Time:
		if v == nil {
			return true
		}
	case time.Time:
		return v.IsZero()
	default:
		return false
	}
	return false
}

//切片去重
func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}
