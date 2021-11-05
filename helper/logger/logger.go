package logger

import (
	"fmt"
	"io"
	"os"
)

// Tee does the same thing as UNIX `tee`
func Tee(s, filePath string) error {
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	mw.Write([]byte(fmt.Sprintf("%s\n", s)))
	return nil
}
