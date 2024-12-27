package types

import (
	"github.com/gorcon/rcon"
	"strconv"
)

func (c *Client) NewSession() (err error) {
	c.Session, err = rcon.Dial(c.Addr+":"+strconv.Itoa(c.Port), c.Password)
	if err != nil {
		return err
	}
	c.Stat = true
	return nil
}
func (c *Client) CloseSession() error {
	if err := c.Session.Close(); err != nil {
		return err
	}
	// cover Stat to false
	if c.Stat {
		c.Stat = !c.Stat
	}
	return nil
}
func (c *Client) ReSession() error {
	if !c.Stat {
		return c.NewSession()
	}
	return nil
}
