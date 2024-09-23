package main

import (
	"github.com/jltobler/go-rcon"
	"github.com/rs/zerolog"
	"os"
	"strconv"
	"time"
)

// execCommand 执行服务器的RCON命令。
// 该函数通过RCON协议连接到服务器，并发送指定的命令，然后返回命令的结果或错误。
// 参数:
//   - clientSetup: 包含连接信息（地址、端口和密码）的客户端设置指针。
//   - cmd: 指向要发送的命令的指针。
//
// 返回值:
//   - string: 服务器对命令的响应结果。
//   - error: 如果连接、发送命令或连接关闭时发生错误，则返回该错误。
func execCommand(clientSetup *client, cmd *string) (result string, err error) {
	// 设置日志输出格式和时间格式。
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// 根据客户端设置，尝试建立与服务器的RCON连接。
	conn, err := rcon.Dial("rcon://"+(*clientSetup).addr+":"+strconv.Itoa(clientSetup.port), (*clientSetup).password)
	if err != nil {
		// 如果连接失败，记录错误并返回。
		log.Error().AnErr("conn error:", err).Msgf("can't connect to server")
		return "", err
	}
	// 确保连接在函数返回前关闭。
	defer func(conn *rcon.Conn) {
		_ = conn.Close()
	}(conn)
	// 发送命令并接收结果。
	result, err = conn.SendCommand(*cmd)
	// 记录发送的命令。
	log.Info().Msgf("command: \"%s\" sended!", *cmd)
	if err != nil {
		// 如果发送命令时发生错误，记录错误。
		log.Error().AnErr("send command error:", err).Msgf("can't send command: %d", cmd)
	}
	if result == "" {
		// 如果命令的响应结果为空，记录警告信息。
		log.Warn().Msgf("response is empty!")
	}
	return
}
