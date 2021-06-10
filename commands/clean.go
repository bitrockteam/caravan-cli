package commands

import (
	"caravan/internal/aws"
	"caravan/internal/caravan"
	"fmt"
	"os"
)

type CleanupCommand struct {
	Workdir string
	Repos   []string
}

func (c CleanupCommand) Help() string {
	return "removes all local repositories and terraform state and structures\n"
}

func (c CleanupCommand) Synopsis() string {
	return "Cleanup of carvan repositories and terraform locks and"
}

func (c CleanupCommand) Run(args []string) (ret int) {
	if len(args) != 4 {
		fmt.Printf(`init: please provide the following args:
                        - project name
                        - infrastructure provider name: e.g. aws,gcp,azure
                        - profile
                        - region`)
		return 1
	}
	cfg := caravan.Config{
		Name:     args[0],
		Provider: args[1],
		Profile:  args[2],
		Region:   args[3],
		Workdir:  c.Workdir,
	}

	os.RemoveAll(cfg.Workdir + "/" + cfg.Name)
	err := cleanCloud(cfg)
	if err != nil {
		return 1
	}
	return
}

func cleanCloud(cfg caravan.Config) (err error) {

	//generate configs and supporting items (bucket and locktable)
	fmt.Printf("removing terraform state and locking structures\n")

	cloud := aws.NewAWS(cfg)

	err = cloud.DeleteBucket(cfg.Name + "-caravan-terraform-state")
	if err != nil {
		return err
	}

	err = cloud.DeleteLockTable(cfg.Name + "-caravan-terraform-state-lock")
	if err != nil {
		return err
	}

	return nil
}
