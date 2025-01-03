package constant

import "time"

const (
	SuccessCode = "0000"

	LocationTimezone = "Asia/Makassar"
)

var (
	DefaultTimezone, _ = time.LoadLocation(LocationTimezone)
)
