package redis

import "time"

type RedisClient interface {
	Get(key string) (value string, found bool, err error)
	Set(key string, value string) error
	SetExpire(key string, value string, exp time.Duration) error
	Delete(key string) error
}

type mockRedisClient struct {
	store map[string]*mockEntry
}

type mockEntry struct {
	value   string
	expTime *time.Time
}

func NewMockRedisClient() RedisClient {
	m := &mockRedisClient{}
	m.store = make(map[string]*mockEntry)
	return m
}

func (m *mockRedisClient) Get(key string) (string, bool, error) {
	entry, ok := m.store[key]
	if !ok || entry.expTime.Before(time.Now()) {
		return "", false, nil
	}
	return entry.value, true, nil
}

func (m *mockRedisClient) Set(key string, value string) error {
	t := time.Now().Add(10 * 365 * 24 * time.Hour)
	m.store[key] = &mockEntry{value, &t}
	return nil
}

func (m *mockRedisClient) SetExpire(key string, value string, exp time.Duration) error {
	t := time.Now().Add(exp)
	m.store[key] = &mockEntry{value, &t}
	return nil
}

func (m *mockRedisClient) Delete(key string) error {
	delete(m.store, key)
	return nil
}
