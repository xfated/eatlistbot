package constants

// Finite state machine for handling adding items
type State int

const (
	Idle State = iota

	/* #### Adding Item #### */
	ReadyForNextAction
	AddNewSetName
	AddNewSetAddress
	AddNewSetNotes
	AddNewSetURL
	AddNewSetImages
	AddNewSetTags
	AddNewRemoveTags
	ConfirmAddItemSubmit
	/* ######## */

	/* #### Query #### */
	QuerySelectType

	QueryOneTagOrName
	QueryOneSetName

	QuerySetTags
	QueryFewSetNum
	QueryRetrieve
	/* ######## */

	/* #### Delete Item #### */
	DeleteSelect
	DeleteConfirm
	/* ######## */

	/* #### EditItem #### */
	GetItemToEdit
	/* ######## */
)

type ItemDetails struct {
	Name    string          `json:"name"`
	Address string          `json:"address"`
	Notes   string          `json:"notes"`
	URL     string          `json:"url"`
	Images  map[string]bool `json:"images"`
	Tags    map[string]bool `json:"tags"`
}

func (itemData *ItemDetails) GetImageIDs() []string {
	if itemData.Images == nil {
		return []string{}
	}
	imageIDs := make([]string, 0)
	for id := range itemData.Images {
		imageIDs = append(imageIDs, id)
	}
	return imageIDs
}

func IsAddingNewItem(state State) bool {
	switch state {
	case ReadyForNextAction,
		AddNewSetName,
		AddNewSetAddress,
		AddNewSetNotes,
		AddNewSetURL,
		AddNewSetImages,
		AddNewSetTags,
		AddNewRemoveTags,
		ConfirmAddItemSubmit:
		return true
	default:
		return false
	}
}

func IsQuery(state State) bool {
	switch state {
	case QuerySelectType,

		QueryOneTagOrName,
		QueryOneSetName,

		QueryFewSetNum,

		QuerySetTags,
		QueryRetrieve:
		return true
	default:
		return false
	}
}

func IsDeleteItem(state State) bool {
	switch state {
	case DeleteSelect,
		DeleteConfirm:
		return true
	default:
		return false
	}
}

func IsEditItem(state State) bool {
	switch state {
	case GetItemToEdit:
		return true
	default:
		return false
	}
}
