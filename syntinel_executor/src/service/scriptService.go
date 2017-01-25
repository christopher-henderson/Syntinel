package service

import "syntinel_executor/DAO"

var ScriptService = &scriptService{}

func (s *scriptService) Register(id int, body []byte) error {
	return DAO.Script.Save(id, body)
}

func (s *scriptService) Update(id int, body []byte) error {
	return DAO.Script.Update(id, body)
}

func (t *scriptService) Delete(id int) error {
	return DAO.Script.Delete(id)
}

type scriptService struct{}
