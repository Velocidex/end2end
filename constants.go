package end2end

import "fmt"

const (
	BASE_PORT = 8010
)

func GetFrontendPort() int {
	return BASE_PORT
}

func GetFrontendUrl() string {
	return fmt.Sprintf("https://127.0.0.1:%d/", GetFrontendPort())
}

func GetAPIPort() int {
	return BASE_PORT + 1
}

func GetGUIPort() int {
	return BASE_PORT + 2
}

func GetMonitoringPort() int {
	return BASE_PORT + 3
}
