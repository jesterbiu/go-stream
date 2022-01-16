package gostream

import (
	"reflect"
	"testing"
)

func Test_sliceStream(t *testing.T) {
	type Case struct {
		name  string
		slice interface{}
	}
	var testCases = []Case{
		{
			name:  "iterate []int",
			slice: []int{1, 2, 3, 4, 5},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				stream   = FromSlice(tc.slice)
				gotSlice []interface{}
			)
			for iter := stream.Iter(); iter.NotNil(); iter = iter.Next() {
				gotSlice = append(gotSlice, iter.Elem())
			}
			sliceVal := reflect.ValueOf(tc.slice)
			for i := 0; i < sliceVal.Len(); i++ {
				exp := sliceVal.Index(i).Interface()
				got := gotSlice[i]
				if !reflect.DeepEqual(exp, got) {
					t.Errorf(
						"exp<%v>: %v, got<%v>: %v",
						reflect.TypeOf(exp), exp, reflect.TypeOf(got), got)
				}

			}
		})
	}
}
