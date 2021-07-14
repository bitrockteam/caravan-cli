package terraform

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Terraform struct {
	Workdir string
}

func (t *Terraform) Init(wd string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	t.Workdir = wd
	fmt.Printf("running init on workdir: %s\n", t.Workdir)
	cmd := exec.CommandContext(ctx, "terraform", "init")
	cmd.Dir = t.Workdir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (t Terraform) ApplyVarMap(config map[string]string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	args := []string{}
	args = append(args, "apply")
	args = append(args, "-auto-approve")
	for k, v := range config {
		args = append(args, fmt.Sprintf("-var=%s=%s", k, v))
	}
	fmt.Printf("running apply on workdir: %s with args: %s\n", t.Workdir, args)
	cmd := exec.CommandContext(ctx, "terraform", args...)
	cmd.Dir = t.Workdir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (t Terraform) ApplyVarFile(file string, timeout time.Duration, env map[string]string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	args := []string{}
	args = append(args, "apply")
	args = append(args, "-auto-approve")
	args = append(args, "-var-file="+file)
	fmt.Printf("running apply on workdir: %s with args: %s\n", t.Workdir, args)
	cmd := exec.CommandContext(ctx, "terraform", args...)
	cmd.Dir = t.Workdir
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (t Terraform) Destroy(file string, env map[string]string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	args := []string{}
	args = append(args, "destroy")
	args = append(args, "-auto-approve")
	args = append(args, "-var-file="+file)
	fmt.Printf("running destroy on workdir: %s with args: %s\n", t.Workdir, args)
	cmd := exec.CommandContext(ctx, "terraform", args...)
	cmd.Dir = t.Workdir
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
