package schedule

import (
	"log"
	"os"
)

var stdLogger = log.New(os.Stdout, "scheduler: ", log.Ldate|log.LstdFlags|log.Lmsgprefix)
var errLogger = log.New(os.Stderr, "scheduler: ", log.Ldate|log.LstdFlags|log.Lmsgprefix)
