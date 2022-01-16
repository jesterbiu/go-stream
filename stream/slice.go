package gostream

import "reflect"

type sliceStream struct {
	slice interface{}
	iter  func() Iterator
}

func (s sliceStream) Iter() Iterator {
	return s.iter()
}

func defaultBeg(s interface{}) func() Iterator {
	return func() Iterator {
		return sliceIter{
			slice: reflect.ValueOf(s),
			index: 0,
		}
	}
}

func FromSlice(s interface{}) Stream {
	if reflect.TypeOf(s).Kind() != reflect.Slice {
		panic(ErrNotASlice)
	}
	return Stream{
		base: &sliceStream{
			slice: s,
			iter:  defaultBeg(s),
		},
	}
}

type sliceIter struct {
	slice reflect.Value
	index int
}

func (si sliceIter) Elem() interface{} {
	return si.slice.Index(si.index).Interface()
}

var emptyVal = reflect.Value{}

func (si sliceIter) NotNil() bool {
	return si.slice != emptyVal && si.index < si.slice.Len()
}

func (si sliceIter) Next() Iterator {
	return &sliceIter{
		slice: si.slice,
		index: si.index + 1,
	}
}
