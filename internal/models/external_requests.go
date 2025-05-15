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
	Current         json.RawMessage `json:"current"`
	Daily           json.RawMessage `json:"daily"`
}

type OpenTripRequest struct {
	Type     string          `json:"type"`
	Features json.RawMessage `json:"features"`
}
type OpenTripPlaceRequest struct {
	XID  string `json:"xid"`
	Name string `json:"name"`

	Address struct {
		City          string `json:"city"`
		State         string `json:"state"`
		County        string `json:"county"`
		Suburb        string `json:"suburb"`
		Country       string `json:"country"`
		Postcode      string `json:"postcode"`
		Pedestrian    string `json:"pedestrian"`
		CountryCode   string `json:"country_code"`
		StateDistrict string `json:"state_district"`
	} `json:"address"`

	Rate string `json:"rate"`
	OSM  string `json:"osm"`

	Bbox struct {
		LonMin float64 `json:"lon_min"`
		LonMax float64 `json:"lon_max"`
		LatMin float64 `json:"lat_min"`
		LatMax float64 `json:"lat_max"`
	} `json:"bbox"`

	Wikidata string `json:"wikidata"`
	Kinds    string `json:"kinds"`

	Sources struct {
		Geometry   string   `json:"geometry"`
		Attributes []string `json:"attributes"`
	} `json:"sources"`

	OTM       string `json:"otm"`
	Wikipidia string `json:"wikipedia"`

	Preview struct {
		Source string `json:"source"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"preview"`

	WikipidiaExtracts struct {
		Title string `json:"title"`
		Text  string `json:"text"`
		Html  string `json:"html"`
	} `json:"wikipedia_extracts"`

	Point struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"point"`
}
