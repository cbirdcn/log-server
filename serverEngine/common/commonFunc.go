package common

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
	
	"io/ioutil"
	"path/filepath"
)

func GetCurrentWorkDir() (string, error){
	return os.Getwd()
}

func GetCurrentExeDir() (string, error){
	// only linux, eg: /workspace/log-server/main
	exe, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return "", err
	}
	pos := strings.LastIndex(exe, "/")
	// eg: eg:/workspace/log-server
	return exe[:pos], nil
}

func SetCurrentWorkDir(strPath string) (error){
	if strPath == "" {
		var err error
		strPath, err = GetCurrentExeDir()
		if err != nil {
			return err
		}
	}
	return os.Chdir(strPath)
}

func GetCurrTime() int64{
	return time.Now().Unix()
}

func GetCurrTmTime() string{
	now := time.Now()
	return fmt.Sprintf("%d%02d%02d%02d:%02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

func GetCurrTmDay() string{
	now := time.Now()
	return fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())
}

func CreateDir(strDir string) (bool, error) {
	if IsExistsDir(strDir) {
		return true, nil
	} else {
		err := os.Mkdir(strDir, 0777)
		if IsExistsDir(strDir) {
			return true, nil
		} else {
			return false, err
		}
	}
}

func IsExistsDir(strDir string) bool {
	_, err := os.Stat(strDir)
	return err == nil
}

func BufferCombineBytes(buffer bytes.Buffer, pBytes ...[]byte) bytes.Buffer{
	for i := 0; i < len(pBytes); i++ {
		buffer.Write(pBytes[i])
	}
	return buffer
}

// docker 中管理 pidfile 的方法

// PIDFile stored the process id
type pidFile struct {
	path string
}

// just suit for linux
func ProcessExists(pid string) bool {
	if _, err := os.Stat(filepath.Join("/proc", pid)); err == nil {
		return true
	}
	return false
}

func CheckPidFileAlreadyExists(path string) error {
	if pidByte, err := ioutil.ReadFile(path); err == nil {
		pid := strings.TrimSpace(string(pidByte))
		if ProcessExists(pid) {
			return fmt.Errorf("ensure the process:%s is not running pid file:%s", pid, path)
		}
	}
	return nil
}

// NewPIDFile create the pid file 
// 指定路径下生产 pidfile， 文件内容为 pid
func NewPIDFile(path string) (*pidFile, error) {
	if err := CheckPidFileAlreadyExists(path); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(path), os.FileMode(0755)); err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(path, []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil {
		return nil, err
	}
	return &pidFile{path: path}, nil
}

// Remove remove the pid file
func (file pidFile) Remove() error {
	return os.Remove(file.path)
}