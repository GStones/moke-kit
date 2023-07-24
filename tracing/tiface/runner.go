package tiface

type Runner interface {
	Start() error
	Stop() error
}
