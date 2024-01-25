package bss

import (
	"Waffle/database"
	"Waffle/helpers"

	"github.com/Waffle-osu/osu-parser/osu_parser"
	"github.com/Waffle-osu/waffle-pp/difficulty"
)

//ill fix the code duplication some time

func RunAndCreateDiffCalc(osu osu_parser.OsuFile, beatmapId int64, beatmapsetId int64) {
	eyupStars := difficulty.CalculateEyupStars(osu)

	for i := 0; i != 4; i++ {
		_, err := database.Database.Exec("INSERT INTO osu_beatmap_difficulty (beatmap_id, beatmapset_id, mode, eyup_stars) VALUES (?, ?, ?, ?)", beatmapId, beatmapsetId, i, eyupStars)

		if err != nil {
			helpers.Logger.Printf("Difficulty Calculation Update failed: %s", err.Error())
		}
	}
}

func RunDiffCalc(osu osu_parser.OsuFile, beatmapId int64) {
	eyupStars := difficulty.CalculateEyupStars(osu)

	_, err := database.Database.Exec("UPDATE osu_beatmap_difficulty SET eyup_stars = ? WHERE beatmap_id = ?", eyupStars, beatmapId)

	if err != nil {
		helpers.Logger.Printf("Difficulty Calculation Update failed: %s", err.Error())
	}
}

func RunDiffCalcMd5(osu osu_parser.OsuFile, md5 string) {
	eyupStars := difficulty.CalculateEyupStars(osu)

	_, err := database.Database.Exec("UPDATE osu_beatmap_difficulty SET eyup_stars = ? WHERE beatmap_id = (SELECT beatmap_id FROM beatmaps WHERE beatmap_md5 = ?)", eyupStars, md5)

	if err != nil {
		helpers.Logger.Printf("Difficulty Calculation Update failed: %s", err.Error())
	}
}
