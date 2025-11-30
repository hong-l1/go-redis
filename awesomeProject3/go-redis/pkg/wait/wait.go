package wait

import (
	"sync"
	"time"
)

type Wait struct {
	wait sync.WaitGroup
}

func (w *Wait) Add(delta int) {
	w.wait.Add(delta)
}
func (w *Wait) WaitWithTimeout(timeout time.Duration) bool {
	c := make(chan struct{}, 1)
	go func() {
		defer close(c)
		w.wait.Wait()
		c <- struct{}{}
	}()
	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}
func (w *Wait) Done() {
	w.wait.Done()
}
func (w *Wait) Wait() {
	w.wait.Wait()
}
