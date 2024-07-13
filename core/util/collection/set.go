package collection

import (
	"golang.org/x/exp/constraints"
)

type Set[T constraints.Ordered] map[T]struct{}

func (s Set[T]) Has(key T) bool {
	_, ok := s[key]
	return ok
}

func (s Set[T]) Add(key T) {
	s[key] = struct{}{}
}

func (s Set[T]) Delete(key T) {
	delete(s, key)
}

func (s Set[T]) AsList() []T {
	var result = make([]T, 0)
	for k, _ := range s {
		result = append(result, k)
	}
	return result
}
