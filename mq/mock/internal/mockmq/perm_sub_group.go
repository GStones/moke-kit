package mockmq

type permSubGroup struct {
	HasSubscriptions
}

func NewPermSubGroup() *permSubGroup {
	return &permSubGroup{}
}
