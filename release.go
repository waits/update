package update

import (
	"fmt"
	"runtime"
)

type Release struct {
	Version Version `json:"version"`
	Assets  []Asset `json:"assets"`
}

func (rel Release) Install() error {
	asset, err := rel.assetForPlatform()
	if err != nil {
		return err
	}

	err = asset.apply()
	if err != nil {
		return err
	}

	return nil
}

func (rel Release) assetForPlatform() (a Asset, err error) {
	platform := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	for _, a = range rel.Assets {
		fmt.Println(a.Platform)
		if a.Platform == platform {
			return
		}
	}
	err = fmt.Errorf("no release for %s", platform)
	return
}
