package utility

import "fmt"

type Deployments string

const (
	DeploymentsLocal Deployments = "local"
	DeploymentsDev   Deployments = "dev"
	DeploymentsProd  Deployments = "prod"
)

func (d Deployments) String() string {
	return string(d)
}

func ParseDeployments(value string) Deployments {
	switch Deployments(value) {
	case DeploymentsLocal:
		return DeploymentsLocal
	case DeploymentsDev:
		return DeploymentsDev
	case DeploymentsProd:
		return DeploymentsProd
	default:
		panic(fmt.Errorf(`"%s" is an unknown deployments`, value))
	}
}

const (
	TokenContextKey = "bearer"
	UIDContextKey   = "uid"
)
