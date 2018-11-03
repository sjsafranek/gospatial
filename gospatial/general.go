package gospatial

import (
	"time"
)

import (
	"./utils"
)

const (
	VERSION string = "1.11.4"
)

var (
	startTime           = time.Now()
	SuperuserKey string = utils.NewAPIKey(12)
)
