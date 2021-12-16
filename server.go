package end2end

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func StartServer(ctx context.Context,
	wg *sync.WaitGroup, binary, config_file string) error {

	argv := []string{"--config", config_file, "frontend", "-v"}
	cmd := exec.CommandContext(ctx, binary, argv...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Running %v %v\n", binary, strings.Join(argv, " "))

	err := cmd.Start()
	if err != nil {
		fmt.Printf("Got error %v\n", err)
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		cmd.Wait()
	}()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	for {
		resp, err := client.Get(GetFrontendUrl() + "server.pem")
		if err == nil {
			defer resp.Body.Close()

			if resp.StatusCode == 200 {
				return nil
			}
		}
		time.Sleep(time.Second)
	}

	return nil
}
