package constants

import "strings"

type Adventure int

const (
	AdventurePond Adventure = iota
	AdventureStream
	AdventureSwamp
	AdventureRiver
	AdventureForest
	AdventureGreatLake
)
const (
	MAX_RATING_POND       = 10
	MAX_RATING_STREAM     = 10
	MAX_RATING_SWAMP      = 10
	MAX_RATING_RIVER      = 10 * 10
	MAX_RATING_FOREST     = 10 * 10 * 10
	MAX_RATING_GREAT_LAKE = 10 * 10 * 10 * 10
)

func (adventure Adventure) String() string {
	switch adventure {
	case AdventurePond:
		return "pond"
	case AdventureStream:
		return "stream"
	case AdventureSwamp:
		return "swamp"
	case AdventureRiver:
		return "river"
	case AdventureForest:
		return "forest"
	case AdventureGreatLake:
		return "great-lake"
	default:
		return ""
	}
}

func AdventureFromString(adventure string) Adventure {
	lowerCased := strings.ToLower(adventure)

	switch lowerCased {
	case "pond":
		return AdventurePond
	case "stream":
		return AdventureStream
	case "swamp":
		return AdventureSwamp
	case "river":
		return AdventureRiver
	case "forest":
		return AdventureForest
	case "great-lake":
		return AdventureGreatLake
	default:
		return AdventurePond
	}
}
