package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"redbudway-api/models"
	"regexp"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/h2non/bimg"
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

func SaveProfileImage(tradespersonID, image string) (string, error) {
	url := ""
	data := strings.Split(image, ",")

	dec, err := base64.StdEncoding.DecodeString(data[1])
	if err != nil {
		log.Println("Failed to decode")
		return url, err
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

	path := fmt.Sprintf("%s/%s", "images", tradespersonID)
	//add to Util package
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return url, err
		}
	}

	fileName := fmt.Sprintf("%s/%s%s", path, tradespersonID, format)
	f, err := os.Create(fileName)
	if err != nil {
		log.Println("Failed to create file with name %s", fileName)
		return url, err
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		return url, err
	}
	if err := f.Sync(); err != nil {
		return url, err
	}

	return fmt.Sprintf("https://"+os.Getenv("SUBDOMAIN")+"redbudway.com/%s", fileName), nil
}

func SaveImage(tradespersonID, image, imageType string) (string, error) {
	url := ""
	data := strings.Split(image, ",")

	dec, err := base64.StdEncoding.DecodeString(data[1])
	if err != nil {
		log.Println("Failed to decode")
		return url, err
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

	path := fmt.Sprintf("%s/%s", "images", tradespersonID)
	//add to Util package
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return url, err
		}
	}

	fileName := fmt.Sprintf("%s/%s%s", path, imageType, format)
	f, err := os.Create(fileName)
	if err != nil {
		log.Println("Failed to create file with name %s", fileName)
		return url, err
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		return url, err
	}
	if err := f.Sync(); err != nil {
		return url, err
	}

	return fmt.Sprintf("https://"+os.Getenv("SUBDOMAIN")+"redbudway.com/%s", fileName), nil
}

func contains(images []string, fileName string) bool {
	exist := false
	for _, image := range images {
		if strings.Contains(image, fileName) {
			exist = true
		}
	}
	return exist
}

func removeImages(servicePath string, images []string) error {
	if _, err := os.Stat(servicePath); !os.IsNotExist(err) {
		files, err := ioutil.ReadDir(servicePath)
		if err != nil {
			log.Printf("Failed to read director %s, %v", servicePath, err)
			return err
		}

		for _, file := range files {
			if !file.IsDir() {
				if !contains(images, file.Name()) {
					filePath := fmt.Sprintf("%s/%s", servicePath, file.Name())
					err := os.Remove(filePath)
					if err != nil {
						log.Printf("Failed to remove file %s, %v", filePath, err)
					}
				}
			}
		}
	}
	return nil
}

func saveImage(path, imageBytes string, index int) (string, error) {
	fileName := fmt.Sprintf("%s/image_%d%s", path, index, ".webp")

	data := strings.Split(imageBytes, ",")

	dec, err := base64.StdEncoding.DecodeString(data[1])
	if err != nil {
		return "", err
	}

	converted, err := bimg.NewImage(dec).Convert(bimg.WEBP)
	if err != nil {
		return fileName, err
	}

	processed, err := bimg.NewImage(converted).Process(bimg.Options{Quality: 70})
	if err != nil {
		return fileName, err
	}

	writeError := bimg.Write(fileName, processed)
	if writeError != nil {
		return fileName, writeError
	}

	return fileName, nil
}

