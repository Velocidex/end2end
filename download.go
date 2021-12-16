package end2end

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Drop any post release numbers
func getBareRelease(release string) string {
	return strings.Split(release, "-")[0]
}

func GetRelease(test_os, arch, release, datastore string) (string, error) {
	filename := fmt.Sprintf("velociraptor-%s-%s-%s", release, test_os, arch)
	if runtime.GOOS == "windows" {
		filename += ".exe"
	}

	file_path := filepath.Join(datastore, "binaries", filename)

	fd, err := os.Open(file_path)
	if err != nil {
		url := fmt.Sprintf(
			"https://github.com/Velocidex/velociraptor/releases/download/%v/%v",
			getBareRelease(release), filename)

		fmt.Printf("File %v not found, will try to download from %v\n",
			filename, url)

		// Ensure the directory exists
		_ = os.MkdirAll(filepath.Join(datastore, "binaries"), 0700)

		resp, err := http.Get(url)
		if err != nil {
			return "", err
		}

		out, err := os.OpenFile(file_path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return "", err
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		return file_path, err
	}
	defer fd.Close()

	return file_path, err
}
