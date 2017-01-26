package service

import "syntinel_executor/DAO"

var ScriptService = &scriptService{}

func (s *scriptService) Register(id int, body []byte) {
	script := DAO.Script{id}
	script.Save(body)
}

func (s *scriptService) Update(id int, body []byte) {
	s.Register(id, body)
}

func (t *scriptService) Delete(id int) {
	script := DAO.Script{id}
	script.Delete()
}

type scriptService struct{}
