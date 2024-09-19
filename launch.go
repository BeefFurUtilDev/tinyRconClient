package main

import (
	"bufio"
	"fmt"
	"github.com/gookit/color"
	"github.com/jltobler/go-rcon"
	"github.com/rs/zerolog"
	"io"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func newSession(clientSetup client) (err error) {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("starting session...")

	conn, err := rcon.Dial("rcon://"+clientSetup.addr+":"+strconv.Itoa(clientSetup.port), clientSetup.password)
	if err != nil {
		log.Error().AnErr("conn error:", err).Msgf("can't connect to server")
		return err
	}
	defer func(conn *rcon.Conn) {
		_ = conn.Close()
	}(conn)

	var stdInput string
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-interruptChan:
			fmt.Println("\nCaught ^C, exiting...")
			return nil
		default:
			// Print prompt
			color.Blue.Print("rcon")
			color.Cyan.Print("@")
			color.Yellow.Print(clientSetup.addr)
			color.Cyan.Print(":")
			color.Yellow.Print(strconv.Itoa(clientSetup.port))
			color.White.Print(" >")
			color.Red.Print("#")
			color.White.Print("> ")

			if scanner.Scan() {
				stdInput = scanner.Text()
			} else {
				err := scanner.Err()
				if err != nil {
					if err == bufio.ErrTooLong {
						log.Error().Err(err).Msg("input too long")
					} else if err == io.EOF {
						log.Info().Msg("EOF detected, exiting...")
						return nil
					} else {
						log.Error().AnErr("scan error:", err).Msg("can't read input")
					}
					continue
				}
			}

			if stdInput == "" {
				fmt.Println("")
				continue
			}
			if stdInput == "exit" {
				return nil
			}

			result, err := conn.SendCommand(stdInput)
			if err != nil {
				log.Error().AnErr("command error:", err).Msg("can't execute command")
				continue
			}
			if result == "" {
				log.Info().Msg("no response.")
			} else {
				log.Info().Msg(result)
			}
		}
	}
	return
}

//func newSession(clientSetup client) (err error) {
//	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
//	log := zerolog.New(output).With().Timestamp().Logger()
//	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
//	log.Info().Msg("starting session...")
//
//	conn, err := rcon.Dial("rcon://"+clientSetup.addr+":"+strconv.Itoa(clientSetup.port), clientSetup.password)
//	if err != nil {
//		log.Error().AnErr("conn error:", err).Msgf("can't connect to server")
//		return err
//	}
//	defer func(conn *rcon.Conn) {
//		_ = conn.Close()
//	}(conn)
//
//	var stdInput string
//	interruptChan := make(chan os.Signal, 1)
//	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
//
//	scanner := bufio.NewScanner(os.Stdin)
//
//	for {
//		select {
//		case <-interruptChan:
//			fmt.Println("\nCaught ^C, exiting...")
//			return nil
//		default:
//			// Print prompt
//			color.Blue.Print("rcon")
//			color.Cyan.Print("@")
//			color.Yellow.Print(clientSetup.addr)
//			color.Cyan.Print(":")
//			color.Yellow.Print(strconv.Itoa(clientSetup.port))
//			color.White.Print(" >")
//			color.Red.Print("#")
//			color.White.Print("> ")
//
//			if scanner.Scan() {
//				stdInput = scanner.Text()
//			} else {
//				err := scanner.Err()
//				if err != nil {
//					log.Error().AnErr("scan error:", err).Msg("can't read input")
//					continue
//				}
//			}
//
//			if stdInput == "" {
//				fmt.Println("")
//				continue
//			}
//			if stdInput == "exit" {
//				return nil
//			}
//
//			result, err := conn.SendCommand(stdInput)
//			if err != nil {
//				log.Error().AnErr("command error:", err).Msg("can't execute command")
//				continue
//			}
//			if result == "" {
//				log.Info().Msg("no response.")
//			} else {
//				log.Info().Msg(result)
//			}
//		}
//	}
//	return
//}

//func newSession(clientSetup client) (err error) {
//	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
//	log := zerolog.New(output).With().Timestamp().Logger()
//	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
//	log.Info().Msg("starting session...")
//	conn, err := rcon.Dial("rcon://"+clientSetup.addr+":"+strconv.Itoa(clientSetup.port), clientSetup.password)
//	if err != nil {
//		log.Error().AnErr("conn error:", err).Msgf("can't connect to server")
//		return err
//	}
//	defer func(conn *rcon.Conn) {
//		_ = conn.Close()
//	}(conn)
//	var stdInput string
//	for {
//		//fmt.Print("rcon@" + clientSetup.addr + ":" + strconv.Itoa(clientSetup.port) + " >#>")
//		color.Blue.Print("rcon")
//		color.Cyan.Print("@")
//		color.Yellow.Print(clientSetup.addr)
//		color.Cyan.Print(":")
//		color.Yellow.Print(strconv.Itoa(clientSetup.port))
//		color.White.Print(" >")
//		color.Red.Print("#")
//		color.White.Print("> ")
//		_, err := fmt.Scanln(&stdInput)
//		if err != nil {
//			log.Error().AnErr("scan error:", err).Msg("can't read input")
//		}
//		if stdInput == "" {
//			fmt.Println("")
//			continue
//		}
//		if stdInput == "exit" {
//			return nil
//		}
//		if stdInput == "" {
//
//		}
//		result, err := conn.SendCommand(stdInput)
//		if result == "" {
//			log.Info().Msg("no response.")
//		}
//		log.Info().Msg(result)
//	}
//	return
//}
