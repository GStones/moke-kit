package sdk

type MockCounterList struct {
}

func (d *MockCounterList) GetCounterCount(key string) (int64, error) {
	return 0, nil
}

func (d *MockCounterList) IncrementCounter(key string, amount int64) error {
	return nil
}

func (d *MockCounterList) DecrementCounter(key string, amount int64) error {
	return nil
}

func (d *MockCounterList) SetCounterCount(key string, amount int64) error {
	return nil
}

func (d *MockCounterList) GetCounterCapacity(key string) (int64, error) {
	return 0, nil
}

func (d *MockCounterList) SetCounterCapacity(key string, amount int64) error {
	return nil
}

func (d *MockCounterList) GetListCapacity(key string) (int64, error) {
	return 0, nil
}

func (d *MockCounterList) SetListCapacity(key string, amount int64) error {
	return nil
}

func (d *MockCounterList) ListContains(key, value string) (bool, error) {
	return true, nil
}

func (d *MockCounterList) GetListLength(key string) (int, error) {
	return 0, nil
}

func (d *MockCounterList) GetListValues(key string) ([]string, error) {
	return nil, nil
}

func (d *MockCounterList) AppendListValue(key, value string) error {
	return nil
}

func (d *MockCounterList) DeleteListValue(key, value string) error {
	return nil
}
