package sdk

type MockAlpha struct {
}

func (d *MockAlpha) GetPlayerCapacity() (int64, error) {
	return 0, nil
}

func (d *MockAlpha) SetPlayerCapacity(capacity int64) error {
	return nil
}

func (d *MockAlpha) PlayerConnect(id string) (bool, error) {
	return true, nil
}

func (d *MockAlpha) PlayerDisconnect(id string) (bool, error) {
	return true, nil
}

func (d *MockAlpha) GetPlayerCount() (int64, error) {
	return 0, nil
}

func (d *MockAlpha) IsPlayerConnected(id string) (bool, error) {
	return true, nil
}

func (d *MockAlpha) GetConnectedPlayers() ([]string, error) {
	return make([]string, 0), nil
}

func (d *MockAlpha) GetCounterCount(key string) (int64, error) {
	return 0, nil
}

func (d *MockAlpha) IncrementCounter(key string, amount int64) (bool, error) {
	return true, nil
}

func (d *MockAlpha) DecrementCounter(key string, amount int64) (bool, error) {
	return true, nil
}

func (d *MockAlpha) SetCounterCount(key string, amount int64) (bool, error) {
	return true, nil
}

func (d *MockAlpha) GetCounterCapacity(key string) (int64, error) {
	return 0, nil
}

func (d *MockAlpha) SetCounterCapacity(key string, amount int64) (bool, error) {
	return true, nil
}

func (d *MockAlpha) GetListCapacity(key string) (int64, error) {
	return 0, nil
}

func (d *MockAlpha) SetListCapacity(key string, amount int64) (bool, error) {
	return true, nil
}

func (d *MockAlpha) ListContains(key, value string) (bool, error) {
	return true, nil
}

func (d *MockAlpha) GetListLength(key string) (int, error) {
	return 0, nil
}

func (d *MockAlpha) GetListValues(key string) ([]string, error) {
	return nil, nil
}

func (d *MockAlpha) AppendListValue(key, value string) (bool, error) {
	return true, nil
}

func (d *MockAlpha) DeleteListValue(key, value string) (bool, error) {
	return true, nil
}
