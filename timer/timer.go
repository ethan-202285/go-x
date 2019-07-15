package timer

import (
	"log"
	"time"
)

// New Timer
func New(name string) *Timer {
	return &Timer{name, time.Now()}
}

// Timer 计时器
type Timer struct {
	name      string
	timestamp time.Time
}

// Count 计时
func (t *Timer) Count(label string) {
	now := time.Now()
	if len(label) > 0 {
		log.Printf(
			"\033[1;35m%s\033[0m %s: %s",
			t.name,
			label,
			now.Sub(t.timestamp),
		)
	}
	t.timestamp = now
}
