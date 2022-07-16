package store

type TestStore struct{}

func (t TestStore) Open() error {
	return nil
}

func (t TestStore) Close() error {
	return nil
}

func (t TestStore) CreateLink(s string) (string, error) {
	return "t3st_h4sh_", nil
}

func (t TestStore) RestoreLink(s string) (string, error) {
	return "https://example.com", nil
}
