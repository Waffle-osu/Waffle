package migrations

import (
	"database/sql"
)

type MigrationInsertAchievements struct{}

func (migration MigrationInsertAchievements) Apply(db *sql.DB) error {
	creationSql :=
		`
		INSERT INTO waffle.osu_achievements (name, image) VALUES
        	("500 Combo",                         "combo500.png"),
        	("750 Combo",                         "combo750.png"),
        	("1000 Combo",                        "combo1000.png"),
        	("2000 Combo",                        "combo2000.png"),
        	("Video Game Pack vol.1",             "gamer1.png"),
        	("Video Game Pack vol.2",             "gamer2.png"),
        	("Video Game Pack vol.3",             "gamer3.png"),
        	("Video Game Pack vol.4",             "gamer4.png"),
        	("Anime Pack vol.1",                  "anime1.png"),
        	("Anime Pack vol.2",                  "anime2.png"),
        	("Anime Pack vol.3",                  "anime3.png"),
        	("Anime Pack vol.4",                  "anime4.png"),
        	("Internet! Pack vol.1",              ".png"),
        	("Internet! Pack vol.2",              ".png"),
        	("Internet! Pack vol.3",              ".png"),
        	("Internet! Pack vol.4",              ".png"),
        	("Rhythm Game Pack vol.1",            "rhythm1.png"),
        	("Rhythm Game Pack vol.2",            "rhythm2.png"),
        	("Rhythm Game Pack vol.3",            "rhythm3.png"),
        	("Rhythm Game Pack vol.4",            "rhythm4.png"),
        	("Catch 20000 Fruits",                "fruitsalad.png"),
        	("Catch 200000 Fruits",               "fruitplatter.png"),
        	("Catch 2000000 Fruits",              "fruitod.png"),
        	("5000 Plays",                        "plays1.png"),
        	("15000 Plays",                       "plays2.png"),
        	("25000 Plays",                       "plays3.png"),
        	("50000 Plays",                       "plays4.png"),
        	("30000 Drum Hits",                   "taiko1.png"),
        	("300000 Drum Hits",                  "taiko2.png"),
        	("3000000 Drum Hits",                 "taiko3.png"),
        	("Don't let the bunny distract you!", "bunny.png"),
        	("S-Ranker",                          "s-ranker.png"),
        	("Most Improved",                     "improved.png"),
        	("Non-stop Dancer",                   "dancer.png");
	`

	return MigrationHelperRunSplitSql(creationSql, db)
}

func (migration MigrationInsertAchievements) Remove(db *sql.DB) error {
	return nil
}
