package update

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

func Check(cur Version, url string) (rel Release, err error) {
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
	if rel.Version.after(cur) {
		err = errors.New("No update available")
	}
	return
}

func (v Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v Version) after(o Version) bool {
	if v.Major != o.Major {
		return v.Major > o.Major
	} else if v.Minor != o.Minor {
		return v.Minor > o.Minor
	}
	return v.Patch > o.Patch
}
