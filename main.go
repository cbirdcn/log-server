package main

import (
	// "fmt"
	// "log-server/serverEngine"
	// "log-server/serverEngine/common"
	"fmt"
	"log-server/gameService"
	"log-server/serverEngine/common"
	"log-server/serverEngine/logger"
)

// "fmt"
// "log-server/serverEngine"
// "strconv"

func main() {
	// _ = serverEngine.GetLogInstancePtr()
	// dir, err := serverEngine.GetCurrentWorkDir()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(dir)
	// strconv.Atoi("s")
	common.SetCurrentWorkDir("")
	logger.GetInstancePtr().Start("LogServer", "./log")
	logger.GetInstancePtr().LogWarn([]byte("hello"))
	
	// 启动服务器进程，监听连接，阻塞状态中
	gameService.GetInstancePtr().Start()

	fmt.Println(1)


	// fmt.Println()

	// 创建CrashReport，输出dmp

	// 解析命令行参数与使用
}
