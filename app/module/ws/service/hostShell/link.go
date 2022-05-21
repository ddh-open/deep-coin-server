package hostShell

import (
	"sync"
)

var (
	once       sync.Once
	linkManage *LinkManage
)

// LinkManage link
type LinkManage struct {
	sync.RWMutex
	links map[string]*Context
}

func NewLinkManage() *LinkManage {
	once.Do(func() {
		linkManage = &LinkManage{
			RWMutex: sync.RWMutex{},
			links:   make(map[string]*Context, 0),
		}
	})
	return linkManage
}

func (l *LinkManage) AddLink(uuid string, link *Context) {
	l.Lock()
	defer l.Unlock()
	l.links[uuid] = link
}

func (l *LinkManage) GetLink(uuid string) *Context {
	l.RLock()
	defer l.RUnlock()
	if link, ok := l.links[uuid]; ok {
		return link
	}
	return nil
}

func (l *LinkManage) RemoveLink(uuid string) {
	l.Lock()
	defer l.Unlock()
	delete(l.links, uuid)
}
