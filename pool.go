package end2end

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func StartPool(ctx context.Context,
	wg *sync.WaitGroup,
	binary, config_file, datastore string, count int) error {

	pool_path := filepath.Join(datastore, "pool")
	err := os.MkdirAll(pool_path, 0700)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, binary, "--config",
		config_file, "pool_client", "--writeback_dir",
		pool_path, "--number", fmt.Sprintf("%d", count), "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		cmd.Wait()
	}()

	return nil
}
