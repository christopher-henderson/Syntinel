package service

import (
	"mime/multipart"
	"syntinel_executor/DAO"
)

var ds = &dockerService{}

var GetDockerService func() *dockerService = func() func() *dockerService {
	return func() *dockerService {
		return ds
	}
}()

type dockerService struct {
}

func (d *dockerService) Register(id int, file multipart.File) {
	docker := &DAO.Docker{id}
	docker.Save(file)
}

func (d *dockerService) Update(id int, file multipart.File) {
	d.Register(id, file)
}

func (d *dockerService) Delete(id int) {
	docker := &DAO.Docker{id}
	docker.Delete()
}
