package util

import "sync"

// Set 集合
type Set[T comparable] struct {
	d map[T]struct{}
}

// NewSet 返回 Set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		d: make(map[T]struct{}),
	}
}

// Add 添加
func (s *Set[T]) Add(k T) {
	s.d[k] = struct{}{}
}

// Del 删除
func (s *Set[T]) Del(k T) {
	delete(s.d, k)
}

// Has 查询
func (s *Set[T]) Has(k T) bool {
	_, ok := s.d[k]
	return ok
}

// SyncSet 同步集合
type SyncSet[T comparable] struct {
	l sync.RWMutex
	d map[T]struct{}
}

// NewSyncSet 返回 SyncSet
func NewSyncSet[T comparable]() *SyncSet[T] {
	return &SyncSet[T]{
		d: make(map[T]struct{}),
	}
}

// Add 添加
func (s *SyncSet[T]) Add(k T) {
	s.l.Lock()
	s.d[k] = struct{}{}
	s.l.Unlock()
}

// Del 删除
func (s *SyncSet[T]) Del(k T) {
	s.l.Lock()
	delete(s.d, k)
	s.l.Unlock()
}

// Has 查询
func (s *SyncSet[T]) Has(k T) bool {
	s.l.RLock()
	_, ok := s.d[k]
	s.l.RUnlock()
	return ok
}

// Slice 返回所有
func (s *SyncSet[T]) Slice() []T {
	var t []T
	s.l.RLock()
	for k := range s.d {
		t = append(t, k)
	}
	s.l.RUnlock()
	return t
}
