package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type ScoreRecord struct {
	Date         time.Time `json:"date"`
	PlayerRed    int       `json:"playerRed"`
	PlayerYellow int       `json:"playerYellow"`
	BoardSizeX   int       `json:"boardSizeX" validate:"required"`
	BoardSizeY   int       `json:"boardSizeY" validate:"required"`
	HostId       uuid.UUID `json:"hostId" validate:"required"`
}

func (s ScoreRecord) Validate() error {
	return validator.New().Struct(s)
}

func (s ScoreRecord) IsHostWinner() bool {
	return s.PlayerYellow > s.PlayerRed
}
