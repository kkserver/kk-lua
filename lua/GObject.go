package lua

import (
	"sync"
	"unsafe"
)

type GObjectId unsafe.Pointer

type GObject struct {
	objects map[GObjectId]interface{}
}

func NewGObject() *GObject {
	v := GObject{}
	v.objects = map[GObjectId]interface{}{}
	return &v
}

func (G *GObject) Get(id GObjectId) interface{} {
	v, ok := G.objects[id]
	if ok {
		return v
	}
	return nil
}

func (G *GObject) Set(id GObjectId, value interface{}) {
	G.objects[id] = value
}

func (G *GObject) Remove(id GObjectId) {
	delete(G.objects, id)
}

type GMutexObject struct {
	GObject
	mutex sync.Mutex
}

func NewGMutexObject() *GMutexObject {
	v := GMutexObject{}
	v.GObject.objects = map[GObjectId]interface{}{}
	return &v
}

func (G *GMutexObject) Get(id GObjectId) interface{} {
	G.mutex.Lock()
	defer G.mutex.Unlock()
	return G.GObject.Get(id)
}

func (G *GMutexObject) Set(id GObjectId, value interface{}) {
	G.mutex.Lock()
	defer G.mutex.Unlock()
	G.GObject.Set(id, value)
}

func (G *GMutexObject) Remove(id GObjectId) {
	G.mutex.Lock()
	defer G.mutex.Unlock()
	G.GObject.Remove(id)
}
