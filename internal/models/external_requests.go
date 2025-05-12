package models

import "encoding/json"

// Get ALL information from the Request
type PhotonRequest struct {
	Type     string          `json:"type"`
	Features json.RawMessage `json:"features"`
}

type RestCountriesRequest []json.RawMessage

type OpenWeatherRequest struct {
	Lat             float64         `json:"lat"`
	Long            float64         `json:"lon"`
	Timezone        string          `json:"timezone"`
	Timezone_Offset float64         `json:"timezone_offset"`
	Hourly          json.RawMessage `json:"hourly"`
	Daily           json.RawMessage `json:"daily"`
}

type OpenTripRequest struct {
	Type     string          `json:"type"`
	Features json.RawMessage `json:"features"`
}
