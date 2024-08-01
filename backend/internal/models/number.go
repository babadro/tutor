package models

type Number interface {
	uint8 | int8 | uint16 | int16 | uint32 | int32 | uint64 | int64 | float32 | float64 | ChatType
}
