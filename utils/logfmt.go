package utils

import (
	"fmt"
	"log"
)

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	//colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	//colorBlue   = "\033[34m"
	//colorPurple = "\033[35m"
	//colorCyan   = "\033[36m"
	//colorWhite  = "\033[37m"
)

func Err(e error) {
	log.Printf("%s[ERR]%s\t%s", colorRed, colorReset, e)
}

func ErrString(m string) string {
	return fmt.Sprintf("%s[ERR]%s\t%s", colorRed, colorReset, m)
}

func ErrStringMsg(m string, e error) string {
	return fmt.Sprintf("%s[ERR]\t%s:%s %s", colorRed, colorReset, m, e)
}

func LogDeferError(s string, err error) {
	log.Printf("%sWARNING:%s\tCannot close %s, check for memory leak\terr:%v", colorYellow, colorReset, s, err)
}
