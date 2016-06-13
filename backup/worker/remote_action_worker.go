package worker

import (
	"fmt"
	"sync"

	"github.com/ptrimble/dreamhost-personal-backup/backup"
)

type RemoteActionWorker struct {
	wg     *sync.WaitGroup
	in     <-chan backup.RemoteAction
	logger backup.BackupLogger

	putToRemote      func(string) error
	removeFromRemote func(string) error
}

func NewRemoteActionWorker(
	putToRemote, removeFromRemote func(string) error,
	wg *sync.WaitGroup,
	in <-chan backup.RemoteAction,
	log backup.BackupLogger,
) RemoteActionWorker {
	return RemoteActionWorker{
		putToRemote:      putToRemote,
		removeFromRemote: removeFromRemote,
		wg:               wg,
		in:               in,
		logger:           log,
	}
}

func (w RemoteActionWorker) Run() {
	for {
		action := <-w.in

		switch action.Type {
		case backup.PUSH:
			w.push(action.File)
		case backup.REMOVE:
			w.remove(action.File)
		}
	}
}

func (w RemoteActionWorker) push(file backup.File) {
	err := w.putToRemote(file.Name)
	if err != nil {
		w.logger.Error(backup.LogEntry{
			Message: fmt.Sprintf("unable to push to remote for file '%s', error: '%s'", file, err.Error()),
			File:    file.Name,
		})
	} else {
		w.logger.Info(backup.LogEntry{
			Message: fmt.Sprintf("%s pushed to remote", file),
			File:    file.Name,
		})
	}

	w.wg.Done()
}

func (w RemoteActionWorker) remove(file backup.File) {
	err := w.removeFromRemote(file.Name)
	if err != nil {
		entry := backup.LogEntry{
			Message: fmt.Sprintf("%s not found locally but unable to remove from remote, error: '%s'", file, err.Error()),
			File:    file.Name,
		}
		w.logger.Error(entry)
	} else {
		entry := backup.LogEntry{
			Message: fmt.Sprintf("%s not found locally, removing from remote", file),
			File:    file.Name,
		}
		w.logger.Info(entry)
	}

	w.wg.Done()
}