func ProcessEmailImages(customerID, quoteID string, incImages []string) ([]string, error) {
	images := []string{}
	emailPath := fmt.Sprintf("%s/%s", "images", "emails")
	if _, err := os.Stat(emailPath); os.IsNotExist(err) {
		err = os.MkdirAll(emailPath, 0755)
		if err != nil {
			return images, err
		}
	}
	quotePath := fmt.Sprintf("%s/%s", emailPath, customerID)
	if quoteID != "" {
		quotePath = fmt.Sprintf("%s/%s", quotePath, quoteID)
	}
	if _, err := os.Stat(quotePath); os.IsNotExist(err) {
		err = os.MkdirAll(quotePath, 0755)
		if err != nil {
			return images, err
		}
	}

	if len(incImages) != 0 {
		for index, binary := range incImages {
			URL, err := saveImage(quotePath, binary, index)
			if err != nil {
				log.Printf("Failed to save images, %s", err)
				continue
			}
			images = append(images, URL)
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

	removeImages(servicePath, service.Images)

	if len(service.Images) != 0 {
		for i := range service.Images {
			if !strings.Contains(service.Images[i], "https://") {
				fileName, err := saveImage(servicePath, service.Images[i], i)
				URL := fmt.Sprintf("%s/%s", "https://"+os.Getenv("SUBDOMAIN")+"redbudway.com", fileName)
				if err != nil {
					log.Printf("Failed to save images, %s", err)
					continue
				}
				images = append(images, &URL)
			} else {
				images = append(images, &service.Images[i])
			}
		}
	}

	return images, nil
}

func GetImage(ID, tradespersonID string) (string, error) {
	url := ""

	fileName := fmt.Sprintf("%s/%s/%s/%s", "images", tradespersonID, ID, "image_0.webp")
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		url = "https://" + os.Getenv("SUBDOMAIN") + "redbudway.com/" + fileName
	}
	if url == "" {
		url = "https://" + os.Getenv("SUBDOMAIN") + "redbudway.com/assets/images/placeholder.svg"
	}

	return url, nil
}

func GetImages(ID, tradespersonID string) ([]string, error) {
	images := []string{}

	servicePath := fmt.Sprintf("%s/%s/%s", "images", tradespersonID, ID)
	if _, err := os.Stat(servicePath); !os.IsNotExist(err) {
		files, err := ioutil.ReadDir(servicePath)
		if err != nil {
			log.Printf("Failed to read director %s, %v", servicePath, err)
			return images, err
		}

		for _, file := range files {
			if !file.IsDir() {
				images = append(images, "https://"+os.Getenv("SUBDOMAIN")+"redbudway.com/"+servicePath+"/"+file.Name())
			}
		}
	}

	return images, nil
}

func GetQuoteImages(customerEmail, quoteID string) ([]string, error) {
	images := []string{}

	servicePath := fmt.Sprintf("images/emails/%s/%s", customerEmail, quoteID)
	if _, err := os.Stat(servicePath); !os.IsNotExist(err) {
		files, err := ioutil.ReadDir(servicePath)
		if err != nil {
			log.Printf("Failed to read director %s, %v", servicePath, err)
			return images, err
		}

		for _, file := range files {
			if !file.IsDir() {
				images = append(images, "https://"+os.Getenv("SUBDOMAIN")+"redbudway.com/"+servicePath+"/"+file.Name())
			}
		}
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

func CreateForm(form []models.FormFields) string {
	var tableRows string
	tableRows += "<table border=\"0\">"
	for _, row := range form {
		tableRows += "<tr>"
		for _, col := range row {
			tableRows += fmt.Sprintf("<td style=\"max-width: 200px; word-wrap: break-word; vertical-align: top;\"><b>%s: </b>%s</td>", col.Field, col.Value)
		}
		tableRows += "</tr>"
	}
	tableRows += "</table><br><br>"
	return tableRows
}

func CreateTimeAndPrice(startTime, endTime, timeZone string, decimalPrice float64) (string, error) {
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
	return fmt.Sprintf("%s<br>%s - %s&nbsp;%s<br>$%.2f<br><br>", date, startTime, endTime, timeZone, decimalPrice), nil
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

func CreateTimeAndPriceFrmDB(startTime, endTime, timeZone string, decimalPrice float64) (string, error) {
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
	return fmt.Sprintf("%s<br>%s - %s&nbsp;%s<br>$%.2f<br><br>", date, startTime, endTime, timeZone, decimalPrice), nil
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

func CreateSubscriptionTimeAndPrice(interval, startTime, endTime, timeZone string, decimalPrice float64) (string, error) {
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
	return fmt.Sprintf("%s<br>%s - %s&nbsp;%s<br>$%.2f<br><br>", day, startTime, endTime, timeZone, decimalPrice), nil
}

func GetDueDate(date string) (time.Time, error) {
	dueDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Printf("Failed to parse dueDate %s", date)
		return dueDate, err
	}
	return dueDate, nil
}
