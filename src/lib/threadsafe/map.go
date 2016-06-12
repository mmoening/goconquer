package threadsafe

import "sync"

// SafeMap encapsulates a map and lock to create a read-write, thread-safe map. The
// structure is used for the server's client and map collections, which change
// frequently when players login, move, and perform actions. 
type SafeMap struct {
	Elements map[interface{}]interface{}
	sync.RWMutex
}

func NewSafeMap() *SafeMap {
	m := &SafeMap {}
	m.Elements = make(map[interface{}]interface{})
	return m
}

func (sm *SafeMap) Add(key interface{}, value interface{}) bool {
	sm.Lock()
	defer sm.Unlock()
	
	_, exists := sm.Elements[key]
	if !exists { sm.Elements[key] = value }
	return !exists
}

func (sm *SafeMap) Contains(key interface{}) bool {
	sm.RLock()
	defer sm.RUnlock()
	
	_, result := sm.Elements[key]
	return result
}

func (sm *SafeMap) Count() int {
	return len(sm.Elements)
}

func (sm *SafeMap) Get(key interface{}) interface{} {
	sm.RLock()
	defer sm.RUnlock()
	
	result, exists := sm.Elements[key]
	if exists { return result }
	return nil
}

func (sm *SafeMap) Remove(key interface{}) interface{} {
	sm.Lock()
	defer sm.Unlock()
	
	result, exists := sm.Elements[key]
	if exists { 
		delete(sm.Elements, key)
		return result 
	}
	return nil
}
