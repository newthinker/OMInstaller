package utl

import (
	"fmt"
	"time"
)

// get current time string
func NowString() string {
	return fmt.Sprint(time.Now().Format("2006-01-02 15:04:05 [MSG] "))
}
