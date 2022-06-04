package common

import (
	"Waffle/database"
	"fmt"
	"strconv"
	"strings"
)

const (
	Achievement500Combo  = 1
	Achievement750Combo  = 2
	Achievement1000Combo = 3
	Achievement2000Combo = 4

	AchievementVideoGamePack1  = 5
	AchievementVideoGamePack2  = 6
	AchievementVideoGamePack3  = 7
	AchievementVideoGamePack4  = 8
	AchievementAnimePack1      = 9
	AchievementAnimePack2      = 10
	AchievementAnimePack3      = 11
	AchievementAnimePack4      = 12
	AchievementInternetPack1   = 13
	AchievementInternetPack2   = 14
	AchievementInternetPack3   = 15
	AchievementInternetPack4   = 16
	AchievementRhythmGamePack1 = 17
	AchievementRhythmGamePack2 = 18
	AchievementRhythmGamePack3 = 19
	AchievementRhythmGamePack4 = 20

	AchievementCatch20000Fruits   = 21
	AchievementCatch200000Fruits  = 22
	AchievementCatch2000000Fruits = 23
	Achievement5000PlaysOsu       = 24
	Achievement15000PlaysOsu      = 25
	Achievement25000PlaysOsu      = 26
	Achievement30000TaikoHits     = 27
	Achievement300000TaikoHits    = 28
	Achievement3000000TaikoHits   = 29

	AchievementMakeUpFC                 = 30
	Achievement5SRanks24Hours           = 31
	AchievementMostImproved             = 32
	AchievementNonstopDancerParaparaMAX = 33
)

