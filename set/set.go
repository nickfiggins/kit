package set

import "encoding/json"

type Set[T comparable] map[T]struct{}

func New[T comparable](vals ...T) Set[T] {
	s := make(Set[T], len(vals))
	for _, v := range vals {
		s[v] = struct{}{}
	}
	return s
}

func FromValues[K, V comparable](m map[K]V) Set[V] {
	s := make(Set[V], len(m))
	for _, v := range m {
		s[v] = struct{}{}
	}
	return s
}

func FromKeys[K comparable, V any](m map[K]V) Set[K] {
	s := make(Set[K], len(m))
	for k := range m {
		s[k] = struct{}{}
	}
	return s
}

func FromList[T comparable](l []T) Set[T] {
	s := make(Set[T], len(l))
	for _, v := range l {
		s[v] = struct{}{}
	}
	return s
}

func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}

func (s Set[T]) Remove(v T) {
	delete(s, v)
}

func (s Set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

func (s Set[T]) Empty() bool {
	return len(s) == 0
}

func (s Set[T]) Clear() {
	for v := range s {
		s.Remove(v)
	}
}

func (s Set[T]) Union(o Set[T]) Set[T] {
	result := New[T]()
	for v := range s {
		result.Add(v)
	}
	for v := range o {
		result.Add(v)
	}
	return result
}

func (s Set[T]) Intersection(o Set[T]) Set[T] {
	result := New[T]()
	for v := range s {
		if o.Contains(v) {
			result.Add(v)
		}
	}
	return result
}

func (s *Set[T]) UnmarshalJSON(b []byte) error {
	var vals []T
	if err := json.Unmarshal(b, &vals); err != nil {
		return err
	}
	n := New[T](vals...)
	*s = n
	return nil
}

func (s *Set[T]) MarshalJSON() ([]byte, error) {
	var vals []T
	for v := range *s {
		vals = append(vals, v)
	}
	return json.Marshal(vals)
}
