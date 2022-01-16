package gostream

import (
	"fmt"
)

var (
	ErrNotASlice = fmt.Errorf("not a slice")
)

// change to chan
type Iterator interface {
	Elem() interface{}
	NotNil() bool
	Next() Iterator
}

type streamer interface {
	Iter() Iterator
}

type Stream struct {
	base streamer
}

func (s Stream) Iter() Iterator {
	return s.base.Iter()
}

type predicate = func(interface{}) bool
type mapper = func(interface{}) interface{}
type sinker = func(interface{})

type filterStream struct {
	source Stream
	beg    func() Iterator
}

func (filterStream) Iter() Iterator {
	return nil
}

func (s Stream) Filter(pred predicate) Stream {
	stream := &filterStream{
		source: s,
		beg: func() Iterator {
			return &filterIter{
				pred: pred,
			}
		},
	}
	return Stream{
		base: stream,
	}
}

type filterIter struct {
	curr, next Iterator
	pred       predicate
}

func (iter filterIter) nextElem() (Iterator, bool) {
	if iter.next != nil {
		return iter.next, true
	}
	for next := iter.curr; next.NotNil(); next = next.Next() {
		if iter.pred(next.Elem()) {
			return next, true
		}
	}
	return nil, false
}

func (iter *filterIter) Elem() interface{} {
	if !iter.pred(iter.curr.Elem()) {
		ok := false
		iter.curr, ok = iter.nextElem()
		if !ok {
			return nil
		}
	}
	return iter.curr.Elem()
}

func (iter *filterIter) NotNil() bool {
	return iter.curr != nil && iter.curr.NotNil()
}

func (iter filterIter) Next() Iterator {
	if iter.next == nil {
		iter.next, _ = iter.nextElem()
	}
	return &filterIter{
		curr: iter.next,
		pred: iter.pred,
	}
}

// func (s sliceStream) Map(f func(interface{})) {
// 	si := s.slice
// 	if reflect.TypeOf(si).Kind() != reflect.Slice {
// 		panic(ErrNotASlice)
// 	}
// 	slice := reflect.ValueOf(si)
// 	for i := 0; i < slice.Len(); i++ {
// 		f(slice.Index(i).Interface())
// 	}
// }
