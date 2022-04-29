package database

type UpdaterItem struct {
	ItemId         int64
	ServerFilename string
	ClientFilename string
	FileHash       string
	ItemName       string
	ItemAction     string
}

func (item *UpdaterItem) FormatUpdaterItem() string {
	return item.ServerFilename + " " + item.FileHash + " " + item.ItemName + " " + item.ItemAction + " " + item.ClientFilename + "\n"
}

func GetUpdaterItems() (result int8, items []UpdaterItem) {
	queryResult, queryErr := database.Query("SELECT item_id, server_filename, client_filename, file_hash, item_name, item_action FROM waffle.updater_items")

	if queryErr != nil {
		return -1, nil
	}

	queryItems := []UpdaterItem{}

	for queryResult.Next() {
		item := UpdaterItem{}

		scanErr := queryResult.Scan(&item.ItemId, &item.ServerFilename, &item.ClientFilename, &item.FileHash, &item.ItemName, &item.ItemAction)

		queryItems = append(queryItems, item)

		if scanErr != nil {
			return -1, nil
		}
	}

	return 0, queryItems
}

func GetOsuExecutableHash() string {
	queryResult, queryErr := database.Query("SELECT server_filename, file_hash FROM updater_items WHERE server_filename = 'osu!.exe'")

	if queryErr != nil {
		return ""
	}

	if queryResult.Next() {
		var serverFilename string
		var fileHash string

		scanErr := queryResult.Scan(&serverFilename, &fileHash)

		if scanErr != nil {
			return ""
		}

		return fileHash
	} else {
		return ""
	}
}
