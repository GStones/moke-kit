package siface

type IAuth interface {
	Auth(token string) (string, error)
}
