package migrations

import "database/sql"

type CreateSiteTokensTablesStruct struct{}

func (migration CreateSiteTokensTablesStruct) Apply(db *sql.DB) {
	creationSql :=
		`
		CREATE TABLE site_tokens (
			token_id      bigint       unsigned NOT NULL AUTO_INCREMENT,
			token_hash    varchar(128)          NOT NULL,
			creation_date datetime              NOT NULL DEFAULT CURRENT_TIMESTAMP,
			
			PRIMARY KEY (token_id, token_hash)
		) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
	`

	db.Query(creationSql)
}

func (migration CreateSiteTokensTablesStruct) Remove(db *sql.DB) {
	db.Query("DROP TABLE waffle.site_tokens")
}
