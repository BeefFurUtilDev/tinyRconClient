package main

import (
	"flag"
	"fmt"
	"github.com/BeefFurUtilDev/tinyRconClient/connFunc"
	"github.com/BeefFurUtilDev/tinyRconClient/printUtil"
	"github.com/BeefFurUtilDev/tinyRconClient/types"
	"github.com/rs/zerolog"
	"os"
	"time"
)

// 全局变量定义
var (
	addr       = flag.String("addr", "localhost", "address of the server")        // 服务器地址，默认为localhost
	port       = flag.Int("port", 25575, "port of the server")                    // 服务器端口，默认为25575
	password   = flag.String("pass", "", "password of the server")                // 服务器密码，必需项
	launchType = flag.String("mode", "console", "launch console or exec command") // 启动模式，可以是console或exec，默认为console
	command    = flag.String("command", "list", "command to execute")             // 执行的命令，默认为list
)

// main函数是程序的入口点
func main() {
	// 初始化日志输出配置
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// 解析命令行标志
	flag.Parse()

	// 设置客户端配置
	clientSetup := &types.Client{Addr: *addr, Port: *port, Password: *password}

	// 根据启动类型执行相应逻辑
	switch *launchType {
	case "exec":
		// 如果命令模式下未指定命令，则提示用户输入
		for *command == "" {
			log.Info().Msg("no command input, please input command:")
			_, err := fmt.Scanln(*command)
			if err != nil {
				log.Fatal().AnErr("command error:", err).Msgf("can't read command")
			}
		}
		// 执行命令并处理结果
		result, err := connFunc.ExecCommand(clientSetup, command)
		if err != nil {
			log.Fatal().AnErr("exec error:", err).Msgf("can't execute command")
		}
		log.Info().Msgf("result: %s", result)
	case "console":
		// 在控制台模式下启动会话
		printUtil.Hello()
		err := connFunc.NewSession(*clientSetup)
		if err != nil {
			log.Warn().AnErr("session error:", err).Msgf("")
		}
	default:
		// 处理未知的启动类型
		log.Fatal().Msgf("unknown launch type: %s", *launchType)
	}
}
