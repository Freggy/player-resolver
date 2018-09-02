package cassandra

import (
	"github.com/gocql/gocql"
	"gitlab.com/luxordynamics/player-resolver/internal/mojang"
)

type CassandraSession struct {
	cSession    *gocql.Session
	clusterConf *gocql.ClusterConfig
}

type Entry struct {
	Mapping    mojang.PlayerNameMapping
	LastUpdate int64
}

var (
	uuidExistsQuery        = "SELECT count(*) FROM uuid_cache WHERE uuid = ?"
	nameExistsQuery        = "SELECT count(*) FROM uuid_cache WHERE name = ?"
	selectByUuidEntryQuery = "SELECT * FROM uuid_cache WHERE uuid = ?"
	selectByNameEntryQuery = "SELECT * FROM uuid_cache WHERE name = ?"
	insertQuery            = "INSERT INTO uuid_cache (uuid, name, last_update) VALUES (?, ?, ?)"
)

// Creates a new instance of CassandraSession and directly connects to the cluster.
func New() (*CassandraSession, error) {
	cluster := gocql.NewCluster("localhost")
	cluster.Keyspace = "luxor_cloud"
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

// Gets an Entry by its UUID from the database.
func (session *CassandraSession) EntryByUuid(uuid string) (*Entry, error) {
	return session.entryFromDatabase(uuid, selectByUuidEntryQuery)
}

// Gets an Entry by its name from the database.
func (session *CassandraSession) EntryByName(name string) (*Entry, error) {
	return session.entryFromDatabase(name, selectByNameEntryQuery)
}

// Checks whether or not an entry with the given uuid exists
func (session *CassandraSession) UuidEntryExists(uuid string) (ret bool, err error) {
	return session.entryExists(uuid, uuidExistsQuery)
}

// Checks whether or not an entry with the given name exists.
func (session *CassandraSession) NameEntryExists(name string) (ret bool, err error) {
	return session.entryExists(name, nameExistsQuery)
}

func (session *CassandraSession) WriteEntry(uuid string, name string, lastUpdated int64) error {
	if err := session.cSession.Query("INSERT INTO uuid_cache (uuid, name, last_update) VALUES (?, ?, ?)", uuid, name, lastUpdated).Exec(); err != nil {
		return err
	}
	return nil
}

// Retrieves an Entry from the database by using the given key and query string. This is a synchronous function call.
func (session *CassandraSession) entryFromDatabase(key string, query string) (entry *Entry, err error) {
	var uuid string
	var name string
	var lastUpdate int64

	if err := session.cSession.Query(query, key).Scan(&uuid, &name, &lastUpdate); err != nil {
		return nil, err
	}

	return &Entry{
		Mapping: mojang.PlayerNameMapping{
			Uuid: uuid,
			Name: name,
		},
		LastUpdate: lastUpdate,
	}, nil
}

// Checks if an Entry with the given key exists in the database.
func (session *CassandraSession) entryExists(key string, query string) (ret bool, err error) {
	count := 0
	if err := session.cSession.Query(query, key).Scan(&count); err != nil {
		return false, err
	} else if count == 1 {
		return true, nil
	}
	return false, nil
}

// Closes the underlying Cassandra session.
func (session *CassandraSession) Close() {
	if !session.cSession.Closed() {
		session.cSession.Close()
	}
}
