package update

import "fmt"

func Auto(cur Version, url string, check func(Version, string) (Release, error)) {
	rel, err := check(cur, url)
	if err != nil {
		fmt.Println("No update available.")
		return
	}

	err = rel.Install()
	if err != nil {
		fmt.Println("Update failed:", err)
		return
	}

	fmt.Printf("Updated to %s.\n", rel.Version)
}
