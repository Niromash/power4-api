package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type GameRecord struct {
	GameId       uuid.UUID `json:"gameId"`
	Date         time.Time `json:"date"`
	PlayerRed    int       `json:"playerRed"`
	PlayerYellow int       `json:"playerYellow"`
	BoardSizeX   int       `json:"boardSizeX" validate:"required"`
	BoardSizeY   int       `json:"boardSizeY" validate:"required"`
	HostId       uuid.UUID `json:"hostId" validate:"required"`
	Boards       []Board   `json:"boards" validate:"required"`
}

type Board struct {
	Board        [][]int `json:"board" validate:"required"`
	Date         string  `json:"date" validate:"required"`
	PlayerRed    int     `json:"playerRed"`
	PlayerYellow int     `json:"playerYellow"`
}

func (s GameRecord) Validate() error {
	return validator.New().Struct(s)
}

func (s GameRecord) IsHostWinner() bool {
	return s.PlayerYellow > s.PlayerRed
}
