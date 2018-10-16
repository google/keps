package section

import (
	"sync"
)

type Collection interface {
	Persist() error
	Erase() error
	Sections() []Info
}

type collection struct {
	sections []section
	readme   section
	locker   sync.Mutex
}

func (c *collection) Persist() error {
	c.locker.Lock()
	defer c.locker.Unlock()

	err := c.persist()
	if err != nil {
		c.erase()
		return err
	}

	return nil
}

func (c *collection) Erase() error {
	c.locker.Lock()
	defer c.locker.Unlock()

	return c.erase()
}

func (c *collection) Sections() []Info {
	c.locker.Lock()
	defer c.locker.Unlock()

	infos := []Info{}
	for i := range c.sections {
		infos = append(infos, c.sections[i])
	}

	return infos
}

func (c *collection) persist() error {
	for _, s := range c.sections {
		err := s.Persist()
		if err != nil {
			return err
		}
	}

	if c.readme != nil {
		return c.readme.Persist()
	}

	return nil
}

func (c *collection) erase() error {
	for _, s := range c.sections {
		err := s.Erase()
		if err != nil {
			return err
		}
	}

	return nil
}
