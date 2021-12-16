package main

import (
	"context"
	"path/filepath"
	"sync"
	"time"

	end2end "github.com/Velocidex/velociraptor-endtoend"
)

var (
	up_cmd = app.Command("up", "Bring test server up")
)

func doUp() error {
	datastore := *app_datastore
	binary, err := end2end.GetRelease(
		*app_os, *app_arch, *app_release, datastore)
	if err != nil {
		return err
	}

	config_file := filepath.Join(datastore, "server.config.yaml")

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer func() {
		cancel()
		time.Sleep(time.Second)
	}()

	err = end2end.StartServer(ctx, wg, binary, config_file)
	if err != nil {
		return err
	}
	time.Sleep(time.Hour)
	return nil
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		if command == up_cmd.FullCommand() {
			doUp()
			return true
		}
		return false
	})
}
