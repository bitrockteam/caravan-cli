package git

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

type Git struct {
	org string
}

func NewGit(org string) (g Git) {
	return Git{org: org}
}

func (g Git) Clone(repo, dest string) (err error) {

	fmt.Printf("cloning repo %s/%s to %s\n", g.org, repo, dest)
	_, err = git.PlainClone("./"+dest, false, &git.CloneOptions{
		URL:      "https://github.com/" + g.org + "/" + repo,
		Progress: os.Stdout,
	})
	if err != nil {
		//TODO: improve error check
		if err.Error() != "repository already exists" {
			return fmt.Errorf("unable to clone repo %s: %s\n", repo, err)
		}
	}
	return nil
}
