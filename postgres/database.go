package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
)

const (
	driverName          = "pgx"
	defaultDatabaseName = "postgres"
)

type Database struct {
	DataSource *DataSource
	DB         *sql.DB
	Metadata   *Metadata
}

type DataSource struct {
	Host         string
	Port         string
	Username     string
	Password     string
	NoPassword   bool
	DatabaseName string
}

// Connect connects to the database.
func (d *Database) Connect() error {
	if d.DataSource.Username == "" {
		return errors.Errorf("user must be set")
	}

	if d.DataSource.Host == "" {
		return errors.Errorf("host must be set")
	}

	if d.DataSource.Port == "" {
		return errors.Errorf("port must be set")
	}

	connStr := fmt.Sprintf("host=%s port=%s", d.DataSource.Host, d.DataSource.Port)
	connConfig, err := pgx.ParseConfig(connStr)
	if err != nil {
		return err
	}
	connConfig.Config.User = d.DataSource.Username
	connConfig.Config.Password = d.DataSource.Password
	connConfig.Config.Database = d.DataSource.DatabaseName

	if d.DataSource.DatabaseName == "" {
		connConfig.Config.Database = defaultDatabaseName
		d.DataSource.DatabaseName = defaultDatabaseName
	}

	connectionString := stdlib.RegisterConnConfig(connConfig)
	db, err := sql.Open(driverName, connectionString)
	if err != nil {
		return err
	}
	d.DB = db
	return nil
}

func (d *Database) Ping() error {
	if err := d.DB.Ping(); err != nil {
		return errors.Wrapf(err, "failed to ping database %q", d.DataSource.DatabaseName)
	}
	return nil
}

// FetchSchemas fetches the schemas from the database.
func (d *Database) FetchSchemas() error {
	txn, err := d.DB.Begin()
	if err != nil {
		return errors.Wrapf(err, "failed to begin transaction")
	}
	defer txn.Rollback()
	schemas, err := getSchemas(txn)
	if err != nil {
		return errors.Wrapf(err, "failed to get schemas from database %q", d.DataSource.DatabaseName)
	}
	d.Metadata = &Metadata{
		Schemas: []Schema{},
	}
	for _, schema := range schemas {
		d.Metadata.Schemas = append(d.Metadata.Schemas, Schema{Name: schema})
	}
	return nil
}

type Metadata struct {
	Schemas []Schema
}

type Schema struct {
	Name   string
	Tables []Table
}

type Table struct {
	Name string
}

// func (d *Database) FetchMetadata() error {
// 	txn, err := d.DB.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer txn.Rollback()

// 	schemaList, err := getSchemas(txn)
// 	if err != nil {
// 		return errors.Wrapf(err, "failed to get schemas from database %q", d.DataSource.DatabaseName)
// 	}
// 	tableMap, err := getTables(txn)
// 	if err != nil {
// 		return errors.Wrapf(err, "failed to get tables from database %q", d.DataSource.DatabaseName)
// 	}

// 	if err := txn.Commit(); err != nil {
// 		return err
// 	}
// }

// var listTableQuery = `
// SELECT tbl.schemaname, tbl.tablename,
// 	pg_table_size(format('%s.%s', quote_ident(tbl.schemaname), quote_ident(tbl.tablename))::regclass),
// 	pg_indexes_size(format('%s.%s', quote_ident(tbl.schemaname), quote_ident(tbl.tablename))::regclass),
// 	GREATEST(pc.reltuples::bigint, 0::BIGINT) AS estimate,
// 	obj_description(format('%s.%s', quote_ident(tbl.schemaname), quote_ident(tbl.tablename))::regclass) AS comment
// FROM pg_catalog.pg_tables tbl
// LEFT JOIN pg_class as pc ON pc.oid = format('%s.%s', quote_ident(tbl.schemaname), quote_ident(tbl.tablename))::regclass` + fmt.Sprintf(`
// WHERE tbl.schemaname NOT IN (%s)
// liii,,,,,,,,,,,,,,,,,kl,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,43e
// ORDER BY tbl.schemaname, tbl.tablename;`, systemSchemas)

// // getTables gets all tables of a database.
// func getTables(txn *sql.Tx) (map[string][]*storepb.TableMetadata, error) {
// 	columnMap, err := getTableColumns(txn)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to get table columns")
// 	}
// 	indexMap, err := getIndexes(txn)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "failed to get indices")
// 	}
// 	foreignKeysMap, err := getForeignKeys(txn)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "failed to get foreign keys")
// 	}

// 	tableMap := make(map[string][]*storepb.TableMetadata)
// 	rows, err := txn.Query(listTableQuery)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		table := &storepb.TableMetadata{}
// 		// var tbl tableSchema
// 		var schemaName string
// 		var comment sql.NullString
// 		if err := rows.Scan(&schemaName, &table.Name, &table.DataSize, &table.IndexSize, &table.RowCount, &comment); err != nil {
// 			return nil, err
// 		}
// 		if comment.Valid {
// 			table.Comment = comment.String
// 		}
// 		key := db.TableKey{Schema: schemaName, Table: table.Name}
// 		table.Columns = columnMap[key]
// 		table.Indexes = indexMap[key]
// 		table.ForeignKeys = foreignKeysMap[key]

// 		tableMap[schemaName] = append(tableMap[schemaName], table)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return tableMap, nil
// }

var listSchemaQuery = `
SELECT nspname
FROM pg_catalog.pg_namespace;
`

func getSchemas(txn *sql.Tx) ([]string, error) {
	rows, err := txn.Query(listSchemaQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var schemaName string
		if err := rows.Scan(&schemaName); err != nil {
			return nil, err
		}
		result = append(result, schemaName)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
