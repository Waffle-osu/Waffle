package database

func ScreenshotsHitScreenshotLimit(id uint64) bool {
	query, queryErr := Database.Query("SELECT COUNT(*) AS 'count' FROM waffle.screenshots WHERE id = ?", id)

	if queryErr != nil {
		if query != nil {
			query.Close()
		}
		return true
	}

	if query.Next() {
		var count uint64

		scanErr := query.Scan(&count)

		if scanErr != nil {
			query.Close()
			return true
		}

		query.Close()

		return count >= 128
	}

	query.Close()

	return true
}

func ScreenshotsInsertNewScreenshot(userId uint64, filename string) bool {
	query, queryErr := Database.Query("INSERT INTO waffle.screenshots (user_id, filename) VALUES (?, ?)", userId, filename)
	defer query.Close()

	return queryErr == nil
}
