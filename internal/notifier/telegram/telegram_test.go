package telegram

import (
	"testing"
	"time"
)

func TestAverage(t *testing.T) {
	now := getNow()
	now2 := time.Now().In(getMoscowLocation())
	if now.Format("2006-01-02") != now2.Format("2006-01-02") {
		t.Errorf("Expected %s, got %s", now.Format("2006-01-02"), now2.Format("2006-01-02"))
	}
}
