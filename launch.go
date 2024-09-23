package main

import (
	"bufio"
	"errors"
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

// newSession 建立一个新的RCON会话。
// 它尝试连接到一个Minecraft服务器，然后在一个循环中读取用户输入的命令并将其发送到服务器，直到会话被中断或用户决定退出。
// 参数:
//
//	clientSetup: 包含连接信息（地址、端口和密码）的结构体。
//
// 返回值:
//
//	错误: 如果在建立连接或执行命令时发生错误，则返回相应的错误。
func newSession(clientSetup client) (err error) {
	// 初始化日志输出格式和时间格式
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("starting session...")

	// 尝试连接到RCON服务器
	conn, err := rcon.Dial("rcon://"+clientSetup.addr+":"+strconv.Itoa(clientSetup.port), clientSetup.password)
	if err != nil {
		log.Error().AnErr("conn error:", err).Msgf("can't connect to server")
		return err
	}
	// 确保在函数结束时关闭连接
	defer func(conn *rcon.Conn) {
		_ = conn.Close()
	}(conn)
	// 初始化变量以读取标准输入和处理中断信号
	var stdInput string
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	scanner := bufio.NewScanner(os.Stdin)
	// 主循环：处理命令输入和中断信号
	for {
		select {
		case <-interruptChan:
			// 当收到中断信号时，退出循环
			fmt.Println("\nCaught ^C, exiting...")
			return nil
		default:
			// 打印提示符
			color.Blue.Print("rcon")
			color.Cyan.Print("@")
			color.Yellow.Print(clientSetup.addr)
			color.Cyan.Print(":")
			color.Yellow.Print(strconv.Itoa(clientSetup.port))
			color.White.Print(" >")
			color.Red.Print("#")
			color.White.Print("> ")
			// 读取并处理用户输入
			if scanner.Scan() {
				stdInput = scanner.Text()
			} else {
				// 处理扫描错误
				err := scanner.Err()
				if err != nil {
					switch {
					case errors.Is(err, bufio.ErrTooLong):
						log.Error().Err(err).Msg("input too long")
					case err == io.EOF:
						log.Info().Msg("EOF detected, exiting...")
						return nil
					default:
						log.Error().AnErr("scan error:", err).Msg("can't read input")
					}
				}
			}

			// 处理空输入或exit命令
			if stdInput == "" {
				fmt.Println("")
				continue
			}
			if stdInput == "exit" || stdInput == "stop" {
				return nil
			}
			// 发送命令并处理结果
			result, err := conn.SendCommand(stdInput)
			switch {
			case err == nil:
				continue
			case errors.Is(err, errors.New("connection closed")):
				log.Error().Msg("connection closed, reconnecting...")
				for i := 3; i == 0 || err != nil; i-- {
					time.Sleep(time.Second * 5)
					log.Info().Msgf("retry num: %d, reconnecting in %d seconds...", i, 5)
					conn, err = rcon.Dial("rcon://"+clientSetup.addr+":"+strconv.Itoa(clientSetup.port), clientSetup.password)
				}
				if err != nil {
					log.Error().AnErr("conn error:", err).Msgf("can't connect to server")
					break
				}
			}
			if err != nil {
				log.Error().AnErr("command error:", err).Msg("can't execute command")
				continue
			}
			if result == "" {
				log.Info().Msg("no response.")
				continue
			}
			log.Info().Msg(result)
		}
	}
	defer func(conn *rcon.Conn) {
		_ = conn.Close()
	}(conn)
	return
}
