package migrations

import (
	"database/sql"
)

type CreateWaffleBotStruct struct{}

func (migration CreateWaffleBotStruct) Apply(db *sql.DB) error {
	creationSql :=
		`
		INSERT INTO waffle.users (username, password) VALUES ("WaffleBot", "no!");
		@@@@
		INSERT INTO waffle.stats (user_id, mode) VALUES (1, 0), (1, 1), (1, 2), (1, 3);
	`

	return MigrationHelperRunSplitSql(creationSql, db)
}

func (migration CreateWaffleBotStruct) Remove(db *sql.DB) error {
	deletionSql :=
		`
		DELETE FROM waffle.users WHERE user_id = 1;
	`

	return MigrationHelperRunSplitSql(deletionSql, db)
}
