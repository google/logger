# logger #
Logger is a simple cross platform Go logging library for Windows, Linux, and
macOS, it can log to the Windows event log, Linux/macOS syslog, and an io.Writer.

This is not an official Google product.

## Usage ##

Set up the default logger to log the system log (event log or syslog) and a
file, include a flag to turn up verbosity:

```go
import (
  "flag"
  "os"

  "github.com/google/logger"
)

const logPath = "/some/location/example.log"

func main() {
  lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
  if err != nil {
    panic("Failed to open log file: %v", err)
  }
  defer lf.Close()

  l := logger.NewLogger("LoggerExample", lf)

  l.Info("I'm about to do something!")
  if err := doSomething(); err != nil {
    l.Errorf("Error running doSomething: %v", err)
  }
}
```

The `NewLogger()` function returns a logger object, you may substantiate
multiple instances if you require multiple logging backends.

```go
lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
if err != nil {
  panic("Failed to open log file: %v", err)
}
defer lf.Close()

// Log to system log and a log file, Info logs don't write to stdout.
lZero := logger.NewLogger("LoggerExample0", lf)
// Don't send to system log or a log file, Info logs writes to stdout..
lOne := logger.NewLogger("LoggerExample1", ioutil.Discard)

lZero.Info("This will log to the log file and the system log")
lOne.Info("This will only log to stdout")
l.Info("This is the same as using loggerOne")

```
