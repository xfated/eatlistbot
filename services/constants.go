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
