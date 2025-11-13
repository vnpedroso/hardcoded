package main

import (
	"fmt"
	"log"
	"math/rand"
)

const (
	ResetColor string = "\033[0m"
	Red        string = "\033[0;31m"
	Yellow     string = "\033[0;33m"
	Green      string = "\033[0;32m"
	Cyan       string = "\033[0;36m"
)

type coloredLogger struct {
	logger *log.Logger
	colors map[string]string
}

func (cl coloredLogger) ColoredPrintf(color string, format string, a ...any) {
	val, ok := cl.colors[color]

	if !ok {
		cl.logger.Printf("%s is an invalid color string", val)
		cl.logger.Printf(format, a...)
		return
	}

	base := fmt.Sprintf(format, a...)
	cl.logger.Printf("%s%s%s", val, base, ResetColor)
}

func NotSoRandomSuccess() bool {
	roulette := rand.Intn(3)
	return roulette < 2 //2 out of 3
}
