package service

import "syntinel_executor/DAO"

var ss = &scriptService{}

var GetScriptService func() *scriptService = func() func() *scriptService {
	return func() *scriptService {
		return ss
	}
}()

type scriptService struct {
}

func (s *scriptService) Register(id int, body []byte) error {
	script := DAO.Script{id}
	script.Save(body)
	return nil
}

func (s *scriptService) Update(id int, body []byte) error {
	return s.Register(id, body)
}

func (t *scriptService) Delete(id int) error {
	script := DAO.Script{id}
	script.Delete()
	return nil
}
