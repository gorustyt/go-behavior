package core

import (
	"cmp"
	"sort"
)

type SortMap[T1 cmp.Ordered, T2 any] struct {
	data  map[T1]T2
	index []T1
}

func NewSortMap[T1 cmp.Ordered, T2 any]() *SortMap[T1, T2] {
	return &SortMap[T1, T2]{data: make(map[T1]T2)}
}

func (s *SortMap[T1, T2]) Set(key T1, value T2) {
	s.data[key] = value
	s.index = append(s.index, key)
	sort.Slice(s.index, func(i, j int) bool {
		return s.index[i] < s.index[j]
	})
}

func (s *SortMap[T1, T2]) Front() (value T2, ok bool) {
	if len(s.index) == 0 {
		return
	}
	return s.Get(s.index[0])
}

func (s *SortMap[T1, T2]) Get(key T1) (value T2, ok bool) {
	v, ok := s.data[key]
	return v, ok
}

func (s *SortMap[T1, T2]) Range(fn func(key T1, value T2) (stop bool)) {
	for _, v := range s.index {
		if fn(v, s.data[v]) {
			break
		}
	}
	return
}

func (s *SortMap[T1, T2]) Delete(key T1) {
	delete(s.data, key)
	for i, v := range s.index {
		if v == key {
			if i != len(s.index)-1 {
				copy(s.index[i:], s.index[i+1:])
			}
			s.index = s.index[:len(s.index)-1]
		}
	}
	sort.Slice(s.index, func(i, j int) bool {
		return s.index[i] < s.index[j]
	})
}

func (s *SortMap[T1, T2]) Reset() {
	s.index = s.index[:0]
	s.data = make(map[T1]T2)
}
