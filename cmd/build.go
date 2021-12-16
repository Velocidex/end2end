package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Velocidex/ordereddict"
	end2end "github.com/Velocidex/velociraptor-endtoend"
	"github.com/alecthomas/kingpin"
)

var (
	build_cmd = app.Command("build", "build an initial datastore")

	app_datastore = app.Flag("datastore", "Datastore directory").
			Default("/tmp").String()

	app_os = app.Flag(
		"os", "OS to test").Default("linux").String()

	app_arch = app.Flag(
		"arch", "Arch to test").Default("amd64").String()

	app_release = app.Flag(
		"release", "Release to test").Default("v0.6.1").String()

	app_release_final = app.Flag(
		"final_release", "Release to test").Default("v0.6.3-rc1").String()
)

func setup(
	state *ordereddict.Dict,
	binary, datastore, config_file, api_client_config string) error {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer func() {
		cancel()
		time.Sleep(time.Second)
	}()

	err := end2end.StartServer(ctx, wg, binary, config_file)
	if err != nil {
		return err
	}

	fmt.Printf("Starting pool client\n")
	count := 20

	err = end2end.StartPool(ctx, wg, binary, config_file, datastore, count)
	if err != nil {
		return err
	}

	err = end2end.WaitForEnrolment(ctx, wg, binary, api_client_config, count)
	if err != nil {
		return err
	}

	err = end2end.AddLabels(ctx, state, binary, api_client_config)
	if err != nil {
		return err
	}

	fmt.Printf("State %v\n", state)

	/*
		err = end2end.CheckServer(
			ctx, state, binary, api_client_config)
		if err != nil {
			return err
		}
	*/
	return nil
}

func check_server(
	state *ordereddict.Dict,
	binary, datastore, config_file, api_client_config string) error {

	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	fmt.Printf("Starting server %v\n", binary)
	err := end2end.StartServer(ctx, wg, binary, config_file)
	if err != nil {
		fmt.Printf("StartServer: Got error %v\n", err)
		return err
	}

	fmt.Printf("Checking server %v\n", binary)
	err = end2end.CheckServer(ctx, state, binary, api_client_config)
	if err != nil {
		fmt.Printf("CheckServer: Got error %v\n", err)
		return err
	}

	return err
}

func doBuild() error {
	datastore := *app_datastore
	binary, err := end2end.GetRelease(
		*app_os, *app_arch, *app_release, datastore)
	if err != nil {
		return err
	}

	fmt.Printf("Loading %v\n", binary)
	config_file, err := end2end.GenerateConfig(binary, datastore)
	if err != nil {
		fmt.Printf("GenerateConfig %v\n", err)
		return err
	}

	api_client_config, err := end2end.GenerateAPIConfig(
		binary, "server", config_file, datastore)

	fmt.Printf("Config at %v and api_client_config %v\n",
		config_file, api_client_config)

	state := ordereddict.NewDict()

	err = setup(state, binary, datastore, config_file, api_client_config)
	if err != nil {
		return err
	}

	// Now spin up the new server.
	new_binary, err := end2end.GetRelease(
		*app_os, *app_arch, *app_release_final, datastore)
	if err != nil {
		return err
	}

	// Apply any known upgrade steps
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	err = end2end.Upgrade(ctx, *app_os, *app_arch, datastore,
		*app_release, *app_release_final)
	if err != nil {
		return err
	}

	err = check_server(state, new_binary, datastore, config_file, api_client_config)
	if err != nil {
		return fmt.Errorf("While checking version %v: %w", new_binary, err)
	}
	return nil
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		if command == build_cmd.FullCommand() {
			err := doBuild()
			time.Sleep(time.Second)
			kingpin.FatalIfError(err, "Checks")
			return true
		}
		return false
	})
}
