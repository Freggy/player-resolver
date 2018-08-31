package cassandra

import (
	"github.com/gocql/gocql"
	"log"
)

type CassandraSession struct {
	cSession    *gocql.Session
	clusterConf *gocql.ClusterConfig
}

func New(keyspace string) (*CassandraSession, error) {
	cluster := gocql.NewCluster("cassandra.luxor.cloud")
	cluster.Keyspace = "lol"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()

	if err != nil {
		return nil, err
	}

	return &CassandraSession{
		session,
		cluster,
	}, nil
}

func (session *CassandraSession) isPresent(uuid string) bool {
	uuid = ""
	if err := session.cSession.Query("SELECT from uuid_cache WHERE uuid = ?", uuid).Scan(&uuid); err != nil {
		log.Fatal(err)
		return false
	}
	if uuid == "" {
		return false
	} else {
		return true
	}
}

func (session *CassandraSession) Close() {
	if !session.cSession.Closed() {
		session.cSession.Close()
	}
}
