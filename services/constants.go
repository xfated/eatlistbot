package services

// Finite state machine for handling adding items
type State int

const (
	Idle State = iota
	SetName
	SetAddress
	SetURL
	SetImages
	SetTags
	Finished
)

type PlaceDetails struct {
	Name    string          `json:"name"`
	Address string          `json:"address"`
	URL     string          `json:"url"`
	Images  map[string]bool `json:"images"`
	Tags    map[string]bool `json:"tags"`
}

func getNextState(cur State) State {
	switch cur {
	case Idle, Finished:
		return cur
	case SetName, SetAddress, SetURL, SetImages, SetTags:
		return cur + 1
	default:
		return cur
	}
}
