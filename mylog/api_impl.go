package mylog

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (l *Logger) Init() {
	// 查看目录是否存在，不存在则创建
	if _, err := os.Stat(l.logCft.LogPath); err != nil {
		_ = CreatDir(l.logCft.LogPath)
	}
	// 创建日志文件
	err := l.createFile()
	if err != nil {
		return
	}
}

func (l *Logger) createFile() error {
	fileName := l.getFileName(0)
	l.fullPath = filepath.Join(l.logCft.LogPath, fileName)
	fmt.Println(l.fullPath)
	var err error
	l.fileFd, err = os.OpenFile(l.fullPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	l.log = log.New(l.fileFd, "", log.Lshortfile|log.LstdFlags|log.Lmsgprefix)
	return nil
}

func (l *Logger) LogFileCheck() error {
	var i = 0
	if l.logCft.MaxBackups == 0 {
		l.logCft.MaxBackups = 10
	}
	fileName := filepath.Join(l.logCft.LogPath, l.getFileName(i))
	if l.fullPath != fileName {
		err := l.createFile()
		if err != nil {
			return err
		}
		return nil
	}
	info, err := os.Stat(l.fullPath)
	if err != nil {
		fmt.Println("failed to get file info:", err)
		return err
	}
	fmt.Println("file size:", info.Size())
	size := 1024 * 1024 * l.logCft.MaxSize
	if int64(size) < info.Size() {
		l.reName(i)
	}
	return nil
}

func (l *Logger) reName(i int) {
	fileName := filepath.Join(l.logCft.LogPath, l.getFileName(i))
	if i >= l.logCft.MaxBackups {
		if _, err := os.Stat(fileName); err == nil {
			_ = os.Remove(fileName)
		}
		return
	}
	if _, err := os.Stat(fileName); err == nil {
		l.reName(i + 1)
		l.fileFd.Close()
		err := os.Rename(fileName, filepath.Join(l.logCft.LogPath, l.getFileName(i+1)))
		if err != nil {
			fmt.Println(err)
		}
		err = l.createFile()
		if err != nil {
			return
		}
	}
}

func (l *Logger) getFileName(i int) string {
	dateString := time.Now().Format("2006-01-02")
	str := l.logCft.Filename
	sep := "."
	index := strings.LastIndex(str, sep)
	if index == -1 {
		index = len(str)
	}
	part1 := str[:index]
	part2 := str[index+len(sep):]
	return part1 + dateString + strconv.Itoa(i) + sep + part2
}

// CreatDir 创建目录
func CreatDir(fileDir string) error {
	fmt.Println("makeDir", fileDir)
	err := os.MkdirAll(fileDir, 0666)
	if err != nil {
		fmt.Println("create log dir fail.")
		return err
	}
	return nil
}

func (l *Logger) printLv(ctx context.Context, level Level) string {
	return fmt.Sprintf("[%s][%s] ", ctx.Value("sysNo").(string), level.String())
}

func (l *Logger) Debug(ctx context.Context, args ...interface{}) {
	if l.checkLevel(DebugLevel) {
		str := fmt.Sprint(args)
		l.log.Printf("%s%s", l.printLv(ctx, DebugLevel), str)
	}
}
func (l *Logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	if l.checkLevel(DebugLevel) {
		str := fmt.Sprintf(format, args)
		l.log.Printf("%s%s", l.printLv(ctx, DebugLevel), str)
	}
}
func (l *Logger) Info(ctx context.Context, args ...interface{}) {
	if l.checkLevel(InfoLevel) {
		str := fmt.Sprint(args)
		l.log.Printf("%s%s", l.printLv(ctx, InfoLevel), str)
	}
}
func (l *Logger) Infof(ctx context.Context, format string, args ...interface{}) {
	if l.checkLevel(InfoLevel) {
		str := fmt.Sprintf(format, args)
		l.log.Printf("%s%s", l.printLv(ctx, InfoLevel), str)
	}
}
func (l *Logger) Warn(ctx context.Context, args ...interface{}) {
	if l.checkLevel(WarnLevel) {
		str := fmt.Sprint(args)
		l.log.Printf("%s%s", l.printLv(ctx, WarnLevel), str)
	}
}
func (l *Logger) Warnf(ctx context.Context, format string, args ...interface{}) {
	if l.checkLevel(WarnLevel) {
		str := fmt.Sprintf(format, args)
		l.log.Printf("%s%s", l.printLv(ctx, WarnLevel), str)
	}
}
func (l *Logger) Error(ctx context.Context, args ...interface{}) {
	if l.checkLevel(ErrorLevel) {
		str := fmt.Sprint(args)
		l.log.Printf("%s%s", l.printLv(ctx, ErrorLevel), str)
	}
}
func (l *Logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	if l.checkLevel(ErrorLevel) {
		str := fmt.Sprintf(format, args)
		l.log.Printf("%s%s", l.printLv(ctx, ErrorLevel), str)
	}
}

func (l *Logger) checkLevel(level Level) bool {
	if level < _minLevel || level > _maxLevel {
		return false
	}
	if level < l.logCft.LogLevel {
		return false
	}
	return true
}

// String returns a lower-case ASCII representation of the log level.
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

