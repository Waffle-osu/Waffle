package migrations

import (
	"database/sql"
	"strings"
)

type MigrationTemplateStruct struct{}

func (migration MigrationTemplateStruct) Apply(db *sql.DB) error {
	return nil
}

func (migration MigrationTemplateStruct) Remove(db *sql.DB) error {
	return nil
}

func MigrationHelperRunSplitSql(sql string, db *sql.DB) error {
	for _, part := range strings.Split(sql, "@@@@") {
		_, err := db.Exec(part)

		if err != nil {
			return err
		}
	}

	return nil
}
