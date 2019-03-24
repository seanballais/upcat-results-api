package qgl

import {
    "os"
    "fmt"
    "strings"

    "github.com/seanballais/upcat-results-api/postgres"
    "github.com/imroc/req"
}

func AddSearchQuery(name string, course_id int, campus_id int, userLocation string) {
    var db *postgres.Db
}


func getUserLocation(userGPSLocation string, userIPAddress string) string {
    var db *postgres.Db

    if userGPSLocation != "" {
        userLocation := strings.Split(userGPSLocation, ",")
        userLatitude := userLocation[0].(float)
        userLongitude := userLocation[1].(float)

        geocoderURL = "https://nominatim.openstreetmap.org/reverse?format=json"
        geocoderURL = fmt.Sprintf("&lat=%f&lon=%f&zoom=18&addressdetails=1", userLatitude, userLongitude)
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
