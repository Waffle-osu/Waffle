package migrations

import "database/sql"

type CreateDatabaseVersionStruct struct{}

func (migration CreateDatabaseVersionStruct) Apply(db *sql.DB) {
	creationQuery :=
		`
	CREATE TABLE waffle.database_state (
		version       INT      NOT NULL,
		last_modified DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP

		PRIMARY KEY (version)
	);

	INSERT INTO waffle.database_state (version) VALUES (1)
`

	db.Query(creationQuery)
}

func (migration CreateDatabaseVersionStruct) Remove(db *sql.DB) {
	db.Query("DROP TABLE waffle.database_state")
}
