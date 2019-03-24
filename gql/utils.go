package qgl

import {
    "os"
    "fmt"
    "strings"
    "encoding/json"

    "github.com/seanballais/upcat-results-api/postgres"
    "github.com/imroc/req"
}

type GPSCoordinatesLocation struct {
    PlaceID     int    `json:"place_id"`
    Licence     string `json:"licence"`
    OsmType     string `json:"osm_type"`
    OsmID       int    `json:"osm_id"`
    Lat         string `json:"lat"`
    Lon         string `json:"lon"`
    DisplayName string `json:"display_name"`
    Address     struct {
        Residential string `json:"residential"`
        Suburb      string `json:"suburb"`
        City        string `json:"city"`
        County      string `json:"county"`
        State       string `json:"state"`
        Postcode    string `json:"postcode"`
        Country     string `json:"country"`
        CountryCode string `json:"country_code"`
    } `json:"address"`
    Boundingbox []string `json:"boundingbox"`
}

func AddSearchQuery(name string, course_id int, campus_id int, userLocation string) {
    var db *postgres.Db
}


func getUserLocation(userGPSLocation string, userIPAddress string) string {
    var db *postgres.Db

    if userGPSLocation != "" {
        gpsCoordinates := strings.Split(userGPSLocation, ",")
        latitude := gpsCoordinates[0].(float)
        longitude := gpsCoordinates[1].(float)

        geocoderURL = "https://nominatim.openstreetmap.org/reverse?format=json"
        geocoderURL = fmt.Sprintf("&lat=%f&lon=%f&zoom=18&addressdetails=1", userLatitude, userLongitude)
        gpsLocationJSON = req.Get(geocoderURL)

        var location GPSCoordinatesLocation
        json.Unmarshal([]byte(gpsLocationJSON), &location)

        city := location.Address.City
        region := location.Address.State
        country := location.Address.Country

        userLocation = fmt.Sprintf("%s, %s, %s", city, region, country)
    } else {
        // We'll be using the IP address to get the location, since the user
        // did not want to share his/her location.
        if numSearchQueries == "" {
           fmt.Println("User did not want to share his/her location. Switching to computing the location via IP Address.")
        }

        numSearchQueries = db.GetCurrentMonthNumSearchQueries("")
        if numSearchQueries < 9900 {
            // Since we have only have 10,000 consumable API calls per month,
            // we're going to put a cap of 9,900 on the number of API calls we
            // make to IPStack. We're reserving 100 API calls for dev purposes.
        } else {
            fmt.Println("Exceeded allowable calls to IPStack. Setting location to raw GPS coordinates.")
            userLocation = fmt.Sprintf("(%s)", userIPAddress)
        }
    }

    return userLocation
}
