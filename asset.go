package update

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const execMode = 0755

type Asset struct {
	Name     string `json:"name"`
	Platform string `json:"platform"`
	Url      string `json:"url"`
}

func (asset Asset) apply() error {
	path := filepath.Join("/", "tmp", asset.Name)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	defer os.Remove(path)

	err = file.Chmod(execMode)
	if err != nil {
		return err
	}

	fmt.Println("Downloading", asset.Url)
	resp, err := http.Get(asset.Url)
	if err != nil {
		return err
	}
	io.Copy(file, resp.Body)

	exec, err := os.Executable()
	if err != nil {
		return err
	}

	os.Rename(exec, exec+".old")
	os.Rename(path, exec)
	os.Remove(exec + ".old")

	return nil
}
