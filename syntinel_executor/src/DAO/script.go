package DAO

import "syntinel_executor/DAO/database"

var Script = &script{}

func (s *script) Get(id int) (*database.ScriptEntity, error) {
	return database.GetScript(id)
}

func (s *script) Save(id int, content []byte) error {
	return database.InsertScript(id, string(content))
}

func (s *script) Delete(id int) error {
	return database.DeleteScript(id)
}

func (s *script) Update(id int, content []byte) error {
	return database.UpdateScript(id, string(content))
}

type script struct{}
