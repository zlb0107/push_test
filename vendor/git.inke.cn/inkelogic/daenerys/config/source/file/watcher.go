package file

import (
	"errors"
	"git.inke.cn/inkelogic/daenerys/config/source"
	"github.com/fsnotify/fsnotify"
	"os"
)

type watcher struct {
	f *file

	fw   *fsnotify.Watcher
	exit chan bool
}

//当修改文件名称时，fsnotify中event.Name仍然是原来的文件名，这就需要我们在重命名事件中，先移除之前的监控，然后添加新的监控
func newWatcher(f *file) (source.Watcher, error) {
	//新建一个监听
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	//添加监听对象
	fw.Add(f.path)

	return &watcher{
		f:    f,
		fw:   fw,
		exit: make(chan bool),
	}, nil
}

func (w *watcher) Next() (*source.ChangeSet, error) {
	// is it closed?
	select {
	case <-w.exit:
		return nil, errors.New("watcher stopped")
	default:
	}

	// try get the event
	select {
	case event := <-w.fw.Events:
		switch event.Op {
		case fsnotify.Remove:
			return nil, nil
		case fsnotify.Rename:
			// check existence of file, and add watch again
			_, err := os.Stat(event.Name)
			if err == nil || os.IsExist(err) {
				w.fw.Add(event.Name)
			}
		default:
		}
		c, err := w.f.Read()
		if err != nil {
			return nil, err
		}
		return c, nil

	case err := <-w.fw.Errors:
		return nil, err
	case <-w.exit:
		return nil, errors.New("watcher stopped")
	}
}

func (w *watcher) Stop() error {
	return w.fw.Close()
}
