package database

type MatchHistoryEventType uint8

const (
	MatchHistoryEventTypeJoin             MatchHistoryEventType = 0
	MatchHistoryEventTypeLeave            MatchHistoryEventType = 1
	MatchHistoryEventTypeKick             MatchHistoryEventType = 2
	MatchHistoryEventTypeMove             MatchHistoryEventType = 3
	MatchHistoryEventTypeLock             MatchHistoryEventType = 4
	MatchHistoryEventTypeUnlock           MatchHistoryEventType = 5
	MatchHistoryEventTypeReady            MatchHistoryEventType = 6
	MatchHistoryEventTypeUnready          MatchHistoryEventType = 7
	MatchHistoryEventTypeChangeTeam       MatchHistoryEventType = 8
	MatchHistoryEventTypeHostChange       MatchHistoryEventType = 9
	MatchHistoryEventTypeSettingsChanged  MatchHistoryEventType = 10
	MatchHistoryEventTypeModsChanged      MatchHistoryEventType = 11
	MatchHistoryEventTypeMatchStarted     MatchHistoryEventType = 12
	MatchHistoryEventTypePlayingStarted   MatchHistoryEventType = 13
	MatchHistoryEventTypeFinalScore       MatchHistoryEventType = 14
	MatchHistoryEventTypePlayerFail       MatchHistoryEventType = 15
	MatchHistoryEventTypeTeamFail         MatchHistoryEventType = 16
	MatchHistoryEventTypeMatchComplete    MatchHistoryEventType = 17
	MatchHistoryEventTypeMatchDisbanded   MatchHistoryEventType = 18
	MatchHistoryEventTypeMatchRefLocked   MatchHistoryEventType = 19
	MatchHistoryEventTypeMatchRefUnlocked MatchHistoryEventType = 20
)

type MatchHistoryElement struct {
	EventId        uint64
	MatchId        string
	EventType      MatchHistoryEventType
	EventInitiator uint64
	ExtraInfo      string
}

func GetMatchHistory(matchId string) (queryResult int8, data []MatchHistoryElement) {
	output := []MatchHistoryElement{}

	query, queryErr := Database.Query("SELECT * FROM osu_match_history WHERE match_id = ?", matchId)

	if queryErr != nil {
		if query != nil {
			query.Close()
		}

		return -2, output
	}

	for query.Next() {
		var result MatchHistoryElement

		scanErr := query.Scan(&result.EventId, &result.MatchId, &result.EventType, &result.EventInitiator, &result.ExtraInfo)

		if scanErr != nil {
			if query != nil {
				query.Close()
			}

			return -1, output
		}

		output = append(output, result)
	}

	query.Close()

	return 0, output
}

func LogMatchHistory(element MatchHistoryElement) bool {
	query, queryErr := Database.Query("INSERT INTO osu_match_history (match_id, event_type, event_initiator_id, extra_info) VALUES (?, ?, ?, ?)", element.MatchId, element.EventType, element.EventInitiator, element.ExtraInfo)

	if queryErr != nil {
		if query != nil {
			query.Close()
		}

		return false
	}

	if query != nil {
		query.Close()
	}

	return true
}
