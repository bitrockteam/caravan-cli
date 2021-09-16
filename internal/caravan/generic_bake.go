package caravan

import (
	"caravan/internal/terraform"
	"fmt"
	"os"
	"time"
)

type GenericBake struct {
	GenericProvider
}

func (g GenericBake) Bake() error {
	if _, err := os.Stat(g.Caravan.WorkdirProject); os.IsNotExist(err) {
		return fmt.Errorf("please run init before bake")
	}

	t := terraform.Terraform{}
	if err := t.Init(g.Caravan.WorkdirBaking); err != nil {
		return err
	}
	env := map[string]string{}
	if err := t.ApplyVarFile(g.Caravan.WorkdirBakingVars, 1200*time.Second, env, "*"); err != nil {
		return err
	}
	return nil
}
