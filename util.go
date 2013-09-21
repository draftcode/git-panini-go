package main

import (
	"fmt"
	"github.com/aybabtme/color/brush"
	"log"
	"os/user"
	"strings"
)

func fixPath(path string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	if strings.HasPrefix(path, "~/") {
		return strings.Replace(path, "~", usr.HomeDir, 1)
	} else {
		return path
	}
}

func printRed(format string, a ...interface{}) (n int, err error) {
	n, err = fmt.Printf(brush.Red(format).String(), a...)
	return
}

func printGreen(format string, a ...interface{}) (n int, err error) {
	n, err = fmt.Printf(brush.Green(format).String(), a...)
	return
}
