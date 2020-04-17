package main

import (
	"os"
	"testing"

	"github.com/fsnotify/fsnotify"
)

const (
	testContent = `
	this is a content this is a content this is a content this is a content
	this is a content this is a content this is a content this is a content
	this is a content this is a content this is a content this is a content
	this is a content this is a content this is a content this is a content
	this is a content this is a content this is a content this is a content
	this is a content this is a content this is a content this is a content
	this is a content this is a content this is a content this is a content
	this is a content this is a content this is a content this is a content
	this is a content this is a content this is a content this is a content
	`
)

func TestInvalidWatcherInstanciation(t *testing.T) {
	if _, err := newWatcher("", nil, nil); err == nil {
		t.Log("watcher instanciation should have failed, invalid path: received empty string")
		t.Fail()
	}

	if _, err := newWatcher("dummy", nil, nil); err == nil {
		t.Log("watcher instanciation should have failed, invalid events channel: received nil")
		t.Fail()
	}

	eventsChannel := make(chan *event)
	if _, err := newWatcher("dummy", eventsChannel, nil); err == nil {
		t.Log("watcher instanciation should have failed, invalid errors channel: received nil")
		t.Fail()
	}

	errorsChannel := make(chan error)
	if _, err := newWatcher("dummy", eventsChannel, errorsChannel); err == nil {
		t.Log("watcher instanciation should have failed, invalid path: does not exist or incorrect permissions")
		t.Fail()
	}
}

func processWatcherResult(op fsnotify.Op, ch chan *event, t *testing.T) {
	event, ok := <-ch

	if !ok || event == nil {
		t.Log("unexpected error occurred on event channel")
		t.Fail()
		return
	}

	if timestamp, ok := event.trace[op]; !ok || len(timestamp) == 0 {
		t.Logf("event %s should have been detected and timestamp should not be empty", op.String())
		t.Fail()
	}
}

func writeFile(path, content string, errchan chan error) {
	var f *os.File
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		errchan <- err
		return
	}

	if err := f.Sync(); err != nil {
		errchan <- err
		return
	}

	if _, err := f.WriteString(content); err != nil {
		errchan <- err
		return
	}

	if err := f.Sync(); err != nil {
		errchan <- err
		return
	}

	if err := f.Close(); err != nil {
		errchan <- err
		return
	}

	return
}

func renameFile(oldpath, newpath string, errchan chan error) {
	if err := os.Rename(oldpath, newpath); err != nil {
		errchan <- err
		return
	}

	return
}

func removeFile(path string, errchan chan error) {
	if err := os.Remove(path); err != nil {
		errchan <- err
		return
	}

	return
}
