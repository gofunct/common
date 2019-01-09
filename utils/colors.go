package utils

import (
	"github.com/fatih/color"
)

func Red(args ...interface{}) string {
	return color.RedString("%s", args)
}

func Green(args ...interface{}) string {
	return color.GreenString("%s", args)
}

func White(args ...interface{}) string {
	return color.WhiteString("%s", args)
}

func Black(args ...interface{}) string {
	return color.BlackString("%s", args)
}

func Magenta(args ...interface{}) string {
	return color.MagentaString("%s", args)
}

func Blue(args ...interface{}) string {
	return color.BlueString("%s", args)
}

func Cyan(args ...interface{}) string {
	return color.CyanString("%s", args)
}

func Yello(args ...interface{}) string {
	return color.YellowString("%s", args)
}
