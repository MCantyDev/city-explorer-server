package models

import (
	"encoding/json"
	"time"
)

type User struct {
	Id        uint      `gorm:"primaryKey;autoIncrement"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`
	Username  string    `gorm:"unique;not null"`
	Email     string    `gorm:"unique; not null"`
	Password  string    `gorm:"not null"`
	IsAdmin   bool      `gorm:"type:boolean"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Country struct {
	Id         uint            `gorm:"primaryKey;autoIncrement"`
	Name       string          `gorm:"not null"`
	IsoCode    string          `gorm:"not null"`
	Data       json.RawMessage `gorm:"not null"`
	CreatedAt  time.Time       `gorm:"autoCreateTime"`
	UpdatedAt  time.Time       `gorm:"autoUpdateTime"`
	ExpiryDate time.Time       `gorm:"type:timestamp"`
}

type City struct {
	Id        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type CityWeather struct {
	Id          uint            `gorm:"primaryKey;autoIncrement"`
	CityName    string          `gorm:"not null"`
	CountryName string          `gorm:"not null"`
	Lat         float64         `gorm:"type:decimal(9,6);not null"`
	Lon         float64         `gorm:"type:decimal(9,6);not null"`
	Data        json.RawMessage `gorm:"type:json;not null"`
	CreatedAt   time.Time       `gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime"`
	ExpiryDate  time.Time       `gorm:"type:timestamp"`
}

type CitySights struct {
	Id          uint            `gorm:"primaryKey;autoIncrement"`
	CityName    string          `gorm:"not null"`
	CountryName string          `gorm:"not null"`
	Lat         float64         `gorm:"not null"`
	Lon         float64         `gorm:"not null"`
	Data        json.RawMessage `gorm:"not null"`
	CreatedAt   time.Time       `gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime"`
	ExpiryDate  time.Time       `gorm:"type:timestamp"`
}

type CityPoi struct {
	Id          uint            `gorm:"primaryKey;autoIncrement"`
	CityName    string          `gorm:"not null"`
	CountryName string          `gorm:"not null"`
	Xid         string          `gorm:"not null"`
	Data        json.RawMessage `gorm:"not null"`
	CreatedAt   time.Time       `gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime"`
	ExpiryDate  time.Time       `gorm:"type:timestamp"`
}
