package models

type MeaningPosition string

const (
	PositionStraight MeaningPosition = "straight"
	PositionReverted MeaningPosition = "reverted"
)

type MeaningRef struct {
	ID       int64  `json:"id"`
	Position string `json:"position"`
	SourceID int64  `json:"source"`
}
