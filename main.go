package main

import (
	"flag"
	"github.com/gookit/color"
	"github.com/rs/zerolog"
	"os"
	"time"
)

var (
	addr       = flag.String("addr", "localhost", "address of the server")
	port       = flag.Int("port", 25575, "port of the server")
	password   = flag.String("pass", "", "password of the server")
	launchType = flag.String("mode", "console", "launch console or exec command")
	command    = flag.String("command", "list", "command to execute")
)

func main() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	flag.Parse()
	clientSetup := &client{*addr, *port, *password}
	switch *launchType {
	case "exec":
		result, err := execCommand(clientSetup, command)
		if err != nil {
			log.Fatal().AnErr("exec error:", err).Msgf("can't execute command")
		}
		log.Info().Msgf("result: %s", result)
	case "console":
		color.BgBlack.Print(" ______ __  __   __  __  __  ______  ______  ______  __   __  ______  __      __  ______  __   __  ______  \n/\\__  _/\\ \\/\\ \"-.\\ \\/\\ \\_\\ \\/\\  == \\/\\  ___\\/\\  __ \\/\\ \"-.\\ \\/\\  ___\\/\\ \\    /\\ \\/\\  ___\\/\\ \"-.\\ \\/\\__  _\\ \n\\/_/\\ \\\\ \\ \\ \\ \\-.  \\ \\____ \\ \\  __<\\ \\ \\___\\ \\ \\/\\ \\ \\ \\-.  \\ \\ \\___\\ \\ \\___\\ \\ \\ \\  __\\\\ \\ \\-.  \\/_/\\ \\/ \n   \\ \\_\\\\ \\_\\ \\_\\\\\"\\_\\/\\_____\\ \\_\\ \\_\\ \\_____\\ \\_____\\ \\_\\\\\"\\_\\ \\_____\\ \\_____\\ \\_\\ \\_____\\ \\_\\\\\"\\_\\ \\ \\_\\ \n    \\/_/ \\/_/\\/_/ \\/_/\\/_____/\\/_/ /_/\\/_____/\\/_____/\\/_/ \\/_/\\/_____/\\/_____/\\/_/\\/_____/\\/_/ \\/_/  \\/_/ \n")
		err := newSession(*clientSetup)
		if err != nil {
			log.Warn().AnErr("session error:", err).Msgf("")
		}
	default:
		log.Fatal().Msgf("unknown launch type: %s", *launchType)
	}
}
