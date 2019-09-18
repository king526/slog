package slog

import (
	"os"
	"runtime"
	"strings"
)

var (
	global = NewBuilder(LevDEBUG).Build()
)

type lev int

const (
	LevDEBUG   lev = 0
	LevVERBOSE lev = 1
	LevINFO    lev = 2
	LevWARN    lev = 3
	LevERROR   lev = 4
	LevNOTE    lev = 5
	LevFATAL   lev = 6
)

func (this lev) String() string {
	switch this {
	case LevDEBUG:
		return "DEBUG"
	case LevVERBOSE:
		return "VERBO"
	case LevINFO:
		return "INFO"
	case LevWARN:
		return "WARN"
	case LevERROR:
		return "ERROR"
	case LevNOTE:
		return "NOTE"
	default:
		return "FATAL"
	}
}

func StringLev(l string) lev {
	l = strings.ToUpper(l)
	switch l {
	case "DEBUG":
		return LevDEBUG
	case "ERROR":
		return LevERROR
	case "INFO":
		return LevINFO
	case "WARN":
		return LevWARN
	case "FATAL":
		return LevFATAL
	case "NOTE":
		return LevNOTE
	default:
		return LevVERBOSE
	}
}

type levPrinter struct {
	l lev
	p Printer
}

type LoggerBuilder struct {
	log Logger
}

func NewBuilder(consoleLev lev) *LoggerBuilder {
	l := &LoggerBuilder{}
	l.AddPrinter(consoleLev, &ConsolePrinter{})
	return l
}

func (l *LoggerBuilder) AddPrinter(lv lev, p Printer) *LoggerBuilder {
	l.log.lp = append(l.log.lp, levPrinter{lv, p})
	if lv < l.log.rootLev {
		l.log.rootLev = lv
	}
	return l
}

func (l *LoggerBuilder) Build() *Logger {
	return &l.log
}

func SetGlobalBuilder(b *LoggerBuilder) {
	global = b.Build()

}

func SetGlobal(consoleLev , logPath string, fileLev lev) {
	b := NewBuilder(StringLev(consoleLev))
	if logPath != "" {
		fp, err := NewFilePrinter(FileConfig{Dir: logPath})
		if err != nil {
			panic(err)
		}
		b.AddPrinter(fileLev, fp)
	}
	global = b.Build()
}

func GlobalLogger()*Logger{
	return global
}

type Logger struct {
	rootLev lev
	f       _DefaultFormater
	lp      []levPrinter
}

func (l *Logger)Write(bs []byte)(int,error){
	for _, n := range l.lp {
		s:=string(bs)
		n.p.Print(&s)
	}
	return len(bs),nil
}

func (l *Logger) Debug(args ...interface{}) {
	l.out(LevDEBUG, args...)
}

func (l *Logger) Verbose(args ...interface{}) {
	l.out(LevVERBOSE, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.out(LevINFO, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.out(LevWARN, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.out(LevERROR, args...)
}
func (l *Logger) Note(args ...interface{}) {
	l.out(LevNOTE, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.printStack(true)
	l.out(LevFATAL, args...)
	os.Exit(3)
}

func (l *Logger) Exit(args ...interface{}) {
	l.out(LevFATAL, args...)
	os.Exit(3)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.outf(LevDEBUG, format, args...)
}

func (l *Logger) Verbosef(format string, args ...interface{}) {
	l.outf(LevVERBOSE, format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.outf(LevINFO, format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.outf(LevWARN, format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.outf(LevERROR, format, args...)
}

func (l *Logger) Notef(format string, args ...interface{}) {
	l.outf(LevNOTE, format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.printStack(true)
	l.outf(LevFATAL, format, args...)
	os.Exit(3)
}

func (l *Logger) Exitf(format string, args ...interface{}) {
	l.outf(LevFATAL, format, args...)
	os.Exit(3)
}

func Debug(args ...interface{}) {
	global.out(LevDEBUG, args...)
}

func Verbose(args ...interface{}) {
	global.out(LevVERBOSE, args...)
}

func Info(args ...interface{}) {
	global.out(LevINFO, args...)
}

func Warn(args ...interface{}) {
	global.out(LevWARN, args...)
}

func Error(args ...interface{}) {
	global.out(LevERROR, args...)
}

func Note(args ...interface{}) {
	global.out(LevNOTE, args...)
}

func Fatal(args ...interface{}) {
	global.printStack(true)
	global.out(LevFATAL, args...)
	os.Exit(3)
}

func Exit(args ...interface{}) {
	global.out(LevFATAL, args...)
	os.Exit(3)
}

func Debugf(format string, args ...interface{}) {
	global.outf(LevDEBUG, format, args...)
}

func Verbosef(format string, args ...interface{}) {
	global.outf(LevVERBOSE, format, args...)
}

func Infof(format string, args ...interface{}) {
	global.outf(LevINFO, format, args...)
}

func Warnf(format string, args ...interface{}) {
	global.outf(LevWARN, format, args...)
}

func Errorf(format string, args ...interface{}) {
	global.outf(LevERROR, format, args...)
}

func Notef(format string, args ...interface{}) {
	global.outf(LevNOTE, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	global.printStack(true)
	global.outf(LevFATAL, format, args...)
	os.Exit(3)
}

func Exitf(format string, args ...interface{}) {
	global.outf(LevFATAL, format, args...)
	os.Exit(3)
}

func Close() {
	global.Close()
}

func (s *Logger) Close() {
	for _, n := range s.lp {
		n.p.Close()
	}
}

func (s *Logger) out(lv lev, args ...interface{}) {
	if s.rootLev <= lv {
		str := s.f.Format(lv, args...)
		for _, n := range s.lp {
			if lv >= n.l {
				n.p.Print(str)
				if lv == LevFATAL {
					n.p.Close()
				}
			}
		}
	}
}

func (s *Logger) outf(lv lev, format string, args ...interface{}) {
	if s.rootLev <= lv {
		str := s.f.Formatf(lv, format, args...)
		for _, n := range s.lp {
			if lv >= n.l {
				n.p.Print(str)
				if lv == LevFATAL {
					n.p.Close()
				}
			}
		}
	}
}

func (s *Logger) printStack(all bool) {
	n := 500
	if all {
		n = 1000
	}
	var trace []byte

	for i := 0; i < 5; i++ {
		n *= 2
		trace = make([]byte, n)
		nbytes := runtime.Stack(trace, all)
		if nbytes <= len(trace) {
			n = nbytes
			break
		}
	}
	ms := string(trace[:n])
	for _, n := range s.lp {
		n.p.Print(&ms)
	}
}
