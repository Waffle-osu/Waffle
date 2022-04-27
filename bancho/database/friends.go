package database

type FriendEntry struct {
	User1 uint64
	User2 uint64
}

// GetFriendsList retrieves all friends from the user with ID `userId`
func GetFriendsList(userId uint64) (result int, friendsList []FriendEntry) {
	var friends = []FriendEntry{}

	queryResult, queryErr := database.Query("SELECT user_1, user_2 FROM waffle.friends WHERE user_1 = ?", userId)

	if queryErr != nil {
		return -1, friends
	}

	for queryResult.Next() {
		friendEntry := FriendEntry{}

		queryResult.Scan(&friendEntry.User1, &friendEntry.User2)

		friends = append(friends, friendEntry)
	}

	return 0, friends
}

// AddFriend stores a new friendship in the database
func AddFriend(userId uint64, friendId uint64) bool {
	_, queryErr := database.Query("INSERT INTO waffle.friends (user_1, user_2) VALUES (?, ?)", userId, friendId)

	if queryErr != nil {
		return false
	}

	return true
}

// RemoveFriend removes a friendship from the database
func RemoveFriend(userId uint64, friendId uint64) bool {
	_, queryErr := database.Query("DELETE FROM waffle.friends WHERE user_1 = ? AND user_2 = ?", userId, friendId)

	if queryErr != nil {
		return false
	}

	return true
}
