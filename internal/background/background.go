package background

import "time"

type Background struct {
	delay time.Duration
	task  func()
}

func NewBackground(delay time.Duration, task func()) *Background {
	return &Background{delay: delay, task: task}
}

func (b Background) Start() {
	ticker := time.NewTicker(b.delay)
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				b.task()
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()

	return
}
