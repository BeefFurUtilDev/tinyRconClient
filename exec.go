package main

import (
	"github.com/jltobler/go-rcon"
	"github.com/rs/zerolog"
	"os"
	"strconv"
	"time"
)

func execCommand(clientSetup *client, cmd *string) (result string, err error) {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	conn, err := rcon.Dial("rcon://"+(*clientSetup).addr+":"+strconv.Itoa(clientSetup.port), (*clientSetup).password)
	if err != nil {
		log.Error().AnErr("conn error:", err).Msgf("can't connect to server")
		return "", err
	}
	defer func(conn *rcon.Conn) {
		_ = conn.Close()
	}(conn)
	result, err = conn.SendCommand(*cmd)
	log.Info().Msgf("command: \"%s\" sended!", *cmd)
	if err != nil {
		log.Error().AnErr("send command error:", err).Msgf("can't send command: %d", cmd)
	}
	if result == "" {
		log.Warn().Msgf("response is empty!")
	}
	return
}
