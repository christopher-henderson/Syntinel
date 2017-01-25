package service

import "syntinel_executor/DAO"

var DockerService = &dockerService{}

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

type dockerService struct{}
