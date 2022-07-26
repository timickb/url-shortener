package store

type StubStore struct{}

func (t StubStore) Open() error {
	return nil
}

func (t StubStore) Close() error {
	return nil
}

func (t StubStore) CreateLink(s string) (string, error) {
	return "t3st_h4sh_", nil
}

func (t StubStore) RestoreLink(s string) (string, error) {
	return "https://example.com", nil
}
