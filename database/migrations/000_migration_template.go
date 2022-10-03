package migrations

import "database/sql"

type MigrationTemplateStruct struct{}

func (migration MigrationTemplateStruct) Apply(db *sql.DB) {

}

func (migration MigrationTemplateStruct) Remove(db *sql.DB) {
	
}
