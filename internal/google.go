package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"redbudway-api/models"
	"time"
)

func GetNextPage(accessToken, pageToken string) (models.GoogleTimeSlots, error) {
	events := models.GoogleTimeSlots{}

	loop := true
	for loop {
		client := &http.Client{}
		URL := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/primary/events?singleEvents=true&showDeleted=%v&timeMin=%s&timeMax=%s&maxResults=10&pageToken=%s", false, time.Now().Format(time.RFC3339), time.Now().AddDate(1, 0, 0).Format(time.RFC3339), pageToken)
		req, err := http.NewRequest(http.MethodGet, URL, nil)
		if err != nil {
			log.Printf("Failed to create new request, %v", err)
			return events, err
		}
		req.Header.Add("content-type", "application/json")
		req.Header.Add("Authorization", "Bearer "+accessToken)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to get user calendar events, %v", err)
			return events, err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read body in %v", err)
			return events, err
		}
		var res map[string]interface{}

		err = json.Unmarshal(body, &res)
		if err != nil {
			log.Printf("Failed to unmarshal response, %v", err)
			return events, err
		}

		if res["error"] != nil {
			log.Printf("Error in retrieveing google events, %v", res)
			return events, errors.New(res["error"].(string))
		}

		items := res["items"].([]interface{})
		for _, item := range items {
			event := models.GoogleTimeSlotsItems0{}
			i := item.(map[string]interface{})
			s := i["start"].(map[string]interface{})
			if s["dateTime"] != nil {
				event.StartTime = s["dateTime"].(string)
				e := i["end"].(map[string]interface{})
				event.EndTime = e["dateTime"].(string)
				event.TimeZone = e["timeZone"].(string)
				if i["recurrence"] != nil {
					r := i["recurrence"].([]interface{})
					event.Recurrence = r[0].(string)
				}
				events = append(events, &event)
			}
		}

		if res["nextPageToken"] != nil {
			pageToken = res["nextPageToken"].(string)
		} else {
			loop = false
		}
	}

	return events, nil
}

func GetGoogleTimeSlots(accessToken string) models.GoogleTimeSlots {
	events := models.GoogleTimeSlots{}

	client := &http.Client{}
	URL := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/primary/events?singleEvents=true&showDeleted=%v&timeMin=%s&timeMax=%s&maxResults=10", false, time.Now().Format(time.RFC3339), time.Now().AddDate(1, 0, 0).Format(time.RFC3339))
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		log.Printf("Failed to create new request, %v", err)
		return events
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to get user calendar events, %v", err)
		return events
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body in %v", err)
	}
	var res map[string]interface{}

	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("Failed to unmarshal response, %v", err)
		return events
	}

	if res["error"] != nil {
		log.Printf("Error in retrieveing google events, %v", res)
		return events
	}

	pagedEvents, err := GetNextPage(accessToken, res["nextPageToken"].(string))
	if err != nil {
		log.Printf("Failed to page through google events, %s", err)
		return events
	}

	items := res["items"].([]interface{})
	for _, item := range items {
		event := models.GoogleTimeSlotsItems0{}
		i := item.(map[string]interface{})
		s := i["start"].(map[string]interface{})
		if s["dateTime"] != nil {
			event.StartTime = s["dateTime"].(string)
			e := i["end"].(map[string]interface{})
			event.EndTime = e["dateTime"].(string)
			event.TimeZone = e["timeZone"].(string)
			if i["recurrence"] != nil {
				r := i["recurrence"].([]interface{})
				event.Recurrence = r[0].(string)
			}
			events = append(events, &event)
		}
	}
	events = append(events, pagedEvents...)

	return events
}
