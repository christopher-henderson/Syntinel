package database

const (
	SelectDockerfile = "SELECT ID, Content, Hash FROM Dockerfile WHERE ID=?"
	SelectScript     = "SELECT ID, Content, Hash FROM Script WHERE ID=?"
	SelectTest       = "SELECT * FROM Test WHERE ID=?"
)
