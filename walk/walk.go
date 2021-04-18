package walk

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"sync"
)

type filesHandler struct {
	ch         chan string
	filterFunc func(string) bool
}

func (fh *filesHandler) all(path string, wg *sync.WaitGroup, c chan bool) {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return
	}

	for _, f := range files {
		currentPath := filepath.Join(path, f.Name())
		if f.IsDir() {
			wg.Add(1)
			go func() {
				c <- true
				fh.all(currentPath, wg, c)
				<-c
				wg.Done()
			}()
		}

		if fh.filterFunc(f.Name()) {
			fh.ch <- currentPath
		}
	}
}

// WalkDir walks a directory concurrently and returns path in string
func WalkDir(root string, filterFunc func(string) bool) <-chan string {
	d := &filesHandler{
		ch:         make(chan string, 128),
		filterFunc: filterFunc,
	}

	go func() {
		wg := &sync.WaitGroup{}
		c := make(chan bool, runtime.NumCPU()*2)
		d.all(root, wg, c)
		wg.Wait()
		close(d.ch)
	}()

	return d.ch
}
