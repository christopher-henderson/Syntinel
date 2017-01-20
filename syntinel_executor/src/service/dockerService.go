package service

import "syntinel_executor/DAO"

var DockerService = newDockerService()

type dockerService struct {
}

func newDockerService() *dockerService {
	return &dockerService{}
}

func (d *dockerService) Register(id int, body []byte) {
	docker := &DAO.Docker{id}
	docker.Save(body)
}

func (d *dockerService) Update(id int, body []byte) {
	d.Register(id, body)
}

func (d *dockerService) Delete(id int) {
	docker := &DAO.Docker{id}
	docker.Delete()
}
