package part

func (s *service) InitWithDummy() error {
	return s.repository.InitWithDummy()
}
