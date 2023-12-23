package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func VerifyReCaptcha(token string) (bool, error) {
	valid := false
	URL := fmt.Sprintf("https://www.google.com/recaptcha/api/siteverify?secret=%s&response=%s", "6Lfi4wopAAAAAEYrv06awJFtSL2NP1vxxCuJYKjC", token)
	resp, err := http.Post(URL, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		log.Printf("Failed to recaptcha token, %v", err)
		return valid, err
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
		return valid, err
	}
	valid = r["success"].(bool)

	return valid, nil
}
