package commands

import (
	"fmt"
	"os"

	"caravan/internal/caravan"
	"caravan/internal/terraform"
)

type BakeVMCommand struct {
	Workdir string
}

func (b BakeVMCommand) Help() string {
	return "Bake a VM image on the given cloud provider\n"
}

func (b BakeVMCommand) Synopsis() string {
	return "Bake VM images"
}

func (b BakeVMCommand) Run(args []string) (ret int) {
	if len(args) != 4 {
		fmt.Printf(`init: please provide the following args:
                        - project name
                        - infrastructure provider name: e.g. aws,gcp,azure
                        - profile
                        - region`)
		return 1
	}

	c := caravan.Config{
		Name:           args[0],
		Provider:       args[1],
		Profile:        args[2],
		Region:         args[3],
		Workdir:        b.Workdir,
		WorkdirProject: b.Workdir + "/" + args[0],
		BucketName:     args[0] + "-caravan-terraform-state",
		TableName:      args[0] + "-caravan-terraform-state-lock",
	}

	if _, err := os.Stat(c.WorkdirProject); os.IsNotExist(err) {
		fmt.Printf("please run init before bake \n")
		return 1
	}

	tf := terraform.NewTerraform(c.Workdir + "/caravan-baking/terraform")
	err := tf.Init()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return 1
	}

	err = tf.ApplyVarFile(c.Provider + "-baking.tfvars")
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return 1
	}
	return
}
