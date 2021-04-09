package helper

import (
	"fmt"
	"io/ioutil"
	"os"
)

type FileFlag int

const (
	// 内容追加
	ContentAppend FileFlag = iota + 1
	// 内容覆盖
	ContentCover
)

// PathExists 判断文件或者目录是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// FilePutContents 写入文件内容
func FilePutContents(path, content string, flag FileFlag) error {

	switch flag {
	case ContentAppend: // 内容追加
		flag = FileFlag(os.O_WRONLY | os.O_APPEND)
	case ContentCover: // 内容覆盖
		flag = FileFlag(os.O_RDWR | os.O_CREATE | os.O_TRUNC)
	}

	// 判断路径是否存在
	exists, err := PathExists(path)
	if err != nil {
		return err
	}
	var f *os.File
	// 文件不存在则创建
	if !exists {
		f, err = os.Create(path)
		if err != nil {
			return err
		}
	} else {
		// 判断是否是文件
		if !IsFile(path) {
			return fmt.Errorf("path is existed, is dirctory please check")
		}
		// 存在则打开
		f, err = os.OpenFile(path, int(flag), os.ModePerm)
		if err != nil {
			return err
		}
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}

// FileGetContents 一口气获取文件内容
func FileGetContents(path string) (string, error) {
	body, err := ioutil.ReadFile(path)
	return string(body), err
}

// IsFile 判断是文件还是目录
func IsFile(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return !fi.IsDir()
}

// CreatFileIfNotExists 如果不存在则创建文件
func CreatFileIfNotExists(path string) error {
	exists, err := PathExists(path)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// CreateDirIfNotExists 如果不存在则创建目录
func CreateDirIfNotExists(path string) error {
	exists, err := PathExists(path)
	if err != nil {
		return err
	}
	if exists && IsFile(path) {
		return fmt.Errorf("错误,路径已存在,并且是文件类型,不可覆盖")
	}
	return os.MkdirAll(path, os.ModePerm)
}

// FormatFileSizeAndUnit 格式化文件大小
func FormatFileSizeAndUnit(fileSize int64) (string, string) {
	if fileSize < 1024 {
		return fmt.Sprintf("%.2f", float64(fileSize)/float64(1)), "B"
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2f", float64(fileSize)/float64(1024)), "KB"
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2f", float64(fileSize)/float64(1024*1024)), "MB"
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2f", float64(fileSize)/float64(1024*1024*1024)), "GB"
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2f", float64(fileSize)/float64(1024*1024*1024*1024)), "TB"
	} else {
		return fmt.Sprintf("%.2f", float64(fileSize)/float64(1024*1024*1024*1024*1024)), "EB"
	}
}
