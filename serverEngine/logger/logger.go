package logger

import (
	"bufio"
	"bytes"
	"fmt"
	"log-server/serverEngine/common"
	"os"
	"sync"
)

const (
	Log_All = iota
	Log_Error
	Log_Warn
	Log_Info
	Log_None
)

var (
	mu sync.RWMutex
	instance *logger // 实例
	pLogFile *os.File // 文件指针
	writer *bufio.Writer // writer句柄
	logLevel int // log等级 
)

type logger struct {
	LogLevel int
}

func GetInstancePtr() (*logger){
	return instance
}

func init() {
	// 单例，实例=指针
	instance = new(logger)
	instance.LogLevel = Log_None
	logLevel = Log_None
}

func (*logger) Start(strPrefix string, strLogDir string) (bool, error){
	if !common.IsExistsDir(strLogDir) {
		if _, err := common.CreateDir(strLogDir); err != nil {
			fmt.Println(err)
			return false, err
		}
	}
	tm := common.GetCurrTmDay()
	fileName := fmt.Sprintf("%s/%s-%s.log", strLogDir, strPrefix, tm)
	var err error
	pLogFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return false, err
	}
	if pLogFile == nil {
		return false, err
	}
	writer = bufio.NewWriter(pLogFile)
	return true, nil
}

func (*logger) Close() bool{
	if (pLogFile == nil) {
		return false
	}
	writer.Flush()
	pLogFile.Close()
	pLogFile = nil
	return true
}

func (*logger) LogWarn(contentBytes []byte) bool{
	if (logLevel < Log_Warn) {
		return false
	}
	if (pLogFile == nil) {
		return false
	}

	var buffer bytes.Buffer
	// content := "hello"
	// contentBytes := []byte(content)
	eol := "\n"
	buffer = common.BufferCombineBytes(buffer, contentBytes)
	buffer = common.BufferCombineBytes(buffer, []byte(eol))

	mu.Lock()
	_, err := writer.Write(buffer.Bytes())
	if err != nil {
		return false
	}
	writer.Flush()
	// 文件超过限定大小就创建新文件
	// CheckAndCreate()
	mu.Unlock()
	return true
}