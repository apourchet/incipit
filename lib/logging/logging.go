package logging

import (
	"fmt"
	"log"
	"os"
)

var (
	INFO  *log.Logger
	ERROR *log.Logger
)

func init() {
	INFO = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	ERROR = log.New(os.Stdout, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(s string, args ...interface{}) {
	INFO.Output(2, fmt.Sprintf(s, args...))
}

func Error(s string, args ...interface{}) {
	ERROR.Output(2, fmt.Sprintf(s, args...))
}

func Fatal(args ...interface{}) {
	ERROR.Fatal(args)
}

func FatalError(err error) {
	if err != nil {
		Fatal(err)
	}
}
