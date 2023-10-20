package time_convert

import (
	"time"
	"fmt"
	"log/slog"
)

var CookieMaxAge int
var RedisTTL time.Duration

func ConvertTimeFormat(ttl string) {

	var (
		err error
		ttlDuration time.Duration
	)

	ttlDuration, err = time.ParseDuration(ttl)
	if err != nil {
		errorMessage := fmt.Sprintf("TTL string convert to duration failed: %v", err)
		slog.Error(errorMessage)
		panic(err)
	}

	CookieMaxAge = int(ttlDuration.Seconds())
	RedisTTL = ttlDuration
}