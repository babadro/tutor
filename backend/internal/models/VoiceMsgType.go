package models

type VoiceMsgType int16

const (
	DefaultVoiceMsgType            = VoiceMsgType(1)
	AwaitingCompletionVoiceMsgType = VoiceMsgType(2)
)
