package types

import "github.com/gorcon/rcon"

type Client struct {
	Addr     string
	Port     int
	Password string
	Session  *rcon.Conn
	Stat     bool
}
