package ylog

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var mu sync.Mutex
var fileSizeLimit int64
var showConsole bool

func init() {
	fileSizeLimit = 1
	showConsole = true
}

func Info(log, path string) {
	Write(log, path, "INFO")
}

func Error(log, path string) {
	Write(log, path, "ERROR")
}

func Write(log, path, level string) {
	if path == "" {
		path = "./log/log.txt"
	}

	paths, _ := filepath.Split(path)

	mkPath(paths)

	if ok, _ := pathExists(path); !ok {
		newfile, err := os.Create(path)
		if err != nil {
			//TODO::创建文件失败处理
			fmt.Println(err)
		}
		newfile.Close()
	}
	mu.Lock()
	defer mu.Unlock()
	now := time.Now()
	if getFileSize(path) > fileSizeLimit*(1<<20) {
		data := now.Format("2006-01-02") + "/"
		newname := now.Format("15-04-05")
		mkPath(paths + data)

		err := os.Rename(path, paths+data+newname+".txt")

		if err != nil {
			//TODO::文件重命名失败处理
			fmt.Println(err)
		}
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}
	r := bufio.NewWriter(file)
	log = "[ " + level + " ][ " + now.Format("2006-01-02 15:04:05") + " ] " + log + "\n"

	if showConsole {
		fmt.Print(log)
	}
	_, _ = r.WriteString(log)
	_ = r.Flush()

}

func mkPath(paths string) {
	if ok, _ := pathExists(paths); !ok {
		err := os.MkdirAll(paths, 0777)
		if err != nil {
			//TODO::创建文件夹失败处理
			fmt.Println(err)
		}
	}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getFileSize(filename string) int64 {

	var filesize int64
	_ = filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
		filesize = info.Size()
		return nil
	})
	return filesize

}
