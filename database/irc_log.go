package database

func ChatInsertNewMessage(userId uint64, target string, message string) {
	query, _ := Database.Query("INSERT INTO waffle.irc_log (sender, target, message) VALUES (?, ?, ?)", userId, target, message)
	defer query.Close()
}
