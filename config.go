package end2end

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	BIND_PORT = 8020
)

func GenerateConfig(binary, datastore string) (string, error) {
	template := fmt.Sprintf(`
{"Datastore":{"location": %q, "filestore_directory": %q},
 "Frontend": {"bind_port": %d},
 "API": {"bind_port": %d},
 "Monitoring": {"bind_port": %d},
 "GUI": {"bind_port": %d},
 "Client": {
    "server_urls": [%q],
    "use_self_signed_ssl": true,
    "max_poll": 1
 }
}`, datastore, datastore, GetFrontendPort(),
		GetAPIPort(),
		GetMonitoringPort(),
		GetGUIPort(), GetFrontendUrl())

	fmt.Println(template)
	cmd := exec.Command(binary, "config", "generate", "--merge", template)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error %v\n", string(out))
		return "", err
	}

	filename := filepath.Join(datastore, "server.config.yaml")
	fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", err
	}

	_, err = fd.Write(out)
	return filename, err
}

func GenerateAPIConfig(binary, name, config_file, datastore string) (string, error) {
	api_config_file := filepath.Join(datastore, "api_config.yaml")
	cmd := exec.Command(binary, "--config", config_file, "config", "api_client",
		"--role=administrator", "--name", name, api_config_file)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error %v\n", string(out))
	}
	return api_config_file, err
}