func UpdateAchievements(userId uint64, beatmapId int32, beatmapsetId int32, ranking string, playmode int8, maxCombo int32) (queryResult int8, achievements []database.Achievement) {
	achievedAchievements := []database.Achievement{}

	achievedIds := []int32{}

	unachievedAchievementsQuery, unachievedAchievementsQueryErr := database.Database.Query("SELECT achievement_id FROM waffle.osu_achievements WHERE achievement_id NOT IN (SELECT achievement_id FROM waffle.osu_achieved_achievements WHERE user_id = ?)", userId)

	if unachievedAchievementsQueryErr != nil {
		return -2, achievedAchievements
	}

	//Retrieve stats
	statGetResultOsu, osuStats := database.UserStatsFromDatabase(userId, 0)
	statGetResultTaiko, taikoStats := database.UserStatsFromDatabase(userId, 1)
	statGetResultCatch, catchStats := database.UserStatsFromDatabase(userId, 2)

	if statGetResultOsu != 0 || statGetResultTaiko != 0 || statGetResultCatch != 0 {
		return -2, achievedAchievements
	}

	AchieveAchievement := func(achievementId int32) {
		achievedIds = append(achievedIds, achievementId)

		insertAchievedQuery, _ := database.Database.Query("INSERT INTO waffle.osu_achieved_achievements (achievement_id, user_id) VALUES (?, ?)", achievementId, userId)

		if insertAchievedQuery != nil {
			insertAchievedQuery.Close()
		}
	}

	beatmapPackSql := "SELECT COUNT(*) AS 'count' FROM (SELECT * FROM waffle.scores WHERE mapset_best = 1 AND user_id = ? AND beatmapset_id IN ( %s ) GROUP BY beatmapset_id) result"

	CheckPackCompleted := func(setIds string, mapAmount int64) bool {
		packCheckQuery, packCheckQueryErr := database.Database.Query(fmt.Sprintf(beatmapPackSql, setIds), userId)

		if packCheckQueryErr != nil {
			return false
		}

		if packCheckQuery.Next() {
			var count int64

			scanErr := packCheckQuery.Scan(&count)

			packCheckQuery.Close()

			if scanErr != nil {
				return false
			}

			if count == mapAmount {
				return true
			}
		}

		return false
	}

	for unachievedAchievementsQuery.Next() {
		var achievementId int32

		scanErr := unachievedAchievementsQuery.Scan(&achievementId)

		if scanErr != nil {
			return -2, achievedAchievements
		}

		switch achievementId {
		case Achievement500Combo:
			if maxCombo >= 500 {
				AchieveAchievement(Achievement500Combo)
			}
		case Achievement750Combo:
			if maxCombo >= 750 {
				AchieveAchievement(Achievement750Combo)
			}
		case Achievement1000Combo:
			if maxCombo >= 1000 {
				AchieveAchievement(Achievement1000Combo)
			}
		case Achievement2000Combo:
			if maxCombo >= 2000 {
				AchieveAchievement(Achievement2000Combo)
			}
		case AchievementVideoGamePack1:
			setIds := "25, 92, 125, 154, 312, 633, 688, 704, 1092, 1211, 1231, 1281, 1635"

			if CheckPackCompleted(setIds, 13) {
				AchieveAchievement(AchievementVideoGamePack1)
			}
		case AchievementVideoGamePack2:
			setIds := "243, 628, 1044, 1123, 1367, 1525, 1818, 2008, 2128, 2147, 2404, 2420, 2619"

			if CheckPackCompleted(setIds, 13) {
				AchieveAchievement(AchievementVideoGamePack2)
			}
		case AchievementVideoGamePack3:
			setIds := "1890, 2085, 2490, 2983, 3150, 3221, 3384, 3511, 3613, 4033, 4299, 4305, 4629"

			if CheckPackCompleted(setIds, 13) {
				AchieveAchievement(AchievementVideoGamePack3)
			}
		case AchievementVideoGamePack4:
			setIds := "7077, 9580, 9668, 9854, 10104, 10880, 13489, 14205, 14458, 16669, 17373, 21836, 23073"

			if CheckPackCompleted(setIds, 13) {
				AchieveAchievement(AchievementVideoGamePack4)
			}
		case AchievementAnimePack1:
			setIds := "35, 147, 301, 442, 511, 584, 842, 897, 1005, 1377, 1414, 1464, 1806"

			if CheckPackCompleted(setIds, 13) {
				AchieveAchievement(AchievementAnimePack1)
			}
		case AchievementAnimePack2:
			setIds := "86, 150, 162, 205, 212, 302, 496, 521, 956, 2207, 2267, 2329, 2425"

			if CheckPackCompleted(setIds, 13) {
				AchieveAchievement(AchievementAnimePack2)
			}
		case AchievementAnimePack3:
			setIds := "2618, 3030, 4851, 4994, 5010, 5235, 5410, 5480, 5963, 6037, 6257, 6535, 6557"

			if CheckPackCompleted(setIds, 13) {
				AchieveAchievement(AchievementAnimePack3)
			}
		case AchievementAnimePack4:
			setIds := "516, 5438, 6301, 8422, 8829, 9556, 12982, 13036, 13673, 14256, 14694, 16252, 21197"

			if CheckPackCompleted(setIds, 13) {
				AchieveAchievement(AchievementAnimePack4)
			}
		/*

			case AchievementInternetPack1:
				setIds := "13177, 17217, 17724, 18568, 23754, 31419, 31811" // 45341, 76115, 102615, 116487, 148979, 239262, 332436"
			case AchievementInternetPack2:
				setIds := "24152, 25198, 31471, 36920" //53810, 54631, 63500, 74110, 106500, 119980, 130725, 196930, 232505, 299643
			case AchievementInternetPack3:
				setIds := "15849, 28222, 28799, 36225" //47517, 53363, 53569, 70259, 105186, 192763, 221414, 336207, 347433
			case AchievementInternetPack4:
				setIds := "13885, 14672, 21581, 22252, 23058, 24084, 25154, 37563" //45698, 47078, 155910, 176702, 213629, 247243

		*/
		case AchievementRhythmGamePack1:
			setIds := "74, 96, 210, 296, 540, 564, 1078, 1201, 1300, 1317, 1338, 1450, 1452"

			if CheckPackCompleted(setIds, 13) {
				AchieveAchievement(AchievementRhythmGamePack1)
			}
		case AchievementRhythmGamePack2:
			setIds := "1207, 1567, 2534, 3302, 3435, 3499, 4887, 5087, 5177, 5275, 5321, 5349, 5577"

			if CheckPackCompleted(setIds, 13) {
				AchieveAchievement(AchievementRhythmGamePack2)
			}
		case AchievementRhythmGamePack3:
			setIds := "1206, 4357, 4617, 4772, 4954, 5180, 5672, 5696, 6598, 7094, 7237, 7612, 7983"

			if CheckPackCompleted(setIds, 13) {
				AchieveAchievement(AchievementRhythmGamePack3)
			}
		case AchievementRhythmGamePack4:
			setIds := "10842, 11135, 11488, 12052, 12190, 12710, 13249, 14572, 14778, 15241, 18492, 19809, 22401"

			if CheckPackCompleted(setIds, 13) {
				AchieveAchievement(AchievementRhythmGamePack4)
			}
		case AchievementCatch20000Fruits:
			if (catchStats.Hit300 + catchStats.Hit100 + catchStats.Hit50) >= 20000 {
				AchieveAchievement(AchievementCatch20000Fruits)
			}
		case AchievementCatch200000Fruits:
			if (catchStats.Hit300 + catchStats.Hit100 + catchStats.Hit50) >= 200000 {
				AchieveAchievement(AchievementCatch200000Fruits)
			}
		case AchievementCatch2000000Fruits:
			if (catchStats.Hit300 + catchStats.Hit100 + catchStats.Hit50) >= 200000 {
				AchieveAchievement(AchievementCatch2000000Fruits)
			}
		case Achievement5000PlaysOsu:
			if osuStats.Playcount >= 5000 {
				AchieveAchievement(Achievement5000PlaysOsu)
			}
		case Achievement15000PlaysOsu:
			if osuStats.Playcount >= 15000 {
				AchieveAchievement(Achievement15000PlaysOsu)
			}
		case Achievement25000PlaysOsu:
			if osuStats.Playcount >= 25000 {
				AchieveAchievement(Achievement25000PlaysOsu)
			}
		case Achievement30000TaikoHits:
			if (taikoStats.Hit300 + taikoStats.Hit100 + taikoStats.Hit50) >= 30000 {
				AchieveAchievement(Achievement30000TaikoHits)
			}
		case Achievement300000TaikoHits:
			if (taikoStats.Hit300 + taikoStats.Hit100 + taikoStats.Hit50) >= 300000 {
				AchieveAchievement(Achievement300000TaikoHits)
			}
		case Achievement3000000TaikoHits:
			if (taikoStats.Hit300 + taikoStats.Hit100 + taikoStats.Hit50) >= 3000000 {
				AchieveAchievement(Achievement3000000TaikoHits)
			}

		case AchievementMakeUpFC:
			if playmode == 0 && ((beatmapId == 8708 && maxCombo == 447) || (beatmapId == 8707 && maxCombo == 371)) {
				AchieveAchievement(AchievementMakeUpFC)
			}
		case Achievement5SRanks24Hours:
			sRankerQuery, sRankerQueryErr := database.Database.Query("SELECT COUNT(*) AS 'count' FROM waffle.scores WHERE passed = 1 AND ranking IN ('S', 'X', 'SH', 'XH') AND date > DATE_ADD(CURRENT_TIMESTAMP, INTERVAL -1 DAY) AND user_id = ?", userId)

			if sRankerQueryErr != nil {
				continue
			}

			if sRankerQuery.Next() {
				var count int64

				scanErr := sRankerQuery.Scan(&count)

				sRankerQuery.Close()

				if scanErr != nil {
					return -2, achievedAchievements
				}

				if count >= 5 {
					AchieveAchievement(Achievement5SRanks24Hours)
				}
			}
		case AchievementMostImproved:
			mostImprovedQuery, mostImprovedQueryErr := database.Database.Query("SELECT COUNT(*) AS 'count' FROM waffle.scores WHERE user_id = ? AND beatmap_id = ? AND passed = 1 AND ranking IN ('S', 'X', 'SH', 'XH', 'A') AND date > (SELECT MIN(date) FROM waffle.scores scores2 WHERE scores2.user_id = ? AND scores2.ranking = 'D' AND scores2.beatmap_id = ? AND scores2.passed = 1)", userId, beatmapId, userId, beatmapId)

			if mostImprovedQueryErr != nil {
				continue
			}

			if mostImprovedQuery.Next() {
				var count int64

				scanErr := mostImprovedQuery.Scan(&count)

				mostImprovedQuery.Close()

				if scanErr != nil {
					return -2, achievedAchievements
				}

				if count > 0 {
					AchieveAchievement(AchievementMostImproved)
				}
			}
		case AchievementNonstopDancerParaparaMAX:
			paraparaMaxQuery, paraparaMaxQueryErr := database.Database.Query("SELECT COUNT(*) FROM waffle.scores WHERE user_id = ? AND beatmapset_id = 972 AND score > 3000000 AND (enabled_mods & 1) = 0")

			if paraparaMaxQueryErr != nil {
				continue
			}

			if paraparaMaxQuery.Next() {
				var count int64

				scanErr := paraparaMaxQuery.Scan(&count)

				paraparaMaxQuery.Close()

				if scanErr != nil {
					return -2, achievedAchievements
				}

				if count > 0 {
					AchieveAchievement(AchievementNonstopDancerParaparaMAX)
				}
			}
		}
	}

	achievementIdsString := "("

	for _, currentAchievementId := range achievedIds {
		achievementIdsString += strconv.FormatInt(int64(currentAchievementId), 10) + ","
	}

	achievementIdsString = strings.TrimSuffix(achievementIdsString, ",")
	achievementIdsString += ")"

	getAchievementsQuery, getAchievementsQueryErr := database.Database.Query("SELECT * FROM osu_achievements WHERE achievement_id IN" + achievementIdsString)

	if getAchievementsQueryErr != nil {
		return -2, achievedAchievements
	}

	for getAchievementsQuery.Next() {
		currentAchievement := database.Achievement{}

		scanErr := getAchievementsQuery.Scan(&currentAchievement.AchievementId, &currentAchievement.Image, &currentAchievement.Name)

		if scanErr != nil {
			getAchievementsQuery.Close()

			return -2, achievedAchievements
		}

		achievedAchievements = append(achievedAchievements, currentAchievement)
	}

	if getAchievementsQuery != nil {
		getAchievementsQuery.Close()
	}

	return 0, achievedAchievements
}
