package engine

import "sync"

type objectsMap struct {
	Items map[*Object]bool
	sync.RWMutex
}

func (m *objectsMap) Set(o *Object) {
	m.Lock()
	m.Items[o] = true
	m.Unlock()
}

// func (m *objectsMap) Get(key int) (Info, bool) {
// 	m.RLock()
// 	item, ok := m.Items[key]
// 	m.RUnlock()
// 	return item, ok
// }

// func (m *objectsMap) Has(key int) bool {
// 	m.RLock()
// 	_, ok := m.Items[key]
// 	m.RUnlock()
// 	return ok
// }

func (m *objectsMap) Del(o *Object) {
	m.Lock()
	delete(m.Items, o)
	m.Unlock()
}
