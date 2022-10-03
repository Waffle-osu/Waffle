package migrations

import (
	"database/sql"
)

type CreateDatabaseVersionStruct struct{}

func (migration CreateDatabaseVersionStruct) Apply(db *sql.DB) error {
	creationQuery :=
		`
	CREATE TABLE waffle.database_state (
		id            INT      NOT NULL AUTO_INCREMENT,
		version       INT      NOT NULL,
		last_modified DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

		PRIMARY KEY (id)
	);
@@@@
    INSERT INTO waffle.database_state (version) VALUES (1);
`

	return MigrationHelperRunSplitSql(creationQuery, db)
}

func (migration CreateDatabaseVersionStruct) Remove(db *sql.DB) error {
	_, err := db.Query("DROP TABLE waffle.database_state")

	return err
}
