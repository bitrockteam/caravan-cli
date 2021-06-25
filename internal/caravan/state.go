package caravan

type Status int

const (
	InitMissing Status = iota
	InitDone
	BakingDone
	InfraDeployRunning
	InfraDeployDone
	PlatformDeployRunning
	PlatformDeployDone
)

func (s Status) String() string {
	return [...]string{"InitMissing", "InitDone", "BakingDone", "InfraDeployRunning", "InfraDeployDone", "PlatformDeployRunning", "PlatformDeployDone"}[s]
}
