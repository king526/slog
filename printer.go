package slog

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	// "cindasoft.com/library/utils"
)

type Printer interface {
	Print(s *string)
	Close()
}

type ConsolePrinter struct {
}

func (this *ConsolePrinter) Print(s *string) {
	fmt.Fprint(os.Stdout, *s)
}

func (this *ConsolePrinter) Close() {
	os.Stdout.Sync()
}

type FilePrinter struct {
	ch          chan *string
	blockMillis time.Duration
	baseName    string
	file        *os.File
	rsize       int
	csize       int
	backup      int
	w           sync.WaitGroup
}

/*
文件输出器
参数：
	maxrn:缓冲队列长度(记录条数)
	dir:日志输出目录,""为当前
	name:日志文件名称,""为程序名
	ksize:日志文件滚动大小，以KB为单位
	backup:日志文件滚动个数
	blockMillis:缓冲队列满时，日志线程阻塞工作线程最大毫秒数。越大则丢日志的可能性越低
*/
type FileConfig struct {
	Maxrn       int
	Dir         string
	Name        string
	SizeKB      int
	Backup      int
	BlockMillis int
}

func NewFilePrinter(cfg FileConfig) (Printer, error) {
	var e error
	if cfg.Name == "" {
		cfg.Name = filepath.Base(os.Args[0])
		if i := strings.LastIndex(cfg.Name, "."); i != -1 {
			cfg.Name = cfg.Name[:i]
		}
	}
	if cfg.Maxrn < 1 {
		cfg.Maxrn = 1024
	}
	if cfg.BlockMillis < 1 {
		cfg.BlockMillis = 100
	}
	if cfg.SizeKB < 1 {
		cfg.SizeKB = 8192
	}
	if cfg.Dir == "" {
		cfg.Dir = "."
	}
	cfg.Dir, e = filepath.Abs(cfg.Dir)
	if e != nil {
		return nil, e
	}
	e = os.MkdirAll(cfg.Dir, os.ModePerm)
	if e != nil {
		return nil, e
	}

	p := &FilePrinter{
		ch:          make(chan *string, cfg.Maxrn),
		rsize:       cfg.SizeKB * 1024,
		blockMillis: time.Millisecond * time.Duration(cfg.BlockMillis),
		backup:      cfg.Backup,
	}
	p.baseName = filepath.Clean(cfg.Dir + string(filepath.Separator) + cfg.Name + ".log")
	p.file, e = os.OpenFile(p.baseName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if e != nil {
		return nil, e
	}
	fs, _ := p.file.Stat()
	p.csize = int(fs.Size())
	t := time.Now()
	s := t.Format("[01-02 15:04:05.999][INFO ] ------------ start ------------\r\n")
	p.Print(&s)
	go p.flush()
	return p, nil
}

func (this *FilePrinter) Print(s *string) {
	select {
	case this.ch <- s:
	case <-time.After(this.blockMillis):
		fmt.Fprint(os.Stderr, "slog timeout:", *s)
	}
}

func (this *FilePrinter) checkfile() {
	if this.csize > this.rsize {
		this.file.Sync()
		this.file.Close()
		this.roll()
		this.file, _ = os.OpenFile(this.baseName, os.O_CREATE|os.O_WRONLY, 0666)
		this.csize = 0
	}
}

func (this *FilePrinter) roll() {
	if this.backup == 0 {
		os.Remove(this.baseName)
		return
	}
	os.Remove(this.baseName + "." + strconv.Itoa(this.backup))
	for i := this.backup - 1; i > 0; i-- {
		o := this.baseName + "." + strconv.Itoa(i)
		n := this.baseName + "." + strconv.Itoa(i+1)
		os.Rename(o, n)
	}
	os.Rename(this.baseName, this.baseName+".1")
}

func (this *FilePrinter) flush() {
	this.w.Add(1)
	defer func() {
		this.file.Close()
		this.w.Done()
	}()
	var n int
	var e error
	for {
		select {
		case s := <-this.ch:
			if s == nil {
				return
			}
			this.checkfile()
			n, e = this.file.WriteString(*s)
			if e != nil {
				this.csize = this.rsize + 1
			} else {
				this.csize += n
			}
		}
	}
}

func (this *FilePrinter) Close() {
	this.ch <- nil
	this.w.Wait()
}
