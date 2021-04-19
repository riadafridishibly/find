package walk

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"sync"
)

type walker struct {
	walkFunc filepath.WalkFunc
	wg       sync.WaitGroup
	jobs     chan string
	done     chan bool
}

func (w *walker) addJob(path string) {
	w.wg.Add(1)
	select {
	case w.jobs <- path:
	default:
		w.process(path)
	}
}

func (w *walker) process(path string) {
	defer w.wg.Done()

	fis, err := ioutil.ReadDir(path)

	if err != nil {
		return
	}

	for _, f := range fis {
		currentPath := filepath.Join(path, f.Name())
		w.walkFunc(currentPath, nil, nil)

		if f.IsDir() {
			w.addJob(currentPath)
		}
	}
}

func (w *walker) worker() {
	for job := range w.jobs {
		w.process(job)
	}
}

func (w *walker) walk(path string, walkFn filepath.WalkFunc) {
	for i := 0; i < 2*runtime.NumCPU(); i++ {
		go w.worker()
	}

	w.addJob(path)
	w.wg.Wait()
	close(w.jobs)

	w.done <- true
}

// Walk walks a directory concurrently and returns path in string
func Walk(root string, walkFn filepath.WalkFunc) <-chan bool {
	w := &walker{
		jobs:     make(chan string, 512),
		walkFunc: walkFn,
		done:     make(chan bool),
	}

	go w.walk(root, walkFn)

	return w.done
}
