package logs

import (
    // "errors"
    "fmt"
    seelog "github.com/cihub/seelog"
    // "io"
)

var Logger seelog.LoggerInterface
var Network seelog.LoggerInterface

func loadAppConfig() {
    // https://github.com/cihub/seelog/wiki/Log-levels
    appConfig := `
<seelog minlevel="info">
    <outputs formatid="common">
        <rollingfile type="size" filename="log/test_roll.log" maxsize="100000" maxrolls="5"/>
        <filter levels="critical">
            <console formatid="stdout"/>
            <file path="log/test_critical.log" formatid="common"/>
        </filter>
        <filter levels="error">
            <console formatid="stdout"/>
            <file path="log/test_error.log" formatid="common"/>
        </filter>
        <filter levels="warn">
            <console formatid="stdout"/> 
            <file path="log/test_warn.log" formatid="common"/>
        </filter>
        <filter levels="info">
            <console formatid="stdout"/>
            <file path="log/test_info.log" formatid="common"/>
        </filter>
        <filter levels="debug">
            <console formatid="stdout"/>
            <file path="log/test_debug.log" formatid="common"/>
        </filter>
        <filter levels="trace">
            <console formatid="stdout"/>
            <file path="log/test_trace.log" formatid="common"/>
        </filter>
    </outputs>
    <formats>
        <format id="common"   format="%Date %Time [%LEVEL] %File %Func %Msg%n" />
        <format id="stdout"   format="%Date %Time [%LEVEL] %File %Func %Msg%n" />
    </formats>
</seelog>
`

    logger, err := seelog.LoggerFromConfigAsBytes([]byte(appConfig))
    if err != nil {
        fmt.Println(err)
        return
    }
    UseLogger(logger)
}

func init() {
    DisableLog()
    loadAppConfig()
}

// DisableLog disables all library log output
func DisableLog() {
    Logger = seelog.Disabled
}

// UseLogger uses a specified seelog.LoggerInterface to output library log.
// Use this func if you are using Seelog logging system in your app.
func UseLogger(newLogger seelog.LoggerInterface) {
    Logger = newLogger
}