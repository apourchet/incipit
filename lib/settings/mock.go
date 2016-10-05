package settings

import "fmt"

type mock struct {
	store map[string]map[string]string
}

func NewMockSettings() SettingsClient {
	m := &mock{}
	m.store = make(map[string]map[string]string)
	return m
}

func (m *mock) Set(userid string, field string, value string) error {
	settings, found := m.store[userid]
	if !found {
		m.store[userid] = make(map[string]string)
		settings = m.store[userid]
	}
	settings[field] = value
	return nil
}

func (m *mock) Get(userid string, field string) (value string, err error) {
	settings, found := m.store[userid]
	if !found {
		return "", fmt.Errorf("Setting not found")
	}

	value, found = settings[field]
	if !found {
		return "", fmt.Errorf("Setting not found")
	}
	return value, nil
}
