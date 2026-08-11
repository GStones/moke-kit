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

// matchDeployment reports whether d equals name or uses name as a prefix
// followed by '_' or '-' (e.g. prod, prod_gcp, prod-aws).
func matchDeployment(d Deployments, name string) bool {
	s := d.String()
	if s == name {
		return true
	}
	return strings.HasPrefix(s, name+"_") || strings.HasPrefix(s, name+"-")
}

// IsProd returns true for "prod" and variants like "prod_gcp", "prod-aws".
func (d Deployments) IsProd() bool {
	return matchDeployment(d, DeploymentsProd.String())
}

// IsDev returns true for "dev" and variants like "dev_test", "dev-load".
func (d Deployments) IsDev() bool {
	return matchDeployment(d, DeploymentsDev.String())
}

// IsLocal returns true for "local" and variants like "local_67", "local-name".
func (d Deployments) IsLocal() bool {
	return matchDeployment(d, DeploymentsLocal.String())
}

func ParseDeployments(value string) Deployments {
	value = strings.ToLower(strings.TrimSpace(value))
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
