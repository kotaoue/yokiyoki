package formatter

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	return fmt.Sprintf("%dd %02dh %02dm", days, hours, minutes)
}
