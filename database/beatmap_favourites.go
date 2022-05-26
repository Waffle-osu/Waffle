package database

type BeatmapFavourite struct {
	FavouriteId  uint64
	BeatmapSetId int32
	UserId       uint64
}

func GetUserFavourites(userId uint64) (result int8, favourites []BeatmapFavourite) {
	var beatmapFavourites = []BeatmapFavourite{}

	queryResult, queryErr := Database.Query("SELECT user_id, beatmapset_id FROM waffle.beatmap_favourites WHERE user_id = ?", userId)

	defer queryResult.Close()

	if queryErr != nil {
		if queryResult != nil {
			queryResult.Close()
		}

		return -1, beatmapFavourites
	}

	for queryResult.Next() {
		beatmapEntry := BeatmapFavourite{}

		queryResult.Scan(&beatmapEntry.UserId, &beatmapEntry.BeatmapSetId)

		beatmapFavourites = append(beatmapFavourites, beatmapEntry)
	}

	if queryResult != nil {
		queryResult.Close()
	}

	return 0, beatmapFavourites
}

func FavouritesAddFavourite(userId uint64, beatmapsetId int32) bool {
	query, queryErr := Database.Query("INSERT INTO waffle.beatmap_favourites (user_id, beatmapset_id) VALUES (?, ?)", userId, beatmapsetId)

	if queryErr != nil {
		if query != nil {
			query.Close()
		}
		return false
	}

	query.Close()
	return true
}

func FavouritesRemoveFavourite(userId uint64, beatmapsetId int32) bool {
	query, queryErr := Database.Query("DELETE FROM waffle.beatmap_favourites WHERE user_id = ? AND beatmapset_id = ?", userId, beatmapsetId)

	if queryErr != nil {
		return false
	}

	if query != nil {
		query.Close()
	}

	return true
}
