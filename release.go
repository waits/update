package update

import "fmt"
import "strings"

type release struct {
	Name   string  `json:"name"`
	Tag    string  `json:"tag_name"`
	Assets []asset `json:"assets"`
}

type asset struct {
	Name string `json:"name"`
	Url  string `json:"browser_download_url"`
}

func (rel *release) assetForPlatform(os string, arch string, repo string) (a asset, err error) {
	repoName := strings.Split(repo, "/")[1]
	assetName := fmt.Sprintf("%s-%s-%s", repoName, os, arch)
	for _, a = range rel.Assets {
		if a.Name == assetName {
			return
		}
	}
	err = fmt.Errorf("no release for %s/%s", os, arch)
	return
}
