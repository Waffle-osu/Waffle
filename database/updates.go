package database

type UpdaterItem struct {
	ItemId         int64
	ServerFilename string
	ClientFilename string
	FileHash       string
	ItemName       string
	ItemAction     string
}

//So nobody has to decipher the peppy updater ever again, update2.txt is expected to be formatted like this:
//psa: this is seperated by spaces
//[0]: Server Filename, this is the filename the updater will try to retrieve if it needs to
//[1]: File Hash, this is the expected file hash
//[2]: Item Name, internally called description, You-can-do-dashes-like-this-to-do-spaces
//[3]: Item Action, what the updater should do with the file:
//   :	"zip"  : expects a .zip file, this file will be unzipped after the updater downloads it
//   :	"noup" : if the file exists, good, keep it there
//   :	"del"  : if the file exists, cool, delete it
//   :	"extra": this is then added to the extras menu in the updater
//   :	"none" : just download and be good
//   :	"diff" : this will try to use a diff file to patch the file directly, peppy uses a special format for this and i don't want to delve into that
//[4]: Client Filename, this is the file the updater will check against to see if it's up-to-date

func (item *UpdaterItem) FormatUpdaterItem() string {
	return item.ServerFilename + " " + item.FileHash + " " + item.ItemName + " " + item.ItemAction + " " + item.ClientFilename + "\n"
}

func GetUpdaterItems() (result int8, items []UpdaterItem) {
	queryResult, queryErr := database.Query("SELECT item_id, server_filename, client_filename, file_hash, item_name, item_action FROM waffle.updater_items")
	defer queryResult.Close()

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

func UpdaterHashFromFilename(filename string) string {
	queryResult, queryErr := database.Query("SELECT server_filename, file_hash FROM updater_items WHERE server_filename = ?", filename)
	defer queryResult.Close()

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
