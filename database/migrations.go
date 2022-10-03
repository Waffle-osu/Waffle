package database

import (
	"Waffle/database/migrations"
	"database/sql"
)

type DatabaseMigration interface {
	Apply(db *sql.DB)
	Remove(db *sql.DB)
}

var Migrations map[uint]DatabaseMigration

func InitializeMigrations() {
	Migrations = make(map[uint]DatabaseMigration)

	Migrations[001] = migrations.CreateDatabaseVersionStruct{}
	Migrations[002] = migrations.CreateUserTablesStruct{}
	Migrations[003] = migrations.CreateBeatmapTablesStruct{}
	Migrations[004] = migrations.CreateScoreTablesStruct{}
	Migrations[005] = migrations.IrcAndUpdaterTablesStruct{}
	Migrations[006] = migrations.AchievementTablesStruct{}
	Migrations[007] = migrations.CreateSiteTokensTablesStruct{}
}
