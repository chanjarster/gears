package schedule

import (
	"log"
	"os"
)

var stdLogger = log.New(os.Stdout, "", log.Ldate|log.LstdFlags|log.Llongfile)
var errLogger = log.New(os.Stderr, "", log.Ldate|log.LstdFlags|log.Llongfile)
