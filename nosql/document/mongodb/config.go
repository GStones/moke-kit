package mongodb

type ClusterConfig struct {
	ConnUrl  string `default:"mongo_db://localhost:27017"`
	Username string
	Password string
}
