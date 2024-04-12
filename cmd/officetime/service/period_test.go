package service

import (
	"testing"
	"time"
)

func TestAverage(t *testing.T) {
	now := GetNow()
	now2 := time.Now().In(GetMoscowLocation())
	if now.Format("2006-01-02") != now2.Format("2006-01-02") {
		t.Errorf("Expected %s, got %s", now.Format("2006-01-02"), now2.Format("2006-01-02"))
	}
}
