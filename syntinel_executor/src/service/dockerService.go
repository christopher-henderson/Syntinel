package service

import "syntinel_executor/DAO"

var DockerService = &dockerService{}

func (d *dockerService) Register(id int, body []byte) error {
	return DAO.Docker.Save(id, body)
}

func (d *dockerService) Update(id int, body []byte) error {
	return DAO.Docker.Update(id, body)
}

func (d *dockerService) Delete(id int) error {
	return DAO.Docker.Delete(id)
}

type dockerService struct{}
