package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"redbudway-api/restapi/operations"
	"strings"

	"github.com/go-openapi/runtime/middleware"
)

func GetLocationHandler(params operations.GetLocationParams) middleware.Responder {
	latitude := params.Latitude
	longitude := params.Longitude

	payload := operations.GetLocationOKBody{City: "", State: ""}
	response := operations.NewGetLocationOK().WithPayload(&payload)
	URL := fmt.Sprintf("https://geocode-api.arcgis.com/arcgis/rest/services/World/GeocodeServer/reverseGeocode?location=%s,%s&f=json&token=%s", longitude, latitude, os.Getenv("ARCGIS_TOKEN"))
	resp, err := http.Get(URL)
	if err != nil {
		log.Printf("Failed to get users location from arcgis api, %v", err)
		return response
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body in %v", err)
	}
	var r map[string]interface{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Printf("Failed to unmarshal response, %v", err)
		return response
	}

	address := r["address"].(map[string]interface{})
	if address["CountryCode"].(string) == "USA" {
		city := address["City"].(string)
		state := address["Region"].(string)
		timeZone := getTimeZone(latitude, longitude)
		payload := operations.GetLocationOKBody{City: city, State: state, TimeZone: timeZone}
		response.SetPayload(&payload)
	}

	return response
}

func GetAddressHandler(params operations.GetAddressParams) middleware.Responder {
	address := strings.ReplaceAll(params.Address, " ", "")

	payload := operations.GetAddressOKBody{City: "", State: ""}
	response := operations.NewGetAddressOK().WithPayload(&payload)
	URL := fmt.Sprintf("https://geocode.arcgis.com/arcgis/rest/services/World/GeocodeServer/findAddressCandidates?SingleLine=%s&countryCode=US&f=json&outFields=city,region&token=%s", address, os.Getenv("ARCGIS_TOKEN"))
	resp, err := http.Get(URL)
	if err != nil {
		log.Printf("Failed to get users location from arcgis api, %v", err)
		return response
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body in %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return response
	}
	var r map[string]interface{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Printf("Failed to unmarshal response, %v", err)
		return response
	}

	candidates := r["candidates"].([]interface{})
	if len(candidates) >= 1 {
		candidate := candidates[0].(map[string]interface{})
		attributes := candidate["attributes"].(map[string]interface{})
		city := attributes["city"].(string)
		state := attributes["region"].(string)
		location := candidate["location"].(map[string]interface{})
		longitude := fmt.Sprintf("%f", location["x"].(float64))
		latitude := fmt.Sprintf("%f", location["y"].(float64))
		timeZone := getTimeZone(latitude, longitude)
		payload := operations.GetAddressOKBody{City: city, State: state, TimeZone: timeZone}
		response.SetPayload(&payload)
	}

	return response
}

func getTimeZone(latitude, longitude string) string {
	timeZone := ""
	URL := fmt.Sprintf("https://api.wheretheiss.at/v1/coordinates/%s,%s", latitude, longitude)
	resp, err := http.Get(URL)
	if err != nil {
		log.Printf("Failed to get users timezone from wheretheiss api, %v", err)
		return timeZone
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body in %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return timeZone
	}
	var r map[string]interface{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Printf("Failed to unmarshal response, %v", err)
		return timeZone
	}

	timeZone = r["timezone_id"].(string)

	return timeZone
}
