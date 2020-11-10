// linux shell 相关

package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
)

// Command 执行 linux 命令
// 示例: Command("ls","/4d56fsf65")
// 返回参数1: cannot access /4d56fsf65: No such file or directory
// 返回参数2: Linux 标准错误 1
func Command(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		return "", fmt.Errorf("执行程序时,内部错误 %s", err)
	}
	// 正常日志和错误日志异步的从管道获取
	successBuf := bytes.NewBufferString("")
	go func() {
		scan := bufio.NewScanner(stdout)
		for scan.Scan() {
			s := scan.Text()
			successBuf.WriteString(s)
			successBuf.WriteString("\n")
		}
	}()

	// 错误日志
	errBuf := bytes.NewBufferString("")
	scan := bufio.NewScanner(stderr)
	for scan.Scan() {
		s := scan.Text()
		errBuf.WriteString(s)
		errBuf.WriteString("\n")
	}
	// 注意事项,错误日志和成功日志 不能在 Wait 下方代码获取获取
	cmd.Wait()

	// 执行失败，返回错误信息
	if !cmd.ProcessState.Success() {
		return errBuf.String(), fmt.Errorf("1")
	}
	return successBuf.String(), nil
}
