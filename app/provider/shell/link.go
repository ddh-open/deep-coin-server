package shell

import (
	"sync"
)

// LinkManage link
type LinkManage struct {
	sync.RWMutex
	links map[string][]*Context
}

func NewLinkManage() *LinkManage {
	return &LinkManage{
		RWMutex: sync.RWMutex{},
		links:   make(map[string][]*Context, 0),
	}
}

func (l *LinkManage) AddLink(hostId string, link *Context) {
	l.Lock()
	defer l.Unlock()
	if links, ok := l.links[hostId]; ok {
		links = append(links, link)
	} else {
		links = append([]*Context{}, link)
	}
}

func (l *LinkManage) GetLink(hostId string, index int) *Context {
	l.RLock()
	defer l.RUnlock()
	if links, ok := l.links[hostId]; ok {
		if len(links) >= index {
			return links[index]
		}
	}
	return nil
}

func (l *LinkManage) RemoveLink(hostId string, index int) {
	l.Lock()
	defer l.Unlock()
	if links, ok := l.links[hostId]; ok {
		if index+1 <= len(links) {
			links = append(links[0:index], links[0:index+1]...)
		} else {
			links = links[0:index]
		}
	}
}
