package sections

import (
	"sync"
)

type Collection interface {
	Persist() error
	Erase() error
	Sections() []string
}

type collection struct {
	sections []section
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

func (c *collection) Sections() []string {
	c.locker.Lock()
	defer c.locker.Unlock()

	sectionFilenames := []string{}
	for i := range c.sections {
		sectionFilenames = append(sectionFilenames, c.sections[i].Filename())
	}

	return sectionFilenames
}

func (c *collection) persist() error {
	for _, s := range c.sections {
		err := s.Persist()
		if err != nil {
			return err
		}
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
