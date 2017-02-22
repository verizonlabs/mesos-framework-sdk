package cassandra

import (
	"github.com/gocql/gocql"
	"mesos-framework-sdk/persistence"
	"strings"
)

type Cassandra struct {
	session *gocql.Session
}

// Creates a new Cassandra client
func NewClient(endpoints []string, keyspace string) persistence.DBStorage {
	cluster := gocql.NewCluster(endpoints...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		panic("Failed to create Cassandra client")
	}

	return &Cassandra{
		session: session,
	}
}

// Dynamically formats columns for use in a query.
func (c *Cassandra) formatColumns(cols []string) string {
	var columns string
	for column := range cols {
		columns += column + ", "
	}

	return strings.TrimRight(columns, ", ")
}

// Generates a bound parameter string based on the number of values that need binding.
func (c *Cassandra) formatBindings(vals []string) string {
	var bindings string
	for _ := range vals {
		bindings += "?, "
	}

	return strings.TrimRight(bindings, ", ")
}

// Generates a formatted WHERE clause with values to be bound for use in a query.
func (c *Cassandra) formatWhereClause(query string, where map[string]string) (string, []string) {
	var vals []string
	query += " WHERE "
	for col, val := range where {
		query += col + " = ? AND "
		vals = append(vals, val)
	}
	return strings.TrimRight(query, " AND "), vals
}

// Inserts data into the database.
func (c *Cassandra) Create(table string, cols, vals []string) error {
	return c.session.Query("INSERT INTO "+table+" ("+c.formatColumns(cols)+") VALUES ("+c.formatBindings(vals)+")", vals...).Exec()
}

// Selects data from the database using an optional WHERE clause.
func (c *Cassandra) Read(table string, cols []string, where map[string]string) ([]map[string]interface{}, error) {
	query := "SELECT " + c.formatColumns(cols) + " FROM " + table + ""
	if where != nil {
		query, vals := c.formatWhereClause(query, where)

		return c.session.Query(query, vals...).Iter().SliceMap(), nil
	}

	return c.session.Query(query).Iter().SliceMap(), nil
}

// Updates data in the database using an optional WHERE clause.
func (c *Cassandra) Update(table string, data, where map[string]string) error {
	query := "UPDATE " + table + " SET "
	for col, val := range data {
		query += col + " = " + val + ", "
	}
	query = strings.TrimRight(query, ", ")
	if where != nil {
		query, vals := c.formatWhereClause(query, where)

		return c.session.Query(query, vals...).Exec()
	}

	return c.session.Query(query).Exec()
}

// Deletes data from the database.
func (c *Cassandra) Delete(table string, cols []string, where map[string]string) error {
	query := "DELETE " + c.formatColumns(cols) + " FROM " + table + ""
	if where != nil {
		query, vals := c.formatWhereClause(query, where)

		return c.session.Query(query, vals...).Exec()
	}

	return c.session.Query(query).Exec()
}
