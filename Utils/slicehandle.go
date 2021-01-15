package Utils

import (
	"math"
	"time"
)

func MaxAndMin(values ...float64) (max float64, min float64) {
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
	max = math.Ceil(max)
	min = math.Floor(min)
	return
}

func NegativeMaxAndMin(values ...float64) (max float64, min float64) {
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
	max = math.Ceil(max)
	min = math.Floor(min)
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
