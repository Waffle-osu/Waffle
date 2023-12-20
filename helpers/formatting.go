package helpers

import "strings"

func FormatMods(mods uint32) string {
	if mods == 0 {
		return "NoMod"
	}

	output := ""

	if (mods & 1) > 1 {
		output += "NoFail, "
	}

	if (mods & 2) > 1 {
		output += "Easy, "
	}

	if (mods & 4) > 1 {
		output += "NoVideo, "
	}

	if (mods & 8) > 1 {
		output += "Hidden, "
	}

	if (mods & 16) > 1 {
		output += "HardRock, "
	}

	if (mods & 32) > 1 {
		output += "SuddenDeath, "
	}

	if (mods & 64) > 1 {
		output += "DoubleTime, "
	}

	if (mods & 128) > 1 {
		output += "Relax, "
	}

	if (mods & 256) > 1 {
		output += "HalfTime, "
	}

	if (mods & 1024) > 1 {
		output += "Flashlight, "
	}

	if (mods & 2048) > 1 {
		output += "Auto, "
	}

	if (mods & 4096) > 1 {
		output += "SpunOut, "
	}

	if (mods & 8192) > 1 {
		output += "Autopilot, "
	}

	return strings.TrimRight(output, ", ")
}

func FormatPlaymodes(playmode uint8) string {
	switch playmode {
	case 0:
		return "osu!"
	case 1:
		return "osu!taiko"
	case 2:
		return "osu!catch"
	case 3:
		return "osu!mania"
	}

	return ""
}

func FormatScoringType(scoringType uint8) string {
	switch scoringType {
	case 0:
		return "Score"
	case 1:
		return "Accuracy"
	}

	return ""
}

func FormatMatchTeamTypes(teamType uint8) string {
	switch teamType {
	case 0:
		return "Head To Head"
	case 1:
		return "Tag Coop"
	case 2:
		return "Team Vs"
	case 3:
		return "Tag Team Vs"
	}

	return ""
}

func FormatSlotStatus(slotStatus uint8) string {
	switch slotStatus {
	case 1:
		return "Open"
	case 2:
		return "Locked"
	case 4:
		return "Not ready"
	case 8:
		return "Ready"
	case 16:
		return "Missing Map"
	case 32:
		return "Playing"
	case 64:
		return "Completed"
	default:
		return ""
	}
}
