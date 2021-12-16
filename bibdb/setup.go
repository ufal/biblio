package bibdb

func (s *Store) setup() error {
	if _, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS authors (
		uid						TEXT PRIMARY KEY,
		parent					TEXT,
		firstname				TEXT,
		middlename				TEXT,
		lastname				TEXT,
		primaryentry			TEXT,
		academicdegreesprefix	TEXT,
		academicdegreessuffix	TEXT,
		nationality				TEXT,
		university				TEXT,
		faculty					TEXT,
		department				TEXT,
		cuninumber				TEXT,
		cuniaffcode				TEXT,
		email					TEXT,
		note					TEXT
	)`); err != nil {
		return err
	}
	return nil
}
