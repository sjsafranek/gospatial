package app

import (
    "fmt"
    seelog "github.com/cihub/seelog"
)

var (
    Verbose bool = false
    ServerLogger seelog.LoggerInterface
    NetworkLogger seelog.LoggerInterface
    DbLogger seelog.LoggerInterface
    LogDirectory string = "log"
    LogLevel string = "trace"
)

func loadServerConfig() {
    // https://github.com/cihub/seelog/wiki/Log-levels
    appConfig := `
<seelog minlevel="` + LogLevel + `">
    <outputs formatid="common">
        <filter levels="critical,error,warn">
            <console formatid="stdout"/>
            <file path="` + LogDirectory + `/error.log" formatid="common"/>
        </filter>
        <filter levels="info,debug,trace">
            <console formatid="stdout"/>
            <file path="` + LogDirectory + `/server.log" formatid="common"/>
        </filter>
    </outputs>
    <formats>
        <format id="common"   format="%Date %Time [%LEVEL] %File %FuncShort:%Line %Msg %n" />
        <format id="stdout"   format="%Date %Time [%LEVEL] %File %FuncShort:%Line %Msg %n" />
    </formats>
</seelog>
`

    logger, err := seelog.LoggerFromConfigAsBytes([]byte(appConfig))
    if err != nil {
        fmt.Println(err)
        return
    }
    // UseLogger(logger)
    ServerLogger = logger
}

func loadNetworkConfig() {
    // https://github.com/cihub/seelog/wiki/Log-levels
    appConfig := `
<seelog minlevel="` + LogLevel + `">
    <outputs formatid="common">
        <rollingfile type="size" filename="` + LogDirectory + `/network.log" maxsize="100000" maxrolls="5"/>
        <filter levels="critical,error,warn,info,debug,trace">
            <console formatid="stdout"/>
        </filter>
    </outputs>
    <formats>
        <format id="common"   format="%Date %Time [%LEVEL] %File %FuncShort:%Line %Msg %n" />
        <format id="stdout"   format="%Date %Time [%LEVEL] %File %FuncShort:%Line %Msg %n" />
    </formats>
</seelog>
`

    logger, err := seelog.LoggerFromConfigAsBytes([]byte(appConfig))
    if err != nil {
        fmt.Println(err)
        return
    }
    NetworkLogger = logger
}

func loadDbConfig() {
    // https://github.com/cihub/seelog/wiki/Log-levels
    appConfig := `
<seelog minlevel="` + LogLevel + `">
    <outputs formatid="common">
        <rollingfile type="size" filename="` + LogDirectory + `/db.log" maxsize="100000" maxrolls="5"/>
        <filter levels="critical,error,warn,info,debug,trace">
            <console formatid="stdout"/>
        </filter>
    </outputs>
    <formats>
        <format id="common"   format="%Date %Time [%LEVEL] %File %FuncShort:%Line %Msg %n" />
        <format id="stdout"   format="%Date %Time [%LEVEL] %File %FuncShort:%Line %Msg %n" />
    </formats>
</seelog>
`
    logger, err := seelog.LoggerFromConfigAsBytes([]byte(appConfig))
    if err != nil {
        fmt.Println(err)
        return
    }
    DbLogger = logger
}


func init() {
    DisableLog()
    loadServerConfig()
    loadNetworkConfig()
    loadDbConfig()
}

func ResetLogging() {
    DisableLog()
    loadServerConfig()
    loadNetworkConfig()
    loadDbConfig()
}


func enable_test_logging() {
    LogLevel = "error"
    DisableLog()
    loadServerConfig()
    loadNetworkConfig()
    loadDbConfig()
}

// DisableLog disables all library log output
func DisableLog() {
    NetworkLogger = seelog.Disabled
    ServerLogger = seelog.Disabled
    DbLogger = seelog.Disabled
}

// UseLogger uses a specified seelog.LoggerInterface to output library log.
// Use this func if you are using Seelog logging system in your app.
// func UseLogger(newLogger seelog.LoggerInterface) {
//     ServerLogger = newLogger
// }







// import (
// 	// "io"
// 	// "io/ioutil"
// 	"log"
// 	"os"
// 	// "path/filepath"
// 	// "strings"
// )

// var (
// 	Info *log.Logger
// 	Debug *log.Logger
// 	Warning *log.Logger
// 	Error                *log.Logger
// 	// serverLoggerWriter   io.Writer
// 	// networkLoggerWriter  io.Writer
// 	// networkLoggerInfo    *log.Logger
// 	// networkLoggerWarning *log.Logger
// 	// networkLoggerError   *log.Logger
// 	// networkLoggerInfoIn  *log.Logger
// 	// networkLoggerInfoOut *log.Logger
// 	Verbose bool = false
// )

// func init() {
// 	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
// 	// if err != nil {
// 	// 	Error.Fatal(err)
// 	// }

// 	// server logging
// 	// serverLogFile := strings.Replace(dir, "bin", "log/server.log", -1)
// 	// serverLoggerWriter, err := os.OpenFile(serverLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
// 	// if err != nil {
// 	// 	serverLoggerWriter, err = os.OpenFile("test.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
// 	// }
// 	// Info = log.New(serverLoggerWriter, "INFO  [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	// Debug = log.New(serverLoggerWriter, "DEBUG [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	// Warning = log.New(serverLoggerWriter, "WARN  [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	// Error = log.New(serverLoggerWriter, "ERROR [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)

// 	// network logging
// 	// networkLogFile := strings.Replace(dir, "bin", "log/network.log", -1)
// 	// networkLoggerWriter, err := os.OpenFile(networkLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
// 	// if err != nil {
// 	// 	networkLoggerWriter, err = os.OpenFile("test.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
// 	// }
// 	// networkLoggerInfoIn = log.New(networkLoggerWriter, "INFO  [NETWORK] [IN] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	// networkLoggerInfoOut = log.New(networkLoggerWriter, "INFO  [NETWORK] [OUT] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	// networkLoggerInfo = log.New(networkLoggerWriter, "INFO  [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	// networkLoggerWarning = log.New(networkLoggerWriter, "WARN  [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	// networkLoggerError = log.New(networkLoggerWriter, "ERROR [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// }

// func testLoggerInit() {
// 	serverLoggerWriter, err := os.OpenFile("test.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
// 	if err != nil {
// 		Error.Fatal("Error opening file: %v", err)
// 	}
// 	// defer serverLoggerWriter.Close()
// 	Info = log.New(serverLoggerWriter, "[TESTING] INFO  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	Debug = log.New(serverLoggerWriter, "[TESTING] DEBUG | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	Warning = log.New(serverLoggerWriter, "[TESTING] WARN  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	Error = log.New(serverLoggerWriter, "[TESTING] ERROR | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// }

// // func StdOutMode() {
// // 	AppMode = "standard"
// // 	Info = log.New(os.Stdout, "INFO  [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// // 	Debug = log.New(os.Stdout, "DEBUG [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// // 	Warning = log.New(os.Stdout, "WARN  [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// // 	Error = log.New(os.Stderr, "ERROR [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// // 	networkLoggerInfo = log.New(os.Stdout, "INFO  [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// // 	networkLoggerWarning = log.New(os.Stdout, "WARN  [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// // 	networkLoggerError = log.New(os.Stderr, "ERROR [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// // 	networkLoggerInfoIn = log.New(os.Stdout, "INFO  [NETWORK] [IN] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// // 	networkLoggerInfoOut = log.New(os.Stdout, "INFO  [NETWORK] [OUT] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// // }
