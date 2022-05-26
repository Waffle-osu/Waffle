package database

type FriendEntry struct {
	User1 uint64
	User2 uint64
}

// FriendsGetFriendsList retrieves all friends from the user with ID `userId`
func FriendsGetFriendsList(userId uint64) (result int, friendsList []FriendEntry) {
	var friends = []FriendEntry{}

	queryResult, queryErr := Database.Query("SELECT user_1, user_2 FROM waffle.friends WHERE user_1 = ?", userId)

	if queryErr != nil {
		if queryResult != nil {
			queryResult.Close()
		}

		return -1, friends
	}

	for queryResult.Next() {
		friendEntry := FriendEntry{}

		scanErr := queryResult.Scan(&friendEntry.User1, &friendEntry.User2)

		if scanErr != nil {
			queryResult.Close()

			return -2, friends
		}

		friends = append(friends, friendEntry)
	}

	if queryResult != nil {
		queryResult.Close()
	}

	return 0, friends
}

// FriendsAddFriend stores a new friendship in the Database
func FriendsAddFriend(userId uint64, friendId uint64) bool {
	query, queryErr := Database.Query("INSERT INTO waffle.friends (user_1, user_2) VALUES (?, ?)", userId, friendId)

	if queryErr != nil {
		if query != nil {
			query.Close()
		}
		return false
	}

	query.Close()
	return true
}

// FriendsRemoveFriend removes a friendship from the Database
func FriendsRemoveFriend(userId uint64, friendId uint64) bool {
	query, queryErr := Database.Query("DELETE FROM waffle.friends WHERE user_1 = ? AND user_2 = ?", userId, friendId)

	if queryErr != nil {
		return false
	}

	if query != nil {
		query.Close()
	}

	return true
}
