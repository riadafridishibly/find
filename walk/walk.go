package walk

import (
	"io/ioutil"
	"path/filepath"
	"sync"
)

type dirGetter struct {
	ch chan string
}

func (d *dirGetter) all(path string, wg *sync.WaitGroup, c chan bool) {
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
				d.all(currentPath, wg, c)
				<-c
				wg.Done()
			}()
		}

		d.ch <- currentPath
	}
}

// WalkDir walks a directory concurrently and returns path in string
func WalkDir(root string) <-chan string {
	d := &dirGetter{ch: make(chan string, 100)}
	go func() {
		wg := &sync.WaitGroup{}
		c := make(chan bool, 8)
		d.all(root, wg, c)
		wg.Wait()
		close(d.ch)
	}()

	return d.ch
}
