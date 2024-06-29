package models

type MessageType uint8

const (
	UnknownMessageType   = MessageType(0)
	JobInterviewQuestion = MessageType(1)
)

type PreparedMessage struct {
	Type        MessageType `firestore:"typ"`
	GermanText  string      `firestore:"de_txt"`
	GermanAudio string      `firestore:"de_audio"`
}
