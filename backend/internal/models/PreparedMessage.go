package models

type MessageType uint8

const (
	UnknownMessageType   = MessageType(0)
	JobInterviewQuestion = MessageType(1)
)

// Variation represents a translation or variation of the message
type Variation struct {
	Language string `firestore:"lang"` // The language code, e.g., "de", "fr"
	Text     string `firestore:"txt"`
	Audio    string `firestore:"audio"`
	ID       string `firestore:"id"`
}

type PreparedMessage struct {
	Type       MessageType `firestore:"typ"`
	BaseText   string      `firestore:"base_text"` // Basic form of the message
	Variations []Variation `firestore:"variations"`
}
