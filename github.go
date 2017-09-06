package update

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ghRelease struct {
	Tag    string    `json:"tag_name"`
	Assets []ghAsset `json:"assets"`
}

type ghAsset struct {
	Name string `json:"name"`
	Url  string `json:"browser_download_url"`
}

const tagFmt = "v%d.%d.%d"

func CheckGithub(cur Version, url string) (rel Release, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var gr ghRelease
	json.Unmarshal(dat, &gr)
	rel = gr.toRelease()

	if !rel.Version.after(cur) {
		err = errors.New("no update available")
	}
	return
}

func (gr ghRelease) toRelease() (rel Release) {
	var major, minor, patch int
	fmt.Sscanf(gr.Tag, tagFmt, &major, &minor, &patch)
	rel.Version = Version{major, minor, patch}

	for _, ga := range gr.Assets {
		names := strings.Split(ga.Name, "-")
		platform := fmt.Sprintf("%s/%s", names[1], names[2])
		rel.Assets = append(rel.Assets, Asset{ga.Name, platform, ga.Url})
	}

	return
}
