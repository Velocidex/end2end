package end2end

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

// Wait here until all the clients are enrolled.
func WaitForEnrolment(ctx context.Context,
	wg *sync.WaitGroup,
	binary, api_config_file string, count int) error {

	for {
		cmd := exec.CommandContext(ctx, binary, "--api_config",
			api_config_file, "query",
			"SELECT count() AS Total FROM clients() WHERE os_info.hostname group by 1")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error %v\n", string(out))
			return err
		}

		fmt.Printf("Checking clients %v\n", string(out))

		result_array := []map[string]interface{}{}
		err = json.Unmarshal(out, &result_array)
		if err != nil {
			return err
		}

		if len(result_array) > 0 {
			value := result_array[0]["Total"]
			fmt.Printf("Count %v clients ready...\n", value)
			if int(value.(float64)) == count {
				return nil
			}
		}
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(2 * time.Second):
		}
	}
}
