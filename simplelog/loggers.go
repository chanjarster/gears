package simplelog

import (
	"log"
	"os"
)

var StdLogger = log.New(os.Stdout, "", log.Ldate|log.LstdFlags|log.Llongfile)
var ErrLogger = log.New(os.Stderr, "", log.Ldate|log.LstdFlags|log.Llongfile)
