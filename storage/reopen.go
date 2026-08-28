package storage

func OpenAndRead(path, id string) (bool, error) {
	s, e := Open(path)
	if e != nil {
		return false, e
	}
	defer s.Close()
	_, e = s.GetRecord(id)
	return e == nil, e
}
