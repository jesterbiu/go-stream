package gostream

import (
	"reflect"
	"testing"
)

func TestJust(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	ch := Just(slice).Chan()
	for _, exp := range slice {
		got, ok := <-ch
		if !ok {
			t.Fatalf("exp<%v>: %v, got: nil", reflect.TypeOf(exp), exp)
		}
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("exp<%v>: %v, got<%v>: %v",
				reflect.TypeOf(exp), exp, reflect.TypeOf(got), got)
		}
	}
	s := Just(slice).Slice()
	for i := range slice {
		exp, got := slice[i], s[i]
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("exp<%v>: %v, got<%v>: %v",
				reflect.TypeOf(exp), exp, reflect.TypeOf(got), got)
		}
	}
}

func TestForEach(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	gotSlice := []int{}
	doAppend := func(i interface{}) {
		gotSlice = append(gotSlice, i.(int))
	}
	s := Just(slice).ForEach(doAppend).Chan()
	for i := range slice {
		sink, ok := <-s
		if !ok {
			t.Fatalf("exp: true, got: false")
		}
		if reflect.TypeOf(sink) != reflect.TypeOf(doAppend) {
			t.Errorf("exp<sink>: %v, got<sink>: %v",
				reflect.TypeOf(doAppend), reflect.TypeOf(sink))
		}
		if exp, got := slice[i], gotSlice[i]; exp != got {
			t.Errorf("exp<%v>: %v, got<%v>: %v",
				reflect.TypeOf(exp), exp, reflect.TypeOf(got), got)
		}
	}
}
