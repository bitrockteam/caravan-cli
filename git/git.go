// Git wraps access to git cli.
package git

import (
	"caravan-cli/cli"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Git struct {
	org      string
	logLevel string
}

func NewGit(org string, logLevel string) (g Git) {
	return Git{org: org, logLevel: logLevel}
}

func (g Git) Clone(name, dest, branch string) (err error) {
	log.Info().Msgf("cloning repo %s/%s to %s - branch: %s", g.org, name, dest, branch)
	cloneOptions := &git.CloneOptions{
		URL: "https://github.com/" + g.org + "/" + name,
	}
	if g.logLevel == cli.LogLevelDebug {
		cloneOptions.Progress = os.Stdout
	}
	repo, err := git.PlainClone("./"+dest, false, cloneOptions)
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
