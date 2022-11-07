package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"redbudway-api/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
)

func GetLocationHandler(params operations.GetLocationParams) middleware.Responder {
	latitude := params.Latitude
	longitude := params.Longitude

	response := operations.NewGetLocationOK()

	URL := fmt.Sprintf("http://www.mapquestapi.com/geocoding/v1/reverse?key=mF8fwnXysBtlAHP7Z2ZNxGWr8cdBNIHA&location=%s,%s&outFormat=json", latitude, longitude)
	resp, err := http.Get(URL)
	if err != nil {
		log.Printf("Failed to get users location from mapquest api, %v", err)
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

	results := r["results"].([]interface{})
	result := results[0].(map[string]interface{})
	locations := result["locations"].([]interface{})
	location := locations[0].(map[string]interface{})
	city := location["adminArea5"].(string)
	state := location["adminArea3"].(string)
	payload := operations.GetLocationOKBody{City: city, State: state}
	response.SetPayload(&payload)

	return response
}
