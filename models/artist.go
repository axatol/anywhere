package models

import (
	"errors"
	"time"

	"github.com/tunes-anywhere/anywhere/database"
	"gorm.io/gorm"
)

type PartialArtist struct {
	Name string `json:"name"`
}

type Artist struct {
	ID        uint      `json:"id"         gorm:"not null;primarykey;autoIncrement;"`
	Name      string    `json:"name"       gorm:"not null;"`
	Tracks    []Track   `json:"tracks"     gorm:"not null;many2many:artists_tracks;"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;"`
}

func ListArtists() ([]Artist, error) {
	var artists []Artist
	err := database.DB.Raw("select * from artists").Scan(&artists).Error

	if err != nil {
		return nil, err
	}

	return artists, nil
}

func CreateArtist(name string) (*Artist, error) {
	artist := Artist{Name: name}
	err := database.DB.Create(&artist).Error

	if err != nil {
		return nil, err
	}

	return &artist, nil
}

func ReadArtist(id uint) (*Artist, error) {
	artist := Artist{ID: id}
	err := database.DB.First(&artist).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &artist, nil
}

func UpdateArtist(id uint, partialArtist *PartialArtist) (*Artist, error) {
	artist := Artist{ID: id}
	err := database.DB.Model(&artist).Select("name").Update("name", partialArtist.Name).Error
	if err != nil {
		return nil, err
	}

	return &artist, nil
}

func DeleteArtist(id uint) error {
	err := database.DB.Delete(&Artist{ID: id}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}
