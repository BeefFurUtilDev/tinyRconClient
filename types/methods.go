package types

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	validator "github.com/asaskevich/govalidator"
	"github.com/gorcon/rcon"
	"github.com/rs/zerolog"
)

func (c *Client) NewSession() (err error) {
	if validator.IsIPv6(c.Addr) {
		c.Addr = "[" + c.Addr + "]"
	}
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
		return c.ProcessError().NewSession()
	}
	return nil
}
func (c *Client) DomainToAddress() *Client {
	// 特别感谢 @TPCW-DRW-TP211
	if net.ParseIP(c.Addr) != nil {
		// 如果为 IP, 则原样返回
		return c
	}
	ips, ipsErr := net.LookupIP(c.Addr)
	if ipsErr != nil {
		// 处理域名解析失败
		c.Errors <- ipsErr
		return c
	}
	if len(ips) == 0 {
		// 没找到对应解析
		c.Errors <- fmt.Errorf("Address not found or is CNAME, CNAME not support yet.")
		return c
	}
	// 他不一定有一个地址，取第一个优先
	c.Addr = ips[0].String()
	return c
}
func (c *Client) ProcessError() *Client {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
	if len(c.Errors) == 0 {
		return c
	}
	printErr := func(err error) {
		log.Error().Msg("")
	}
	for err := range c.Errors {
		go printErr(err)
	}
	return c
}
