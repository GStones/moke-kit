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

// IsProd returns true if the deployment contains "prod", e.g. "prod", "prod_gcp", "prod_aws"... etc
func (d Deployments) IsProd() bool {
	return strings.Contains(d.String(), DeploymentsProd.String())
}

// IsDev returns true if the deployment contains "dev", e.g. "dev", "dev_test", "dev_load"... etc
func (d Deployments) IsDev() bool {
	return strings.Contains(d.String(), DeploymentsDev.String())
}

// IsLocal returns true if the deployment contains "local", e.g. "local", "local_67", "local_name"... etc
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
