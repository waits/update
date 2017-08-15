package update

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const execMode = 0755
const urlFmt = "https://api.github.com/repos/%s/releases/latest"

func Apply(repo string, cur string) {
	url := fmt.Sprintf(urlFmt, repo)
	rel, err := getLatestRelease(url)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if rel.Tag == cur {
		fmt.Println("No update available.")
		return
	}

	asset, err := rel.assetForPlatform(runtime.GOOS, runtime.GOARCH, repo)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	err = applyPatch(asset)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func getLatestRelease(url string) (rel release, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	json.Unmarshal(dat, &rel)
	return
}

func applyPatch(asset asset) error {
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
