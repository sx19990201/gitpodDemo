package utils

import "sync"

type Worker struct {
	Source chan func()
}

func (w *Worker) Start() {
	w.Source = make(chan func())
	go func() {
		for {
			select {
			case task := <-w.Source:
				go task()
			}
		}
	}()
}

var tss *ThreadSafeSlice
var once sync.Once

func GetTss() *ThreadSafeSlice {
	if tss == nil {
		once.Do(func() {
			tss = &ThreadSafeSlice{}
		})
	}
	return tss
}

func InitWorks() {
	tts := GetTss()
	w := &Worker{}
	go w.Start()
	tts.Push(w)
}

type ThreadSafeSlice struct {
	sync.Mutex
	workers []*Worker
}

func (slice *ThreadSafeSlice) Push(w *Worker) {
	slice.Lock()
	defer slice.Unlock()

	slice.workers = append(slice.workers, w)
}

func (slice *ThreadSafeSlice) Process(routine func(*Worker)) {
	slice.Lock()
	defer slice.Unlock()

	for _, worker := range slice.workers {
		routine(worker)
	}
}
