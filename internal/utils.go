package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"redbudway-api/models"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Printf("%s", err)
	}
	return err == nil
}

func FilterSanitizeNumber(number string) string {
	regexp, _ := regexp.Compile(`[^\d-]`)
	return regexp.ReplaceAllString(number, "")
}

func GenerateUUID() (uuid.UUID, error) {
	u2, err := uuid.NewV4()
	return u2, err
}

func saveImage(path, image string, index int) (string, error) {
	data := strings.Split(image, ",")

	dec, err := base64.StdEncoding.DecodeString(data[1])
	if err != nil {
		return "", err
	}
	format := ""
	switch data[0] {
	case "data:image/jpeg;base64":
		format = ".jpeg"
	case "data:image/png;base64":
		format = ".png"
	case "data:image/webp;base64":
		format = ".webp"
	}

	fileName := fmt.Sprintf("%s/images_%d%s", path, index, format)
	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		return "", err
	}
	if err := f.Sync(); err != nil {
		return "", err
	}

	return fileName, nil
}

func ProcessQuoteImages(customerID, quoteID string, incImages []string) ([]string, error) {
	images := []string{}
	customerPath := fmt.Sprintf("%s/%s", "images", customerID)
	if _, err := os.Stat(customerPath); os.IsNotExist(err) {
		err = os.MkdirAll(customerPath, 0755)
		if err != nil {
			return images, err
		}
	}
	quotePath := fmt.Sprintf("%s/%s", customerPath, customerID)
	if _, err := os.Stat(customerPath); os.IsNotExist(err) {
		err = os.MkdirAll(customerPath, 0755)
		if err != nil {
			return images, err
		}
	}

	if len(incImages) != 0 {
		for index, binary := range incImages {
			if !strings.Contains(binary, "http") {
				URL, err := saveImage(quotePath, binary, index)
				if err != nil {
					log.Printf("Failed to save images, %s", err)
					continue
				}
				images = append(images, URL)
			} else {
				images = append(images, binary)
			}
		}
	}

	return images, nil
}

func ProcessEmailImages(customerID string, incImages []string) ([]string, error) {
	images := []string{}
	emailPath := fmt.Sprintf("%s/%s", "images", "emails")
	if _, err := os.Stat(emailPath); os.IsNotExist(err) {
		err = os.MkdirAll(emailPath, 0755)
		if err != nil {
			return images, err
		}
	}
	customerPath := fmt.Sprintf("%s/%s", emailPath, customerID)
	if _, err := os.Stat(customerPath); os.IsNotExist(err) {
		err = os.MkdirAll(customerPath, 0755)
		if err != nil {
			return images, err
		}
	}

	if len(incImages) != 0 {
		for index, binary := range incImages {
			if !strings.Contains(binary, "http") {
				URL, err := saveImage(customerPath, binary, index)
				if err != nil {
					log.Printf("Failed to save images, %s", err)
					continue
				}
				images = append(images, URL)
			} else {
				images = append(images, binary)
			}
		}
	}

	return images, nil
}

func ProcessImages(tradespersonID, serviceID string, service *models.ServiceDetails) ([]*string, error) {
	images := []*string{}
	tradespersonPath := fmt.Sprintf("%s/%s", "images", tradespersonID)
	if _, err := os.Stat(tradespersonPath); os.IsNotExist(err) {
		err = os.MkdirAll(tradespersonPath, 0755)
		if err != nil {
			return images, err
		}
	}
	servicePath := fmt.Sprintf("%s/%s", tradespersonPath, serviceID)
	if _, err := os.Stat(servicePath); os.IsNotExist(err) {
		err = os.MkdirAll(servicePath, 0755)
		if err != nil {
			return images, err
		}
	}

	if len(service.Images) != 0 {
		for index, binary := range service.Images {
			if !strings.Contains(binary, "http") {
				fileName, err := saveImage(servicePath, binary, index)
				URL := fmt.Sprintf("%s/%s", "https://"+os.Getenv("SUBDOMAIN")+"redbudway.com", fileName)
				if err != nil {
					log.Printf("Failed to save images, %s", err)
					continue
				}
				images = append(images, &URL)
			} else {
				images = append(images, &binary)
			}
		}
	} else {
		URL := "https://" + os.Getenv("SUBDOMAIN") + "redbudway.com/assets/images/deal.svg"
		images = append(images, &URL)
	}

	return images, nil
}

func GenerateQuoteSuffix() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 24)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

type City struct {
	Name string
}

func SelectedCities(citiesJson, city string) (bool, error) {
	cities := []City{}
	err := json.Unmarshal([]byte(citiesJson), &cities)
	if err != nil {
		return false, err
	}

	for _, selectedCity := range cities {
		if selectedCity.Name == city {
			return true, nil
		}
	}

	return false, nil
}

