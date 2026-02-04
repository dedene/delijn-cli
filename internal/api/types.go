//nolint:tagliatelle // JSON tags must match De Lijn API response format (Dutch camelCase)
package api

import (
	"slices"
	"time"
)

// Entity represents a regional entity (1-5: Antwerpen, Oost-Vlaanderen, etc.)
type Entity struct {
	Number      int    `json:"entiteitnummer"`
	Code        string `json:"code"`
	Description string `json:"omschrijving"`
}

// Stop represents a De Lijn stop (halte).
type Stop struct {
	Number           int       `json:"haltenummer"`
	Description      string    `json:"omschrijving"`
	Municipality     string    `json:"omschrijvingGemeente"`
	MunicipalityCode int       `json:"gemeentenummer"`
	EntityNumber     int       `json:"entiteitnummer"`
	GeoCoordinate    *GeoCoord `json:"geoCoordinaat,omitempty"`
	Links            []Link    `json:"links,omitempty"`
}

// GeoCoord represents geographic coordinates in Lambert72 and WGS84.
type GeoCoord struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Link represents an API hyperlink.
type Link struct {
	Rel string `json:"rel"`
	URL string `json:"url"`
}

// Line represents a De Lijn line (lijn).
type Line struct {
	EntityNumber  int    `json:"entiteitnummer"`
	LineNumber    int    `json:"lijnnummer"`
	PublicNumber  string `json:"lijnnummerPubliek"`
	Description   string `json:"omschrijving"`
	TransportType string `json:"vervoertype"` // BUS, TRAM, METRO
	IsPublic      bool   `json:"publpiekeVervoer"`
	Links         []Link `json:"links,omitempty"`
}

// LineDirection represents a line direction (richting).
type LineDirection struct {
	EntityNumber int    `json:"entiteitnummer"`
	LineNumber   int    `json:"lijnnummer"`
	Direction    string `json:"richting"` // HEEN, TERUG
	Destination  string `json:"bestemming"`
}

// Departure represents a realtime departure at a stop.
type Departure struct {
	EntityNumber     int        `json:"entiteitnummer"`
	LineNumber       int        `json:"lijnnummer"`
	LinePublicNumber string     `json:"lijnnummerPubliek,omitempty"`
	Direction        string     `json:"richting"`
	Destination      string     `json:"bestemming"`
	ScheduledTime    time.Time  `json:"-"`
	RealTime         *time.Time `json:"-"`
	ScheduledTimeRaw string     `json:"dienstregelingTijdstip"`
	RealTimeRaw      string     `json:"real-timeTijdstip,omitempty"`
	PredictionStatus []string   `json:"predictionStatussen"`
	TransportType    string     `json:"vervoertype,omitempty"`
}

// IsRealTime returns whether this departure has realtime data.
func (d *Departure) IsRealTime() bool {
	return slices.Contains(d.PredictionStatus, "REALTIME")
}

// DelaySeconds returns the delay in seconds (positive = late, negative = early).
func (d *Departure) DelaySeconds() int {
	if d.RealTime == nil {
		return 0
	}

	return int(d.RealTime.Sub(d.ScheduledTime).Seconds())
}

// StopPassage represents passages at a specific stop.
type StopPassage struct {
	StopNumber int         `json:"haltenummer"`
	Departures []Departure `json:"doorkomsten"`
}

// RealtimeResponse is the response from the real-time endpoint.
type RealtimeResponse struct {
	StopPassages []StopPassage `json:"halteDoorkomsten"`
}

// StopsResponse is the response from stops search.
type StopsResponse struct {
	Stops []Stop `json:"haltes"`
	Links []Link `json:"links,omitempty"`
}

// LinesResponse is the response from lines search.
type LinesResponse struct {
	Lines []Line `json:"lijnen"`
	Links []Link `json:"links,omitempty"`
}

// Colour represents a colour code with hex value.
type Colour struct {
	Code string `json:"code"`
	Hex  string `json:"hex"`
}

// LineColours represents line foreground/background colours.
type LineColours struct {
	Foreground       Colour `json:"voorgrond"`
	Background       Colour `json:"achtergrond"`
	ForegroundBorder Colour `json:"voorgrondRand"`
	BackgroundBorder Colour `json:"achtergrondRand"`
}

// Disruption represents a service disruption (storing/omleiding).
type Disruption struct {
	ID           string    `json:"id"`
	Title        string    `json:"titel"`
	Description  string    `json:"omschrijving"`
	Type         string    `json:"type"` // STORING, OMLEIDING
	StartDate    time.Time `json:"-"`
	EndDate      time.Time `json:"-"`
	StartDateRaw string    `json:"startDatum"`
	EndDateRaw   string    `json:"eindDatum"`
	Lines        []Line    `json:"lijnen,omitempty"`
}

// DisruptionsResponse is the response from disruptions endpoint.
type DisruptionsResponse struct {
	Disruptions []Disruption `json:"storingen"`
}
