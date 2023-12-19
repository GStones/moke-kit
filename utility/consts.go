package utility

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
	return d == DeploymentsProd
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
