package printUtil

import (
	"github.com/gookit/color"
	"strconv"
)

func PS1(addr string, port int) {
	color.Blue.Print("rcon")
	color.Cyan.Print("@")
	color.Yellow.Print(addr)
	color.Cyan.Print(":")
	color.Yellow.Print(strconv.Itoa(port))
	color.White.Print(" >")
	color.Red.Print("#")
	color.White.Print("> ")
}
