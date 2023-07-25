package couchbase

type ClusterConfig struct {
	ConnUrl  string `default:"couchbase://couchbase"`
	Username string
	Password string
}
