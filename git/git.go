// Git wraps access to git cli.
package git

import (
	"fmt"
	"os"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Git struct {
	org string
}

func NewGit(org string) (g Git) {
	return Git{org: org}
}

func (g Git) Clone(name, dest, branch string) (err error) {
	fmt.Printf("cloning repo %s/%s to %s - branch: %s\n", g.org, name, dest, branch)

	repo, err := git.PlainClone("./"+dest, false, &git.CloneOptions{
		URL:      "https://github.com/" + g.org + "/" + name,
		Progress: os.Stdout,
	})
	if err != nil {
		if err.Error() != "repository already exists" {
			return fmt.Errorf("unable to clone repo %s: %w", name, err)
		}
		repo, err = git.PlainOpen("./" + dest)
		if err != nil {
			return fmt.Errorf("unable to open repo %s: %w", name, err)
		}
	}

	if branch == "" {
		return nil
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("error getting worktree: %w", err)
	}

	b := plumbing.NewRemoteReferenceName("origin", branch)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: b,
	})
	if err != nil {
		// TODO better error check
		if !strings.Contains(err.Error(), "worktree contains unstaged changes") {
			return fmt.Errorf("error checking out: %w", err)
		}
	}

	return nil
}
