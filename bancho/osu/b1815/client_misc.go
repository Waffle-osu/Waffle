package b1815

import (
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
	"fmt"
	"strings"
)

// Handles the OsuBeatmapInfoRequest packet.
func (client *Client) HandleBeatmapInfoRequest(infoRequest base_packet_structures.BeatmapInfoRequest) {
	go func() {
		infoReply := base_packet_structures.BeatmapInfoReply{}

		//Initially store the user ids for the prepared statement
		queryArguments := []interface{}{
			client.OsuStats.UserID, client.OsuStats.UserID, client.OsuStats.UserID,
		}

		//will store the prepared statement question marks for the filenames
		questionMarks := ""

		//edge case for no filenames, it will still have at least 1 filename, even though itll be empty
		if len(infoRequest.Filenames) == 0 {
			questionMarks = "?"
			queryArguments = append(queryArguments, "")
		} else {
			//for every filename add a question mark
			for i := 0; i != len(infoRequest.Filenames); i++ {
				questionMarks += "?, "
			}
		}

		//trim off the last comma to not cause massive issues
		questionMarks = strings.TrimSuffix(questionMarks, ", ")

		//will store the beatmap ids for the sql
		beatmapIds := ""

		//edge case for no beatmap ids, it will still have at least 1 beatmap id, even though itll be 0
		if len(infoRequest.BeatmapIds) == 0 {
			beatmapIds = "0"
		} else {
			//add every beatmap id
			for i := 0; i != len(infoRequest.BeatmapIds); i++ {
				beatmapIds += string(infoRequest.BeatmapIds[i]) + ", "
			}
		}

		//trim off comma to not have a extra one
		beatmapIds = strings.TrimSuffix(beatmapIds, ", ")

		sqlString := `
SELECT				
	result.beatmap_id,				
	result.beatmapset_id,
	result.filename,
	result.beatmap_md5,
	result.ranking_status,
	result.final_osu_ranking AS 'osu_ranking',
	result.final_taiko_ranking AS 'taiko_ranking',
	result.final_catch_ranking AS 'catch_ranking'
FROM (
	SELECT beatmaps.beatmap_id, 
			beatmaps.beatmapset_id, 
			beatmaps.filename, 
			beatmaps.beatmap_md5, 
			beatmaps.ranking_status, 
			osuResult.ranking AS 'osu_ranking', 
			osuResult.user_id AS 'osu_user_id', 
			taikoResult.ranking AS 'taiko_ranking',
			taikoResult.user_id AS 'taiko_user_id', 
			catchResult.ranking AS 'catch_ranking', 
			catchResult.user_id AS 'catch_user_id',
		CASE WHEN osuResult.ranking IS NULL THEN 'N' ELSE osuResult.ranking END AS 'final_osu_ranking',
		CASE WHEN taikoResult.ranking IS NULL THEN 'N' ELSE taikoResult.ranking END AS 'final_taiko_ranking', 
		CASE WHEN catchResult.ranking IS NULL THEN 'N' ELSE catchResult.ranking END AS 'final_catch_ranking'
	FROM waffle.beatmaps 
		LEFT JOIN scores osuResult ON osuResult.beatmap_id = beatmaps.beatmap_id AND osuResult.mapset_best = 1 AND osuResult.playmode = 0 AND osuResult.user_id = ? 
		LEFT JOIN scores taikoResult ON taikoResult.beatmap_id = beatmaps.beatmap_id AND taikoResult.mapset_best = 1 AND taikoResult.playmode = 1 AND taikoResult.user_id = ? 
		LEFT JOIN scores catchResult ON catchResult.beatmap_id = beatmaps.beatmap_id AND catchResult.mapset_best = 1 AND catchResult.playmode = 2 AND catchResult.user_id = ? 
	WHERE beatmaps.filename IN ( %s ) 
	OR beatmaps.beatmap_id IN ( %s )
) result
`
		//the absolutely gigantic sql
		sql := fmt.Sprintf(sqlString, questionMarks, beatmapIds)

		//add the filenames as query arguments
		for i := 0; i != len(infoRequest.Filenames); i++ {
			queryArguments = append(queryArguments, infoRequest.Filenames[i])
		}

		//query
		databaseQuery, databaseQueryErr := database.Database.Query(sql, queryArguments...)

		//momentarily nonexistant error handling
		if databaseQueryErr != nil {
			if databaseQuery != nil {
				databaseQuery.Close()
			}
		}

		if databaseQuery != nil {
			//for every database result
			for databaseQuery.Next() {
				beatmapInfo := base_packet_structures.BeatmapInfo{}

				var osuRank, taikoRank, catchRank string
				var osuFilename string

				scanErr := databaseQuery.Scan(&beatmapInfo.BeatmapId, &beatmapInfo.BeatmapSetId, &osuFilename, &beatmapInfo.BeatmapChecksum, &beatmapInfo.Ranked, &osuRank, &taikoRank, &catchRank)

				if scanErr != nil {
					continue
				}

				rankToNum := map[string]uint8{
					"XH": 0,
					"SH": 1,
					"X":  2,
					"S":  3,
					"A":  4,
					"B":  5,
					"C":  6,
					"D":  7,
					"F":  8,
					"N":  9,
				}

				//convert string rank to peppys enum ranking
				beatmapInfo.OsuRank = rankToNum[osuRank]
				beatmapInfo.TaikoRank = rankToNum[taikoRank]
				beatmapInfo.CatchRank = rankToNum[catchRank]

				//will store the index of the info request the client gave us
				infoPosition := int16(-1)

				//find it
				for k := 0; k != len(infoRequest.Filenames); k++ {
					if infoRequest.Filenames[k] == osuFilename {
						infoPosition = int16(k)
					}
				}

				beatmapInfo.InfoId = infoPosition

				//append to the reply list
				infoReply.BeatmapInfos = append(infoReply.BeatmapInfos, beatmapInfo)
			}

			infoReply.Count = int32(len(infoReply.BeatmapInfos))

			//send off
			client.BanchoBeatmapInfoReply(infoReply)

			//make sure to close the connection
			databaseQuery.Close()
		}
	}()
}
