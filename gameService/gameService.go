package gameService

import (
	"bytes"
	"fmt"
	"log-server/serverEngine/common"
	"log-server/serverEngine/config"
	"log-server/serverEngine/logger"
	"net"
	"strconv"
)

var instance *gameService
// var connCnt int

type gameService struct {
}

func GetInstancePtr() *gameService {
	return instance
}

func init() {
	instance = new(gameService)
	// connCnt = 0
}

func (*gameService) Start() (bool, error) {
	// 切换工作路径
	common.SetCurrentWorkDir("")
	if res, err := logger.GetInstancePtr().Start("LogServer", "./log"); res == false {
		return false, err
	}
	logger.GetInstancePtr().LogWarn([]byte("---------服务器开始启动--------"))
	// 加载配置文件
	areaid, err := config.GetInstancePtr().String("areaid")
	if err != nil {
		panic(err)
	}

	// 判断是否有已经运行的程序
	var bufferSign bytes.Buffer
	bufferSign = common.BufferCombineBytes(bufferSign, []byte("LogServer"))
	bufferSign = common.BufferCombineBytes(bufferSign, []byte(areaid))
	if common.ProcessExists(bufferSign.String()) {
		logger.GetInstancePtr().LogWarn([]byte("LogServer is already running!"))
		return false, nil
	}

	port, err := config.GetInstancePtr().Int("log_svr_port")
	if err != nil {
		logger.GetInstancePtr().LogWarn([]byte("---------服务器启动失败，获取完成端口配置错误--------"))
		panic(err)
	}
	if port <= 0 {
		logger.GetInstancePtr().LogWarn([]byte("---------服务器启动失败，端口配置无效--------"))
		logger.GetInstancePtr().LogWarn([]byte("LogServer config log_svr_port error!"))
	}

	maxConn, err := config.GetInstancePtr().Int("log_svr_max_con")
	if err != nil {
		logger.GetInstancePtr().LogWarn([]byte("---------服务器启动失败，获取最大连接数配置错误--------"))
		panic(err)
	}

	var address bytes.Buffer
	address = common.BufferCombineBytes(address, []byte("0.0.0.0:"))
	address = common.BufferCombineBytes(address, []byte(strconv.Itoa(port)))


	// 创建socket、监听端口、协程处理多路复用的连接请求
	listener, err := net.Listen("tcp", address.String())
	if err != nil {
		logger.GetInstancePtr().LogWarn([]byte("---------服务器启动失败，监听错误--------"))
		panic(err)
	}

	logger.GetInstancePtr().LogWarn([]byte("---------服务器启动完成，监听中--------"))

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go process(conn, maxConn)
	}
}

func process(conn net.Conn, maxConn int) {
	defer conn.Close()
	for {
		// if connCnt < maxConn {
		// 	connCnt++
		// } else {
		// 	fmt.Printf("reached max connection")
		// 	break
		// }
		var buf [128]byte
		// 接收数据
		n, err := conn.Read(buf[:])
		if err != nil {
			// connCnt--
			fmt.Printf("read from connect failed, err: %v\n", err)
			break
		}
		fmt.Printf("receive from client, data: %v\n", string(buf[:n]))
		// 响应数据
		if _, err = conn.Write([]byte("Send from server")); err != nil {
			// connCnt--
			fmt.Printf("write to client failed, err: %v\n", err)
			break
		}
	}
}
