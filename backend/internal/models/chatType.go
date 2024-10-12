package models

import "fmt"

type ChatType int16

const (
	UnknownChatType                       = ChatType(0)
	GeneralChatType                       = ChatType(1)
	JobInterviewSeparateQuestionsChatType = ChatType(2)
)

func GetChatTypeFromNumber[T1 Number](n T1) (ChatType, error) {
	switch n {
	case 0:
		return UnknownChatType, nil
	case 1:
		return GeneralChatType, nil
	case 2:
		return JobInterviewSeparateQuestionsChatType, nil
	default:
		return UnknownChatType, fmt.Errorf("unexpected number %v", n)
	}
}
