package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"redbudway-api/restapi/operations"

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
		payload := operations.GetLocationOKBody{City: city, State: state}
		response.SetPayload(&payload)
	}

	return response
}
