package amethyst

import (
	"github.com/phrara/amethyst/interval"
	"sync"
)

type ConnManager struct {
	ConnSet sync.Map
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		ConnSet: sync.Map{},
	}
}

func (c *ConnManager) Append(conn interval.Connecter) {
	c.ConnSet.Store(conn.ConnID(), conn)
}

func (c *ConnManager) Remove(connID uint) {
	c.ConnSet.Delete(connID)
}

func (c *ConnManager) GetConn(connID uint) interval.Connecter {
	value, ok := c.ConnSet.Load(connID)
	if ok {
		return value.(interval.Connecter)
	} else {
		return nil
	}
}

func (c *ConnManager) Sum() int {
	sum := 0
	c.ConnSet.Range(func(key, value any) bool {
		if value.(interval.Connecter).State() == OPENING {
			sum++
		} else {
			c.ConnSet.Delete(key)
		}
		return true
	})
	return sum
}

func (c *ConnManager) Clear() {
	c.ConnSet.Range(func(key, value any) bool {
		value.(interval.Connecter).Close()
		c.ConnSet.Delete(key)
		return true
	})
}
