package main

import (
	// "fmt"
	"reflect"
)

func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func GenerateCombinations(alphabet []interface{}, length int) <-chan []interface{} {
	c := make(chan []interface{})

	go func(c chan []interface{}) {
		defer close(c)
		emptySlice := make([]interface{}, 0)
		AddOption(c, emptySlice, alphabet, length)
	}(c)

	return c
}

func AddOption(c chan []interface{}, combo []interface{}, alphabet []interface{}, length int) {
	if length <= 0 {
		c <- combo
		return
	}
	var newCombo []interface{}
	for _, ch := range alphabet {
		newCombo = append(combo, ch)
		AddOption(c, newCombo, alphabet, length-1)
	}
}
