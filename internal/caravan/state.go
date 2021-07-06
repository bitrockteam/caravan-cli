package caravan

import "fmt"

type Status int

const (
	InitMissing Status = iota
	InitDone
	BakingDone
	InfraCleanDone
	InfraCleanRunning
	InfraDeployRunning
	InfraDeployDone
	PlatformCleanDone
	PlatformCleanRunning
	PlatformDeployRunning
	PlatformDeployDone
	ApplicationCleanDone
	ApplicationCleanRunning
	ApplicationDeployRunning
	ApplicationDeployDone
)

func (s Status) String() string {
	return [...]string{
		fmt.Sprintf("%d-InitMissing", s),
		fmt.Sprintf("%d-InitDone", s),
		fmt.Sprintf("%d-BakingDone", s),
		fmt.Sprintf("%d-InfraCleanDone", s),
		fmt.Sprintf("%d-InfraCleanRunning", s),
		fmt.Sprintf("%d-InfraDeployRunning", s),
		fmt.Sprintf("%d-InfraDeployDone", s),
		fmt.Sprintf("%d-PlatformCleanDone", s),
		fmt.Sprintf("%d-PlatformCleanRunning", s),
		fmt.Sprintf("%d-PlatformDeployRunning", s),
		fmt.Sprintf("%d-PlatformDeployDone", s),
		fmt.Sprintf("%d-ApplicationCleanDone", s),
		fmt.Sprintf("%d-ApplicationCleanRunning", s),
		fmt.Sprintf("%d-ApplicationDeployRunning", s),
		fmt.Sprintf("%d-ApplicationDeployDone", s),
	}[s]
}
