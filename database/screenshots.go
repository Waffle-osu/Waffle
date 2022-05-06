package database

func ScreenshotsHitScreenshotLimit(id uint64) bool {
	query, queryErr := database.Query("SELECT COUNT(*) AS 'count' FROM waffle.screenshots WHERE id = ?", id)
	defer query.Close()

	if queryErr != nil {
		return true
	}

	if query.Next() {
		var count uint64

		scanErr := query.Scan(&count)

		if scanErr != nil {
			return true
		}

		return count >= 128
	}

	return true
}

func ScreenshotsInsertNewScreenshot(userId uint64, filename string) bool {
	query, queryErr := database.Query("INSERT INTO waffle.screenshots (id, filename) VALUES (?, ?)", userId, filename)
	defer query.Close()

	return queryErr == nil
}
