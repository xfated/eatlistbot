package services

// Finite state machine for handling adding items
type State int

const (
	Idle State = iota

	ReadyForNextAction
	SetName
	SetAddress
	SetURL
	SetImages
	SetTags
	Preview
	FinishAdding

	Finished
)

type PlaceDetails struct {
	Name    string          `json:"name"`
	Address string          `json:"address"`
	URL     string          `json:"url"`
	Images  map[string]bool `json:"images"`
	Tags    map[string]bool `json:"tags"`
}

func IsAddingNewPlace(state State) bool {
	switch state {
	case SetName,
		SetAddress,
		SetImages,
		SetTags,
		SetURL:
		return true
	default:
		return false
	}
}
