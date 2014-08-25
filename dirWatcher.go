package main

import (
	"code.google.com/p/go.exp/fsnotify"
	"github.com/codegangsta/cli"
	"log"
	"os"
	"os/signal"
	"strings"
)

func processEvents(w *fsnotify.Watcher, e string) {
	for {
		select {
		case ev := <-w.Event:
			if e == "" {
				log.Println("event:", ev)
			} else {
				if strings.Contains(ev.String(), ".yaml\"") {
					log.Println("event:", ev)
				}
			}
		case err := <-w.Error:
			log.Println("error:", err)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "dwatcher"
	app.Usage = "Watch a directory for filesystem changes"

	app.Flags = []cli.Flag{
		cli.StringFlag{"directory, d", "/tmp/foo", "directory dwatcher watches"},
		cli.StringFlag{"extension, e", "", "only report files with this extension"},
	}

	app.Action = func(c *cli.Context) {
		dir := c.String("directory")
		ext := c.String("extension")

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}

		go processEvents(watcher, ext)

		err = watcher.Watch(dir)
		if err != nil {
			log.Fatal(err)
		}
		// kill event, ctrl+c
		onkill := make(chan os.Signal, 1)
		signal.Notify(onkill, os.Interrupt, os.Kill)
		<-onkill // wait for event

	}

	app.Run(os.Args)

}
