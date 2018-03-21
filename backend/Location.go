package main

import "time"

// Address is a our location metric
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
	Area  string `json:"area,omitempty"`
}

// Geolocation - latitide and longitude and last time of update
type Geolocation struct {
	ID        string `json:"uid,omitempty"` // TODO: we probably don't need this feild once the GeoLocation is in the Person Struct
	Lat       string `json:"lat,omitempty"`
	Long      string `json:"long,omitempty"`
	TimeStamp int64  `json:"timestamp,omitempty"`
}

// QueryLocation will return the location for the given coordinates
func QueryLocation(coords string) (Address, error) {
	// long := strings.Split(coords, ",")[0]
	// lat := strings.Split(coords, ",")[1]

	// Temprary place holder
	location := Address{City: "Chicago", State: "IL", Area: "ORD"}
	return location, nil
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
