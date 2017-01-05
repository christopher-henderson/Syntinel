package service

import (
	"io"
	"syntinel_executor/DAO"
)

var DockerService = newDockerService()

type dockerService struct {
}

func newDockerService() *dockerService {
	return &dockerService{}
}

func (d *dockerService) Register(id int, data io.Reader) {
	docker := &DAO.Docker{id}
	docker.Save(data)
}

func (d *dockerService) Update(id int, data io.Reader) {
	d.Register(id, data)
}

func (d *dockerService) Delete(id int) {
	docker := &DAO.Docker{id}
	docker.Delete()
}
