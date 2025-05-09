package models

import "encoding/json"

type PhotonRequest struct {
	Type     string          `json:"type"`
	Features json.RawMessage `json:"features"`
}

type RestCountriesRequest []json.RawMessage

type OpenWeatherRequest struct {
}

type OpenTripRequest struct {
}
