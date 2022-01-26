package util

import "sync"

type Set struct {
	m map[interface{}]bool
	sync.RWMutex
}

func NewSet() *Set {
	return &Set{
		m: map[interface{}]bool{},
	}
}

func (s *Set) Add(item interface{}) {
	//写锁
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

func (s *Set) Remove(item interface{}) {
	//写锁
	s.Lock()
	defer s.Unlock()
	delete(s.m, item)
}

func (s *Set) Has(item interface{}) bool {
	//允许读
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

func (s *Set) List() []interface{} {
	//允许读
	s.RLock()
	defer s.RUnlock()
	var outList []interface{}
	for value := range s.m {
		outList = append(outList, value)
	}
	return outList
}

func (s *Set) Len() int {
	return len(s.List())
}

func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[interface{}]bool{}
}

func (s *Set) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}
