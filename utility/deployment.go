package utility

import "strings"

type Deployments string

const (
	DeploymentsLocal Deployments = "local"
	DeploymentsDev   Deployments = "dev"
	DeploymentsProd  Deployments = "prod"
)

func (d Deployments) String() string {
	return string(d)
}

func (d Deployments) IsProd() bool {
	return strings.Contains(d.String(), DeploymentsProd.String())
}

func (d Deployments) IsDev() bool {
	return strings.Contains(d.String(), DeploymentsDev.String())
}

func (d Deployments) IsLocal() bool {
	return strings.Contains(d.String(), DeploymentsLocal.String())
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
		return Deployments(value)
	}
}
