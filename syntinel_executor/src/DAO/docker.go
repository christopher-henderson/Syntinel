package DAO

import "syntinel_executor/DAO/database"

var Docker = &docker{}

func (d *docker) Get(id int) (*database.DockerEntity, error) {
	return database.GetDockerfile(id)
}

func (d *docker) Save(id int, content []byte) error {
	return database.InsertDockerfile(id, string(content))
}

func (d *docker) Delete(id int) error {
	return database.DeleteDockerfile(id)
}

func (d *docker) Update(id int, content []byte) error {
	return database.UpdateDockerfile(id, string(content))
}

// // Path returns the absolute Dockerfile of the script on the filesystem.
// func (d *docker) Path() string {
// 	return fmt.Sprintf("%v%v", utils.DockerFileDirectory(), strconv.Itoa(d.ID))
// }
//
// // tmpPath returns what Path does, but as a hidden file.
// func (d *docker) tmpPath() string {
// 	return fmt.Sprintf("%v.%v", utils.DockerFileDirectory(), strconv.Itoa(d.ID))
// }

type docker struct{}
