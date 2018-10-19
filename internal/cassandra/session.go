package cassandra

import (
	"github.com/gocql/gocql"
	"gitlab.com/luxordynamics/player-resolver/internal/mojang"
)

type Session struct {
	session     *gocql.Session
	clusterConf *gocql.ClusterConfig
}

type Entry struct {
	Mapping    *mojang.PlayerNameMapping
	LastUpdate int64
}

var (
	uuidExistsQuery        = "SELECT count(*) FROM uuid_cache WHERE uuid = ?"
	nameExistsQuery        = "SELECT count(*) FROM uuid_cache WHERE name = ?"
	selectByUuidEntryQuery = "SELECT * FROM uuid_cache WHERE uuid = ?"
	selectByNameEntryQuery = "SELECT * FROM uuid_cache WHERE name = ?"
	insertQuery            = "INSERT INTO uuid_cache (uuid, name, last_update) VALUES (?, ?, ?)"
)

// New creates a new instance of Session and directly connects to the cluster.
func New() (*Session, error) {
	cluster := gocql.NewCluster("localhost")
	cluster.Keyspace = "luxor_cloud"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()

	if err != nil {
		return nil, err
	}

	return &Session{
		session,
		cluster,
	}, nil
}

// EntryByUuid returns an Entry by its UUID from the database.
func (session *Session) EntryByUuid(uuid string) (*Entry, error) {
	return session.entryFromDatabase(uuid, selectByUuidEntryQuery)
}

// EntryByName returns an Entry by its name from the database.
func (session *Session) EntryByName(name string) (*Entry, error) {
	return session.entryFromDatabase(name, selectByNameEntryQuery)
}

// UuidEntryExists returns whether or not an entry with the given uuid exists
func (session *Session) UuidEntryExists(uuid string) (ret bool, err error) {
	return session.entryExists(uuid, uuidExistsQuery)
}

// NameEntryExists returns whether or not an entry with the given name exists.
func (session *Session) NameEntryExists(name string) (ret bool, err error) {
	return session.entryExists(name, nameExistsQuery)
}

// WriteEntry inserts a name with the associated uuid and the last update time.
func (session *Session) WriteEntry(uuid string, name string, lastUpdated int64) error {
	err := session.session.Query("INSERT INTO uuid_cache (uuid, name, last_update) VALUES (?, ?, ?)", uuid, name, lastUpdated).Exec();
	if err != nil {
		return err
	}
	return nil
}

// Retrieves an Entry from the database by using the given key and query string. This is a synchronous function call.
func (session *Session) entryFromDatabase(key string, query string) (entry *Entry, err error) {
	var uuid string
	var name string
	var lastUpdate int64

	if err := session.session.Query(query, key).Scan(&uuid, &name, &lastUpdate); err != nil {
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
func (session *Session) entryExists(key string, query string) (ret bool, err error) {
	count := 0
	if err := session.session.Query(query, key).Scan(&count); err != nil {
		return false, err
	} else if count == 1 {
		return true, nil
	}
	return false, nil
}

// Closes the underlying Cassandra session.
func (session *Session) Close() {
	if !session.session.Closed() {
		session.session.Close()
	}
}
