package caravan

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
)

func (s Status) String() string {
	return [...]string{"InitMissing", "InitDone", "BakingDone", "InfraCleanDone", "InfraCleanRunning", "InfraDeployRunning", "InfraDeployDone", "PlatformCleanDone", "PlatformCleanRunning", "PlatformDeployRunning", "PlatformDeployDone"}[s]
}
