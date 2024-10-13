package background

import "time"

// Background Сервис для фонового выполнения функции через определенные интервалы
type Background struct {
	delay time.Duration
	task  func()
}

// NewBackground Конструктор сервиса для фонового выполнения функций
func NewBackground(delay time.Duration, task func()) *Background {
	return &Background{delay: delay, task: task}
}

// Start Выполнить функцию
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
