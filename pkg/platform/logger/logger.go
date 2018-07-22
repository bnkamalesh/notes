// Package logger provides multi-level log functionality
package logger

import (
	"log"
	"os"
)

const (
	// Svr1 defines severity level 1, highest level.
	Svr1 = severity("Fatal:")
	// Svr2 defines severity level 2, medium level.
	Svr2 = severity("Error:")
	// Svr3 defines severity level 3, low level.
	Svr3 = severity("Warning:")
	// Svr4 defines severity level 4, Info.
	Svr4 = severity("Info:")
	// Svr5 defines severity level 5, Debug.
	Svr5 = severity("Debug:")
)

var (
	svr1Logger = log.New(os.Stderr, string(Svr1), log.LstdFlags|log.Llongfile)
	svr2Logger = log.New(os.Stderr, string(Svr2), log.LstdFlags|log.Llongfile)
	svr3Logger = log.New(os.Stderr, string(Svr3), log.LstdFlags|log.Llongfile)
	svr4Logger = log.New(os.Stdout, string(Svr4), log.LstdFlags|log.Llongfile)
	svr5Logger = log.New(os.Stdout, string(Svr5), log.LstdFlags|log.Llongfile)
)

type severity string

// Service defines all the logging methods to be implemented
type Service interface {
	Debug(data ...interface{})
	Info(data ...interface{})
	Warn(data ...interface{})
	Error(data ...interface{})
	Fatal(data ...interface{})
}

// Log handles all the dependencies for logger
type Log struct{}

// New returns an instance of Log with all the dependencies initialized
func New() Log {
	return Log{}
}

// Debug prints log of severity 5
func (l Log) Debug(data ...interface{}) {
	svr5Logger.Println(data...)
}

// Info prints logs of severity 4
func (l Log) Info(data ...interface{}) {
	svr4Logger.Println(data...)
}

// Warn prints log of severity 3
func (l Log) Warn(data ...interface{}) {
	svr3Logger.Println(data...)
}

//  Error prints log of severity 2
func (l Log) Error(data ...interface{}) {
	svr2Logger.Println(data...)
}

// Fatal prints log of severity 1
func (l Log) Fatal(data ...interface{}) {
	svr1Logger.Println(data...)
}
