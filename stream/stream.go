package gostream

import (
	"fmt"
	"reflect"
)

type Stream struct {
	recv func() (interface{}, bool)
}

var (
	ErrNotStream = fmt.Errorf("not a stream")
	dummy        = func() (interface{}, bool) { return nil, false }
)

func assertSliceOrArray(s interface{}) bool {
	kind := reflect.TypeOf(s).Kind()
	return kind == reflect.Slice || kind == reflect.Array
}

func assertStream(s interface{}) bool {
	return reflect.TypeOf(s) == reflect.TypeOf(dummy)
}

func fromSlice(s interface{}) Stream {
	var (
		sliceVal = reflect.ValueOf(s)
		i        = 0
		recv     = func() (interface{}, bool) {
			curr := i
			if curr >= sliceVal.Len() {
				return nil, false
			}
			i++
			return sliceVal.Index(curr).Interface(), true
		}
	)
	return Stream{recv}
}

func Just(strm interface{}) Stream {
	switch {
	case assertSliceOrArray(strm):
		return fromSlice(strm)
	case assertStream(strm):
		return strm.(Stream)
	} // end of switch
	panic(ErrNotStream)
}

func (s Stream) Chan() <-chan interface{} {
	ch := make(chan interface{})
	go func(s Stream, ch chan<- interface{}) {
		defer close(ch)
		for {
			s, ok := s.recv()
			if ok {
				ch <- s
			} else {
				break
			}
		}
	}(s, ch)
	return ch
}

func (s Stream) Slice() []interface{} {
	ret := []interface{}{}
	for {
		i, ok := s.recv()
		if !ok {
			break
		} else {
			ret = append(ret, i)
		}
	}
	return ret
}

func (s Stream) Sink() {
	for {
		_, ok := s.recv()
		if !ok {
			return
		}
	}
}

type Consumer = func(interface{})

func (s Stream) ForEach(c Consumer) Stream {
	recv := func() (interface{}, bool) {
		i, ok := s.recv()
		if !ok {
			return nil, false
		}
		c(i)
		return c, true
	}
	return Stream{recv}
}

type Predicate = func(interface{}) bool

func (s Stream) Filter(pred Predicate) Stream {
	recv := func() (interface{}, bool) {
		i, ok := s.recv()
		if !ok {
			return nil, false
		}
		if pred(i) {
			return i, true
		} else {
			return nil, false
		}
	}
	return Stream{recv}
}

type Mapper = func(interface{}) interface{}

func (s Stream) Map(m Mapper) Stream {
	recv := func() (interface{}, bool) {
		i, ok := s.recv()
		if !ok {
			return nil, false
		}
		return m(i), true
	}
	return Stream{recv}
}
