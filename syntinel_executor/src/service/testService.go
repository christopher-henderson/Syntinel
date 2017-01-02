package service

import "syntinel_executor/DAO"

var TestService = testService{}

type testService struct {
}

func (t *testService) Run(id int) error {
	process := PAO.
}

func (t *testService) Kill(id int) error {
	return nil
}

func (t *testService) Query(id int) (interface{}, error) {
	return id, nil
}

func (t *testService) Delete(id int) error {
	script := DAO.Script{id}
	script.Delete()
	return nil
}

func (t *testService) Register(id int, body []byte) error {
	script := DAO.Script{id}
	script.Save(body)
	return nil
}
