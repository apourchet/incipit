package etcd

type mockEtcdClient struct {
	store map[string]string
}

func NewMockEtcdClient() EtcdClient {
	m := &mockEtcdClient{}
	m.store = make(map[string]string)
	return m
}

func (m *mockEtcdClient) Get(key string) (string, bool, error) {
	val, ok := m.store[key]
	if !ok {
		return "", false, nil
	}
	return val, true, nil
}

func (m *mockEtcdClient) Set(key string, value string) error {
	m.store[key] = value
	return nil
}

func (m *mockEtcdClient) Delete(key string) error {
	delete(m.store, key)
	return nil
}
