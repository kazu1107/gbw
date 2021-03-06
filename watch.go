package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func (app *App) watchFiles() {
	chExit := make(chan bool, 1)
	go app.commandRunner(chExit)

	if !app.waitFirst {
		app.chFileEvent <- FileEvent{
			action: "init",
		}
	}

	for {
		select {
		case event := <-app.watcher.Events:
			if len(app.exts) != 0 && !searchExt(event.Name, app.exts) {
				break
			}

			cmdRun := false

			var fEvent FileEvent

			switch {
			case event.Op&fsnotify.Write == fsnotify.Write && app.writed:
				fEvent.action = "writed"
				cmdRun = true

			case event.Op&fsnotify.Create == fsnotify.Create && app.created:
				fEvent.action = "created"
				cmdRun = true

			case event.Op&fsnotify.Remove == fsnotify.Remove && app.removed:
				fEvent.action = "removed"
				cmdRun = true

			case event.Op&fsnotify.Rename == fsnotify.Rename && app.renamed:
				fEvent.action = "renamed"
				cmdRun = true

			case event.Op&fsnotify.Chmod == fsnotify.Chmod && app.chmod:
				fEvent.action = "chmod"
				cmdRun = true
			}

			if cmdRun {
				fEvent.file = event.Name
				app.chDirEvent <- fEvent
				app.chFileEvent <- fEvent
			}
		case <-chExit:
			app.chDone <- true

		case err := <-app.watcher.Errors:
			if err != nil {
				log.Println(err)
			}
			app.chDone <- true
		}
	}
}
