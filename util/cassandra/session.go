package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/luxordynamics/player-resolver/util/mojang"
)


type Session struct {
	session     *gocql.Session
	clusterConf *gocql.ClusterConfig
}

type Entry struct {
	Mapping    *mojang.PlayerNameMapping
	LastUpdate int64
	LastQuery  int64
}

var (
	createKeySpace         = "CREATE KEYSPACE luxor WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor' : 1};"
	createTableQuery       = "CREATE TABLE uuid_cache (uuid UUID PRIMARY KEY, name text, last_change timestamp, last_query timestamp)"
	updateNameQuery        = "UPDATE uuid_cache SET name = ? WHERE uuid = ?"
	updateLastQuery        = "UPDATE uuid_cache SET last_query = ? WHERE uuid = ?"
	uuidExistsQuery        = "SELECT count(*) FROM uuid_cache WHERE uuid = ?"
	nameExistsQuery        = "SELECT count(*) FROM uuid_cache WHERE name = ?"
	selectByUuidEntryQuery = "SELECT * FROM uuid_cache WHERE uuid = ?"
	selectByNameEntryQuery = "SELECT * FROM uuid_cache WHERE name = ?"
	insertQuery            = "INSERT INTO uuid_cache (uuid, name, last_change, last_query) VALUES (?, ?, ?, ?)"
)

// New creates a new instance of Session and directly connects to the cluster.
func New(host string) (*Session, error) {
	cluster := gocql.NewCluster(host)
	cluster.Keyspace = "luxor"
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

// Setup creates the keyspace and the table.
func (session *Session) Setup() error {
	if err := session.session.Query(createKeySpace).Exec(); err != nil {
		return err
	}

	if err := session.session.Query(createTableQuery).Exec(); err != nil {
		return err
	}

	return nil
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
func (session *Session) UuidEntryExists(uuid string) (bool, error) {
	return session.entryExists(uuid, uuidExistsQuery)
}

// NameEntryExists returns whether or not an entry with the given name exists.
func (session *Session) NameEntryExists(name string) (bool, error) {
	return session.entryExists(name, nameExistsQuery)
}

// WriteEntry inserts a name with the associated uuid and the last update time.
func (session *Session) WriteEntry(uuid string, name string, lastUpdated int64, lastQuery int64) error {
	if err := session.session.Query(insertQuery, uuid, name, lastUpdated, lastQuery).Exec(); err != nil {
		return err
	}
	return nil
}

// UpdateName updates the name column.
func (session *Session) UpdateName(name string, uuid string) error {
	if err := session.session.Query(updateNameQuery, name, uuid).Exec(); err != nil {
		return err
	}
	return nil
}

// UpdateLastQuery updates the last_query column.
func (session *Session) UpdateLastQuery(lastQuery int64, uuid string) error {
	if err := session.session.Query(updateLastQuery, lastQuery, uuid).Exec(); err != nil {
		return err
	}
	return nil
}

// Retrieves an Entry from the database by using the given key and query string. This is a synchronous function call.
func (session *Session) entryFromDatabase(key string, query string) (*Entry, error) {
	var uuid string
	var name string
	var lastUpdate int64
	var lastQuery int64

	if err := session.session.Query(query, key).Scan(&uuid, &name, &lastUpdate, &lastQuery); err != nil {
		return nil, err
	}

	return &Entry{
		Mapping: &mojang.PlayerNameMapping{
			Uuid: uuid,
			Name: name,
		},
		LastUpdate: lastUpdate,
		LastQuery: lastQuery,
	}, nil
}

// Checks if an Entry with the given key exists in the database.
func (session *Session) entryExists(key string, query string) (bool, error) {
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
