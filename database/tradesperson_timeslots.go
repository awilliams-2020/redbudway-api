package database

import (
	"database/sql"
	"encoding/json"
	"log"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"
	"time"
)

func InsertTimeSlots(fixedPriceID int64, fixedPrice *models.ServiceDetails) error {
	for _, timeSlot := range fixedPrice.TimeSlots {
		t, err := time.Parse("1/2/2006, 3:04:00 PM", timeSlot.StartTime)
		if err != nil {
			return err
		}
		formattedDate := t.Format("2006-01-02 15:04:00")
		stmt, err := db.Prepare("INSERT INTO fixed_price_time_slots (fixedPriceId, startTime, segmentSize, maxPeople) VALUES (?, ?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(fixedPriceID, formattedDate, timeSlot.SegmentSize, timeSlot.MaxPeople)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateTimeSlots(fixedPriceID int64, fixedPrice *models.ServiceDetails) error {
	stmt, err := db.Prepare("SELECT startTime, segmentSize FROM fixed_price_time_slots WHERE fixedPriceId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(fixedPriceID)
	if err != nil {
		return err
	}

	existingTimeSlots := []*models.TimeSlotsItems0{}
	var startTime, segmentSize string
	for rows.Next() {
		if err := rows.Scan(&startTime, &segmentSize); err != nil {
			return err
		}
		existingTimeSlot := &models.TimeSlotsItems0{}
		existingTimeSlot.StartTime = startTime
		existingTimeSlot.SegmentSize = segmentSize
		existingTimeSlots = append(existingTimeSlots, existingTimeSlot)
	}

	for _, existingTimeSlot := range existingTimeSlots {
		found := false
		for _, timeSlot := range fixedPrice.TimeSlots {
			t, err := time.Parse("1/2/2006, 3:04:00 PM", timeSlot.StartTime)
			if err != nil {
				return err
			}
			formattedDate := t.Format("2006-01-02 15:04:00")
			if formattedDate == existingTimeSlot.StartTime {
				found = true
				if timeSlot.SegmentSize != existingTimeSlot.SegmentSize {
					stmt, err := db.Prepare("UPDATE fixed_price_time_slots SET segmentSize=? AND maxPeople=? WHERE fixedPriceId=? AND startTime=?")
					if err != nil {
						return err
					}
					defer stmt.Close()
					_, err = stmt.Exec(timeSlot.SegmentSize, timeSlot.MaxPeople, fixedPriceID, formattedDate)
					if err != nil {
						return err
					}
				}
			}
		}
		if !found {
			stmt, err := db.Prepare("DELETE FROM fixed_price_time_slots WHERE fixedPriceId=? AND startTime=?")
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(fixedPriceID, existingTimeSlot.StartTime)
			if err != nil {
				return err
			}
		}
	}

	for _, timeSlot := range fixedPrice.TimeSlots {
		t, err := time.Parse("1/2/2006, 3:04:00 PM", timeSlot.StartTime)
		if err != nil {
			return err
		}
		formattedDate := t.Format("2006-01-02 15:04:00")
		found := false
		for _, existingTimeSlot := range existingTimeSlots {
			if formattedDate == existingTimeSlot.StartTime {
				found = true
			}
		}
		if !found {
			stmt, err := db.Prepare("INSERT INTO fixed_price_time_slots (fixedPriceId, startTime, segmentSize, maxPeople) VALUES (?, ?, ?, ?)")
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(fixedPriceID, formattedDate, timeSlot.SegmentSize, timeSlot.MaxPeople)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ResetTakenTimeSlotByCustomer(cuStripeID string) error {
	stmt, err := db.Prepare("UPDATE fixed_price_time_slots SET taken=False, takenBy=NULL, cuStripeId=NULL WHERE cuStripeId=? and startTime < CURDATE()")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cuStripeID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTimeSlotByInvoice(invoiceID string) error {

	stmt, err := db.Prepare("UPDATE fixed_price_time_slots fpts INNER JOIN customer_time_slots cts ON fpts.id=cts.timeSlotId SET fpts.curPeople=fpts.curPeople-1 WHERE cts.invoiceId=? AND fpts.startTime > CURDATE()")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(invoiceID)
	if err != nil {
		return err
	}

	stmt, err = db.Prepare("DELETE FROM customer_time_slots WHERE invoiceId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(invoiceID)
	if err != nil {
		return err
	}

	return nil
}

func GetTakenTimeSlot(startTime string, fixedPriceID int64) ([]string, []string, error) {

	cuStripeIDs := []string{}
	takenBy := []string{}

	stmt, err := db.Prepare("SELECT takenBy, cuStripeIds, maxPeople, curPeople FROM fixed_price_time_slots WHERE startTime=? AND fixedPriceId=?")
	if err != nil {
		return cuStripeIDs, takenBy, err
	}
	defer stmt.Close()

	var t, c sql.NullString
	var maxPeople, curPeople int64
	err = stmt.QueryRow(startTime, fixedPriceID).Scan(&t, &c, &maxPeople, &curPeople)
	if err != nil {
		return cuStripeIDs, takenBy, err
	}

	if t.Valid {
		var t2 interface{}
		if err := json.Unmarshal([]byte(t.String), &t2); err != nil {
			return cuStripeIDs, takenBy, err
		}
		takenBy = t2.([]string)
	}

	if c.Valid {
		var c2 interface{}
		if err := json.Unmarshal([]byte(c.String), &c2); err != nil {
			return cuStripeIDs, takenBy, err
		}
		cuStripeIDs = c2.([]string)
	}

	return cuStripeIDs, takenBy, nil
}

func updateCustomerInvoiceTimeSlot(timeSlotID int64, stripeInvoiceID, cuStripeID string) error {

	stmt, err := db.Prepare("INSERT INTO customer_time_slots (timeSlotId, invoiceId, cuStripeId) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(timeSlotID, stripeInvoiceID, cuStripeID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTakenTimeSlot(stripeInvoiceID, cuStripeID, startTime string, fixedPriceID int64, timeSlotID int64) (bool, error) {
	startDate, err := time.Parse("1/2/2006, 3:04:00 PM", startTime)
	if err != nil {
		log.Printf("Failed to parse startTime %s", startTime)
		return false, err
	}

	startTime = startDate.Format("2006-1-2 15:04:00")

	if err := updateCustomerInvoiceTimeSlot(timeSlotID, stripeInvoiceID, cuStripeID); err != nil {
		log.Printf("Failed to insert customer time slot, %v", err)
		return false, err
	}

	stmt, err := db.Prepare("UPDATE fixed_price_time_slots SET curPeople=curPeople+1 WHERE startTime=? AND fixedPriceId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(startTime, fixedPriceID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func updateCustomerSubscriptionTimeSlot(timeSlotID int64, subscriptionID, cuStripeID string) error {

	stmt, err := db.Prepare("INSERT INTO customer_time_slots (timeSlotId, subscriptionId, cuStripeId) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(timeSlotID, subscriptionID, cuStripeID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateWeeklyTimeSlot(subscriptionID, cuStripeID, startTime string, fixedPriceID, timeSlotID int64) (bool, error) {
	startDate, err := time.Parse("1/2/2006, 3:04:00 PM", startTime)
	if err != nil {
		log.Printf("Failed to parse startTime %s", startTime)
		return false, err
	}
	startTime = startDate.Format("2006-1-2 15:04:00")

	if err := updateCustomerSubscriptionTimeSlot(timeSlotID, subscriptionID, cuStripeID); err != nil {
		log.Printf("Failed to insert customer time slot, %v", err)
		return false, err
	}

	stmt, err := db.Prepare("UPDATE fixed_price_time_slots SET curPeople=curPeople+1 WHERE DAYNAME(startTime)=DAYNAME(?) AND TIME(startTime)=TIME(?) AND fixedPriceId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(startTime, startTime, fixedPriceID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func UpdateMonthlyTimeSlot(subscriptionID, cuStripeID, startTime string, fixedPriceID, timeSlotID int64) (bool, error) {
	startDate, err := time.Parse("1/2/2006, 3:04:00 PM", startTime)
	if err != nil {
		log.Printf("Failed to parse startTime %s", startTime)
		return false, err
	}
	startTime = startDate.Format("2006-1-2 15:04:00")

	if err := updateCustomerSubscriptionTimeSlot(timeSlotID, subscriptionID, cuStripeID); err != nil {
		log.Printf("Failed to insert customer time slot, %v", err)
		return false, err
	}

	stmt, err := db.Prepare("UPDATE fixed_price_time_slots SET curPeople=curPeople+1 WHERE DAYOFMONTH(startTime)=DAYOFMONTH(?) AND TIME(startTime)=TIME(?) AND fixedPriceId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(startTime, startTime, fixedPriceID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func UpdateYearlyTimeSlot(subscriptionID, cuStripeID, startTime string, fixedPriceID, timeSlotID int64) (bool, error) {
	startDate, err := time.Parse("1/2/2006, 3:04:00 PM", startTime)
	if err != nil {
		log.Printf("Failed to parse startTime %s", startTime)
		return false, err
	}
	startTime = startDate.Format("2006-1-2 15:04:00")

	if err := updateCustomerSubscriptionTimeSlot(timeSlotID, subscriptionID, cuStripeID); err != nil {
		log.Printf("Failed to insert customer time slot, %v", err)
		return false, err
	}

	stmt, err := db.Prepare("UPDATE fixed_price_time_slots SET curPeople=curPeople+1 WHERE MONTH(startTime)=MONTH(?) AND DAYOFMONTH(startTime)=DAYOFMONTH(?) AND TIME(startTime)=TIME(?) AND fixedPriceId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(startTime, startTime, startTime, fixedPriceID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func getServiceTimeSlots(fixedPriceId int) ([]*operations.GetTradespersonTradespersonIDTimeSlotsOKBodyItems0TimeSlotsItems0, error) {
	timeSlots := []*operations.GetTradespersonTradespersonIDTimeSlotsOKBodyItems0TimeSlotsItems0{}

	stmt, err := db.Prepare("SELECT startTime, segmentSize FROM fixed_price_time_slots WHERE fixedPriceId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return timeSlots, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(fixedPriceId)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return timeSlots, err
	}

	var startTime, segmentSize string
	for rows.Next() {
		if err := rows.Scan(&startTime, &segmentSize); err != nil {
			return timeSlots, err
		}
		timeSlot := &operations.GetTradespersonTradespersonIDTimeSlotsOKBodyItems0TimeSlotsItems0{}
		timeSlot.StartTime = startTime
		timeSlot.SegmentSize = segmentSize
		timeSlots = append(timeSlots, timeSlot)
	}

	return timeSlots, nil
}

func GetTradespersonTimeslots(tradespersonID string) (*operations.GetTradespersonTradespersonIDTimeSlotsOK, error) {
	db := GetConnection()

	response := operations.NewGetTradespersonTradespersonIDTimeSlotsOK()
	services := []*operations.GetTradespersonTradespersonIDTimeSlotsOKBodyItems0{}
	response.SetPayload(services)

	stmt, err := db.Prepare("SELECT id, subscription, subInterval FROM fixed_prices WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return response, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return response, err
	}

	var id int
	var subscription bool
	var interval string
	for rows.Next() {
		if err := rows.Scan(&id, &subscription, &interval); err != nil {
			return response, err
		}
		service := operations.GetTradespersonTradespersonIDTimeSlotsOKBodyItems0{}
		service.Interval = interval
		service.Subscription = subscription
		service.TimeSlots, err = getServiceTimeSlots(id)
		if err != nil {
			log.Printf("Failed to get service timeslots %s", err)
		}
		services = append(services, &service)
	}
	response.SetPayload(services)

	return response, nil
}

func getCustomerTimeSlots(timeSlotID int64) []*models.TimeSlotCustomersItems0 {
	customers := []*models.TimeSlotCustomersItems0{}

	stmt, err := db.Prepare("SELECT invoiceId, subscriptionId, cuStripeId FROM customer_time_slots WHERE timeSlotId=?")
	if err != nil {
		return customers
	}
	defer stmt.Close()

	rows, err := stmt.Query(timeSlotID)
	if err != nil {
		return customers
	}

	var invoiceID, subscriptionID sql.NullString
	var cuStripeID string
	for rows.Next() {
		if err := rows.Scan(&invoiceID, &subscriptionID, &cuStripeID); err != nil {
			return customers
		}
		customer := models.TimeSlotCustomersItems0{}
		if invoiceID.Valid {
			customer.InvoiceID = invoiceID.String
		} else if subscriptionID.Valid {
			customer.SubscriptionID = subscriptionID.String
		}
		customer.CuStripeID = cuStripeID
		customers = append(customers, &customer)
	}

	return customers
}

func GetTimeSlots(fixedPriceID int64) ([]*models.TimeSlot, error) {
	timeSlots := []*models.TimeSlot{}

	stmt, err := db.Prepare("SELECT id, startTime, segmentSize, curPeople, maxPeople FROM fixed_price_time_slots WHERE fixedPriceId=?")
	if err != nil {
		return timeSlots, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(fixedPriceID)
	if err != nil {
		return timeSlots, err
	}

	var ID, curPeople, maxPeople int64
	var startTime, segmentSize string
	for rows.Next() {
		if err := rows.Scan(&ID, &startTime, &segmentSize, &curPeople, &maxPeople); err != nil {
			return timeSlots, err
		}
		timeSlot := &models.TimeSlot{}
		timeSlot.ID = ID
		timeSlot.StartTime = startTime
		timeSlot.SegmentSize = segmentSize
		timeSlot.CurPeople = curPeople
		timeSlot.MaxPeople = maxPeople
		timeSlot.Customers = getCustomerTimeSlots(ID)
		timeSlots = append(timeSlots, timeSlot)
	}

	return timeSlots, nil
}

func GetPublicTimeSlots(fixedPriceID int64) ([]*models.TimeSlot, error) {
	timeSlots := []*models.TimeSlot{}

	stmt, err := db.Prepare("SELECT id, startTime, segmentSize, curPeople, maxPeople FROM fixed_price_time_slots WHERE fixedPriceId=?")
	if err != nil {
		return timeSlots, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(fixedPriceID)
	if err != nil {
		return timeSlots, err
	}

	var ID, curPeople, maxPeople int64
	var startTime, segmentSize string
	for rows.Next() {
		if err := rows.Scan(&ID, &startTime, &segmentSize, &curPeople, &maxPeople); err != nil {
			return timeSlots, err
		}
		timeSlot := &models.TimeSlot{}
		timeSlot.ID = ID
		timeSlot.StartTime = startTime
		timeSlot.SegmentSize = segmentSize
		timeSlot.CurPeople = curPeople
		timeSlot.MaxPeople = maxPeople
		timeSlots = append(timeSlots, timeSlot)
	}

	return timeSlots, nil
}

func GetAvailableTimeSlots(fixedPriceId int64, subscription bool) (int64, error) {
	availableTimeSlots := int64(0)

	var sqlStmt string
	if subscription {
		sqlStmt = "SELECT COUNT(id) FROM fixed_price_time_slots WHERE fixedPriceId=? AND NOT curPeople <=> maxPeople"
	} else {
		sqlStmt = "SELECT COUNT(id) FROM fixed_price_time_slots WHERE fixedPriceId=? AND startTime > CURDATE() AND NOT curPeople <=> maxPeople"
	}

	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return availableTimeSlots, err
	}
	defer stmt.Close()

	var count int64
	row := stmt.QueryRow(fixedPriceId)
	switch err = row.Scan(&count); err {
	case sql.ErrNoRows:
		return availableTimeSlots, err
	case nil:
		availableTimeSlots = count
	default:
		log.Printf("Unknown %v", err)
	}

	return availableTimeSlots, nil
}

func GetInvoiceStartTimeSegmentSize(invoiceID string) (string, string, error) {
	var startTime, segmentSize string
	stmt, err := db.Prepare("SELECT fpts.startTime, fpts.segmentSize FROM fixed_price_time_slots fpts INNER JOIN customer_time_slots cts ON fpts.id=cts.timeSlotId WHERE cts.invoiceId=?")
	if err != nil {
		return startTime, segmentSize, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(invoiceID)
	switch err = row.Scan(&startTime, &segmentSize); err {
	case sql.ErrNoRows:
		return startTime, segmentSize, err
	case nil:
		//
	default:
		log.Printf("Unknown %v", err)
	}

	return startTime, segmentSize, nil
}

func GetSubscriptionTimeSlot(subscriptionID string, fixedPriceID int64) (*models.TimeSlot, error) {
	timeSlot := &models.TimeSlot{}
	stmt, err := db.Prepare("SELECT fpts.startTime, fpts.segmentSize FROM fixed_price_time_slots fpts INNER JOIN customer_time_slots cts ON fpts.id=cts.timeSlotId WHERE cts.subscriptionId=? AND fpts.fixedPriceId=?")
	if err != nil {
		return timeSlot, err
	}
	defer stmt.Close()

	var startTime, segmentSize string
	row := stmt.QueryRow(subscriptionID, fixedPriceID)
	switch err = row.Scan(&startTime, &segmentSize); err {
	case sql.ErrNoRows:
		return timeSlot, err
	case nil:
		timeSlot.StartTime = startTime
		timeSlot.EndTime, err = internal.CreateEndTime(startTime, segmentSize)
		if err != nil {
			log.Printf("Failed to create endTime, %v", err)
		}
	default:
		log.Printf("Unknown %v", err)
	}

	return timeSlot, nil
}
