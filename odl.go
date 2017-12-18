package odl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Basic API Endpoint
const (
	baseURL string = "https://odlinfo.bfs.de/daten/json"
)

// State Status of the measuring point
type State int

// Constants for the status-code of a measuring point defined by the Bundesamt für Strahlschutz
// https://info.bfs.de/downloads/Datenbereitstellung-2016-04-21.pdf
const (
	BROKEN      State = 0
	OK          State = 1
	TEST        State = 128
	MAINTENANCE State = 2048
)

// MgmtNodeID management-node for a measuring point
type MgmtNodeID int

// Management-node constants
const (
	FREIBURG   MgmtNodeID = 1
	BERLIN     MgmtNodeID = 2
	MUENCHEN   MgmtNodeID = 3
	BONN       MgmtNodeID = 4
	SALZGITTER MgmtNodeID = 5
	RENDSBURG  MgmtNodeID = 6
)

// StationInfo StationInfo for a measuring point
type StationInfo struct {
	ID         string     `json:"kenn"`
	Place      string     `json:"ort"`
	Zip        string     `json:"plz"`
	Altitude   int        `json:"hoehe"`
	Lon        float64    `json:"lon"`
	Lat        float64    `json:"lat"`
	Radiation  float64    `json:"mw"`
	State      State      `json:"status"`
	MgmtNodeID MgmtNodeID `json:"kid"`
}

// Timestamp String representation of the timestamp
type Timestamp string

// Values ...
type Values struct {
	Times                []Timestamp `json:"t"`
	Radiation            []float64   `json:"mw"`
	RadiationCosmic      []float64   `json:"cos"`
	RadiationTerrestrial []float64   `json:"ter"`
	ValueState           []int       `json:"ps"`
	RainChangeTimes      []Timestamp `json:"tr"`
	RainChance           []float64   `json:"r"`
}

// Station StationInfo and 1/24 hour mean-data
type Station struct {
	Info  StationInfo `json:"stamm"`
	MW1h  Values      `json:"mw1h"`
	MW24h Values      `json:"mw24h"`
}

// Stations Map of all stations
type Stations map[string]StationInfo

// Info ...
type Info struct {
	username string
	password string
	client   *http.Client
}

func (info *Info) requestData(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(info.username, info.password)
	return info.client.Do(req)
}

// NewInfo Create new Info for accessing radiation Information
func NewInfo(username, password string) *Info {
	return &Info{
		username: username,
		password: password,
		client:   http.DefaultClient}
}

// GetStation Request detailed information about a station
func (info *Info) GetStation(id string) (*Station, error) {
	url := fmt.Sprintf("%s/%sct.json", baseURL, id)
	resp, err := info.requestData(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(http.StatusText(resp.StatusCode))
		return nil, err
	}

	s := &Station{}
	json.NewDecoder(resp.Body).Decode(s)
	return s, nil
}

// ListStations Lists all stations with their basedata
func (info *Info) ListStations() (*Stations, error) {
	url := fmt.Sprintf("%s/stamm.json", baseURL)
	resp, err := info.requestData(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(http.StatusText(resp.StatusCode))
		return nil, err
	}

	stations := &Stations{}
	json.NewDecoder(resp.Body).Decode(&stations)
	return stations, nil
}
