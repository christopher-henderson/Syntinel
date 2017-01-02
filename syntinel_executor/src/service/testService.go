package service

import (
	"syntinel_executor/DAO"
	"syntinel_executor/PAO"
)

var ts = &testService{}

var GetTestService func() *testService = func() func() *testService {
	PAO.StartPAO()
	return func() *testService {
		return ts
	}
}()

type testService struct {
}

func (t *testService) Run(id int) error {
	PAO.Run(id)
	return nil
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
