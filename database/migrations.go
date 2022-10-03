package database

import (
	"Waffle/database/migrations"
	"Waffle/helpers"
	"database/sql"
)

type DatabaseMigration interface {
	Apply(db *sql.DB) error
	Remove(db *sql.DB) error
}

var Migrations map[int]DatabaseMigration
var DatabaseVersion int = -1

func InitializeMigrations() {
	Migrations = make(map[int]DatabaseMigration)

	Migrations[001] = migrations.CreateDatabaseVersionStruct{}
	Migrations[002] = migrations.CreateUserTablesStruct{}
	Migrations[003] = migrations.CreateBeatmapTablesStruct{}
	Migrations[004] = migrations.CreateScoreTablesStruct{}
	Migrations[005] = migrations.IrcAndUpdaterTablesStruct{}
	Migrations[006] = migrations.AchievementTablesStruct{}
	Migrations[007] = migrations.CreateSiteTokensTablesStruct{}
}

func InitializeDatabaseVersion() {
	//Check for database version existing, if not then we have to run *everything*
	databaseStateResult, err := Database.Query("SHOW TABLES LIKE 'database_state'")

	if err != nil {
		panic("Failed to query for database_state")
	}

	if databaseStateResult.Next() {
		versionResult, versionErr := Database.Query("SELECT version FROM database_state LIMIT 1")

		if versionErr != nil {
			panic("Failed to query for database version")
		}

		if versionResult.Next() {
			versionResult.Scan(&DatabaseVersion)
		}
	} else {
		DatabaseVersion = 0
	}
}

func RunNecessaryMigrations() {
	for i := DatabaseVersion + 1; i != len(Migrations)+1; i++ {
		migration, exists := Migrations[i]

		if !exists {
			continue
		}

		err := migration.Apply(Database)

		if err != nil {
			helpers.Logger.Panicf("[Database] Migration %03d Failed!\n Error: %s\n", i, err)
		} else {
			helpers.Logger.Printf("[Database] Migration %03d successfully applied.\n", i)
			Database.Query("UPDATE database_state SET version = ? WHERE id = 1", i)
		}
	}
}
