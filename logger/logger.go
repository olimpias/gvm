package logger

import "os"

func Info(text string) {
	os.Stdout.WriteString(text)
}

func ExitWithInfo(text string) {
	os.Stdout.WriteString(text)
	os.Exit(0)
}

func Error(text string) {
	os.Stderr.WriteString(text)
}

func ExitWithError(text string) {
	os.Stderr.WriteString(text)
	os.Exit(1)
}
