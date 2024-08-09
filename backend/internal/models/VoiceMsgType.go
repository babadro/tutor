package models

type VoiceMsgType int16

const (
	DefaultVoiceMsgType            = VoiceMsgType(0)
	AwaitingCompletionVoiceMsgType = VoiceMsgType(1)
)
