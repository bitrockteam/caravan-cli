package cli

type DeployLayer int

const (
	Infrastructure DeployLayer = iota
	Platform
	ApplicationSupport
)
