package models

import (
	"time"
)

type Track struct {
	ID        uint          `json:"id"         gorm:"not null;primarykey;autoIncrement;"`
	Title     string        `json:"title"      gorm:"not null;"`
	SourceURL string        `json:"source_url" gorm:"not null;unique;"`
	Duration  time.Duration `json:"duration"   gorm:"not null;"`
	Artists   []Artist      `json:"artists"    gorm:"not null;many2many:artists_tracks;"`
	CreatedAt time.Time     `json:"created_at" gorm:"not null;"`
	UpdatedAt time.Time     `json:"updated_at" gorm:"not null;"`
}

func ListTracks() ([]Track, error)             { return nil, nil }
func CreateTrack() (*Track, error)             { return nil, nil }
func ReadTrack(id uint) (*Track, error)        { return nil, nil }
func UpdateTrack(track *Track) (*Track, error) { return nil, nil }
func DeleteTrack(id uint) error                { return nil }
