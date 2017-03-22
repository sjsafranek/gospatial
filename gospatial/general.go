package gospatial

import (
	"time"
)

import (
	"./utils"
)

const (
	VERSION string = "1.11.3"
)

var (
	startTime           = time.Now()
	SuperuserKey string = utils.NewAPIKey(12)
)
