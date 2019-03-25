package qgl

import {
    "os"
    "fmt"
    "strings"
    "encoding/json"

    "github.com/imroc/req"
    "github.com/harwoeck/ipstack"

    "github.com/seanballais/upcat-results-api/postgres"
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

func getUserLocationID(userGPSLocation string, userIPAddress string) int {
    var db *postgres.Db
    var userLocationID int

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

        userLocation := fmt.Sprintf("%s, %s, %s", city, region, country)
        userLocationID = db.AddLocation(userLocation)
    } else {
        // We'll be using the IP address to get the location, since the user
        // did not want to share his/her location.
        if numSearchQueries == "" {
           fmt.Println("User did not want to share his/her location. Switching to computing the location via IP Address.")
        }

        numIPAddresses = db.GetCurrentMonthMappedIPAddresses()
        if numIPAddresses < 9900 {
            // Since we have only have 10,000 consumable API calls per month,
            // we're going to put a cap of 9,900 on the number of API calls we
            // make to IPStack. We're reserving 100 API calls for dev purposes.
            if db.IsIPAddressCached(userIPAddress) {
                userLocationID = db.GetIPAddressLocationID(userIPAddress)
            } else {
                // Cache the IP address if we haven't already.
                ipStackKey := os.Getenv("UPCAT_RESULTS_API_IPSTACK_API_KEY")
                ipStackClient := ipstack.Client.NewClient(ipStackKey, true, 30)
                ipStackResp, err := ipStackClient.Check(userIPAddress)
                if err != nil {
                    fmt.Println("IP Stack Client Error: ", err)
                } else {
                    city := ipStackResp.City
                    region := ipStackResp.RegionName
                    country := ipStackResp.CountryName
                    
                    userLocation := fmt.Sprintf("%s, %s, %s,", city, region, country)

                    userLocationID = db.AddLocation(userLocation)
                    db.AddIPAddressLocationMapping(userIPAddress, userLocationID)
                }
            }
        } else {
            fmt.Println("Exceeded allowable calls to IPStack. Setting location to raw GPS coordinates.")
            userLocationID = fmt.Sprintf("(%s)", userIPAddress)
        }
    }

    return userLocationID
}
