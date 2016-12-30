package service

var TestService = testService{}

type testService struct {
}

func (t *testService) Kill(id int) error {
	return nil
}

func (t *testService) Run(id int) error {
	return nil
}

func (t *testService) Query(id int) (interface{}, error) {
	return id, nil
}

func (t *testService) Delete(id int) error {
	return nil
}
