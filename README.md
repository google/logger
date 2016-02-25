# logger #
Logger is a simple cross platform Go logging library for Windows and Linux, it 
can log to the Windows event log, Linux syslog, and an io.Writer. 

This is not an official Google product.

## Usage ##

Set up logger to log the system log (event log or syslog) and a file, 
include a flag to turn up verbosity:

```go
import (
  "flag"
  "os"
  
  "github.com/google/logger"
)

const logPath = "/some/location/example.log"

var verbose = flag.Bool("verbose", false, "print info level logs to stdout")

func main() {
  flag.Parse()

  lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
  if err != nil {
    logger.Fatalf("Failed to open log file: %v", err)
  }
  defer lf.Close()

  logger.Init("LoggerExample", *verbose, true, lf)

  logger.Info("I'm about to do something!")
  if err := doSomething(); err != nil {
    logger.Errorf("Error running doSomething: %v", err)
  }
}
```