func CreateTimeAndPrice(startTime, endTime string, decimalPrice float64) (string, error) {
	startDate, err := GetStartDate(startTime)
	if err != nil {
		log.Printf("Failed to get start date, %v", err)
		return "", err
	}
	endDate, err := GetEndDate(endTime)
	if err != nil {
		log.Printf("Failed to get end date, %v", err)
		return "", err
	}
	date := startDate.Format("Monday, January 2, 2006")
	startTime = startDate.Format("3:04")
	endTime = endDate.Format("3:04PM")
	return fmt.Sprintf("%s<br>%s - %s<br>$%.2f<br><br>", date, startTime, endTime, decimalPrice), nil
}

func CreateEndTime(startTime, segmentSize string) (string, error) {
	endTime, err := time.Parse("2006-1-2 15:04:00", startTime)
	if err != nil {
		log.Printf("Failed to parse endTime")
		return "", err
	}
	segment, err := strconv.Atoi(segmentSize)
	if err != nil {
		log.Printf("Failed to cast string to int")
		return "", err
	}
	milliSegment := segment * 60 * 60 * 1000
	endTime = endTime.Add(time.Millisecond * time.Duration(milliSegment))
	return endTime.Format("2006-1-2 15:04:00"), nil
}

func GetEndDate(endTime string) (time.Time, error) {
	endDate, err := time.Parse("1/2/2006, 3:04:00 PM", endTime)
	if err != nil {
		log.Printf("Failed to parse endTime %s", endTime)
		return endDate, err
	}
	endDate = endDate.AddDate(0, 0, 1)
	return endDate, nil
}

func GetStartDate(startTime string) (time.Time, error) {
	startDate, err := time.Parse("1/2/2006, 3:04:00 PM", startTime)
	if err != nil {
		log.Printf("Failed to parse startTime %s", startTime)
		return startDate, err
	}

	return startDate, nil
}

func CreateTimeAndPriceFrmDB(startTime, endTime string, decimalPrice float64) (string, error) {
	startDate, err := GetStartDateFrmDB(startTime)
	if err != nil {
		log.Printf("Failed to get start date, %v", err)
		return "", err
	}
	endDate, err := GetEndDateFrmDB(endTime)
	if err != nil {
		log.Printf("Failed to get end date, %v", err)
		return "", err
	}
	date := startDate.Format("Monday, January 2, 2006")
	startTime = startDate.Format("3:04")
	endTime = endDate.Format("3:04PM")
	return fmt.Sprintf("%s<br>%s - %s<br>$%.2f<br><br>", date, startTime, endTime, decimalPrice), nil
}

func GetEndDateFrmDB(endTime string) (time.Time, error) {
	endDate, err := time.Parse("2006-1-2 15:04:00", endTime)
	if err != nil {
		log.Printf("Failed to parse endTime %s", endTime)
		return endDate, err
	}
	endDate.AddDate(0, 0, 1)
	return endDate, nil
}

func GetStartDateFrmDB(startTime string) (time.Time, error) {
	startDate, err := time.Parse("2006-1-2 15:04:00", startTime)
	if err != nil {
		log.Printf("Failed to parse startTime %s", startTime)
		return startDate, err
	}

	return startDate, nil
}

func CreateSubscriptionTimeAndPrice(interval, startTime, endTime string, decimalPrice float64) (string, error) {
	startDate, err := GetStartDate(startTime)
	if err != nil {
		log.Printf("Failed to get start date, %v", err)
		return "", err
	}

	var day string
	if interval == "week" {
		day = "Every " + startDate.Format("Monday")
	} else if interval == "month" {
		day = "Every " + startDate.Format("2") + " of the month"
	} else if interval == "year" {
		day = "Every " + startDate.Format("January 2")
	}

	endDate, err := GetEndDate(endTime)
	if err != nil {
		log.Printf("Failed to get end date, %v", err)
		return "", err
	}
	startTime = startDate.Format("3:04")
	endTime = endDate.Format("3:04PM")
	return fmt.Sprintf("%s<br>%s - %s<br>$%.2f<br><br>", day, startTime, endTime, decimalPrice), nil
}

func CreateSubscriptionTimeAndPriceFrmDB(interval, startTime, endTime string, decimalPrice float64) (string, error) {
	startDate, err := GetStartDateFrmDB(startTime)
	if err != nil {
		log.Printf("Failed to get start date, %v", err)
		return "", err
	}

	var day string
	if interval == "week" {
		day = "Every " + startDate.Format("Monday")
	} else if interval == "month" {
		day = "Every " + startDate.Format("2") + " of the month"
	} else if interval == "year" {
		day = "Every " + startDate.Format("January 2")
	}

	endDate, err := GetEndDateFrmDB(endTime)
	if err != nil {
		log.Printf("Failed to get end date, %v", err)
		return "", err
	}
	startTime = startDate.Format("3:04")
	endTime = endDate.Format("3:04PM")
	return fmt.Sprintf("%s<br>%s - %s<br>$%.2f<br><br>", day, startTime, endTime, decimalPrice), nil
}

func GetDueDate(date string) (time.Time, error) {
	dueDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Printf("Failed to parse dueDate %s", date)
		return dueDate, err
	}
	return dueDate, nil
}
