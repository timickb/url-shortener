package store

type MockStore struct{}

func (t MockStore) Open() error {
	return nil
}

func (t MockStore) Close() error {
	return nil
}

func (t MockStore) CreateLink(s string) (string, error) {
	return "t3st_h4sh_", nil
}

func (t MockStore) RestoreLink(s string) (string, error) {
	return "https://example.com", nil
}
