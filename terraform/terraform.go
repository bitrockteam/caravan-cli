// Terraform wraps acces to terrraform cli utility.
package terraform

import (
	"caravan-cli/cli"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
)

type Terraform struct {
	Workdir  string
	logLevel string
}

func New(logLevel string) (t *Terraform) {
	return &Terraform{logLevel: logLevel}
}

func (t *Terraform) Init(ctx context.Context, wd string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Second)
	defer cancel()

	t.Workdir = wd
	log.Info().Msgf("running init on workdir: %s", t.Workdir)
	cmd := exec.CommandContext(ctx, "terraform", "init", "-upgrade")
	cmd.Dir = t.Workdir
	if t.logLevel == cli.LogLevelDebug {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (t Terraform) ApplyVarMap(ctx context.Context, config map[string]string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Second)
	defer cancel()

	args := []string{}
	args = append(args, "apply")
	args = append(args, "-auto-approve")
	for k, v := range config {
		args = append(args, fmt.Sprintf("-var=%s=%s", k, v))
	}
	log.Info().Msgf("running apply on workdir: %s with args: %s", t.Workdir, args)
	cmd := exec.CommandContext(ctx, "terraform", args...)
	cmd.Dir = t.Workdir

	if t.logLevel == cli.LogLevelDebug {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (t Terraform) ApplyVarFile(ctx context.Context, file string, timeout time.Duration, env map[string]string, target string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	args := []string{}
	args = append(args, "apply")
	args = append(args, "-auto-approve")
	args = append(args, "-var-file="+file)
	if target != "*" {
		args = append(args, "-target="+target)
	}
	log.Info().Msgf("running apply on workdir: %s with args: %s", t.Workdir, args)
	cmd := exec.CommandContext(ctx, "terraform", args...)
	cmd.Dir = t.Workdir
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	if t.logLevel == cli.LogLevelDebug {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (t Terraform) Destroy(ctx context.Context, file string, env map[string]string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()

	args := []string{}
	args = append(args, "destroy")
	args = append(args, "-auto-approve")
	args = append(args, "-var-file="+file)
	log.Info().Msgf("running destroy on workdir: %s with args: %s", t.Workdir, args)
	cmd := exec.CommandContext(ctx, "terraform", args...)
	cmd.Dir = t.Workdir
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	if t.logLevel == cli.LogLevelDebug {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
