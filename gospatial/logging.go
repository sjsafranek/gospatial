package gospatial

import (
	"fmt"
	seelog "github.com/cihub/seelog"
)

var (
	Verbose       bool = false
	ServerLogger  seelog.LoggerInterface
	NetworkLogger seelog.LoggerInterface
	LogDirectory  string = "log"
	LogLevel      string = "trace"
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

func init() {
	DisableLog()
	loadServerConfig()
	loadNetworkConfig()
}

func ResetLogging() {
	DisableLog()
	loadServerConfig()
	loadNetworkConfig()
}

func enable_test_logging() {
	LogLevel = "error"
	DisableLog()
	loadServerConfig()
	loadNetworkConfig()
}

// DisableLog disables all library log output
func DisableLog() {
	NetworkLogger = seelog.Disabled
	ServerLogger = seelog.Disabled
}
