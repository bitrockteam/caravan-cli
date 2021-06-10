package commands

import (
	"fmt"

	"caravan/internal/aws"
	"caravan/internal/caravan"
	"caravan/internal/git"
)

type InitCommand struct {
	Name      string
	CloudName string
	Repos     []string
	Profile   string
	Workdir   string
}

func (i InitCommand) Help() string {
	return "init help\n"
}

func (i InitCommand) Synopsis() string {
	return `
	init: please provide the following args:
          - project name
          - infrastructure provider name: e.g. aws,gcp,azure
          - profile
          - region
	  e.g ./caravan init demo aws default eu-south-1
	  `
}

// Creates the needed supporting structures as follows:
//
//- clone git projects on .caravan/<project-name>
//
//- creates a bucket and lock table for terraform
//
//- generates the terrform files needed to deploy the infrastructure
func (i InitCommand) Run(args []string) (ret int) {

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
		Workdir:        i.Workdir,
		WorkdirProject: i.Workdir + "/" + args[0],
		BucketName:     args[0] + "-caravan-terraform-state",
		TableName:      args[0] + "-caravan-terraform-state-lock",
	}

	// checkout repos
	git := git.NewGit("bitrockteam")
	for _, repo := range i.Repos {
		if repo == "caravan-infra" {
			repo = repo + "-" + c.Provider
		}
		err := git.Clone(repo, ".caravan/"+i.Name+"/"+repo)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			ret = 1
		}
	}

	// init AWS
	err := initCloud(c)
	if err != nil {
		fmt.Printf("error during init: %s\n", err)
		return 1
	}

	return
}

func initCloud(c caravan.Config) (err error) {

	//generate configs and supporting items (bucket and locktable)
	fmt.Printf("initializing cloud resources\n")
	cloud := aws.NewAWS(c)

	fmt.Printf("creating bucket: %s\n", c.BucketName)
	err = cloud.CreateBucket(c.BucketName)
	if err != nil {
		return err
	}

	fmt.Printf("creating lock table: %s\n", c.TableName)
	err = cloud.CreateLockTable(c.TableName)
	if err != nil {
		return err
	}

	err = cloud.GenerateConfig()
	if err != nil {
		return err
	}
	return nil
}
