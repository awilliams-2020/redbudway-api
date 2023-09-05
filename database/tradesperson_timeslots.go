package database

import (
	"database/sql"
	"encoding/json"
	"log"
	"redbudway-api/models"
	"time"
)

func InsertTimeSlots(fixedPriceID int64, fixedPrice *models.ServiceDetails) error {
	for _, timeSlot := range fixedPrice.TimeSlots {
		t, err := time.Parse("1/2/2006, 3:04:00 PM", timeSlot.StartTime)
		if err != nil {
			return err
		}
		startTime := t.Format("2006-01-02 15:04:00")
		t, err = time.Parse("1/2/2006, 3:04:00 PM", timeSlot.EndTime)
		if err != nil {
			return err
		}
		endTime := t.Format("2006-01-02 15:04:00")
		stmt, err := db.Prepare("INSERT INTO fixed_price_time_slots (fixedPriceId, startTime, endTime, bookings) VALUES (?, ?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(fixedPriceID, startTime, endTime, timeSlot.Bookings)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateTimeSlots(fixedPriceID int64, fixedPrice *models.ServiceDetails) error {
	stmt, err := db.Prepare("SELECT startTime, endTime FROM fixed_price_time_slots WHERE fixedPriceId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(fixedPriceID)
	if err != nil {
		return err
	}

	existingTimeSlots := []*models.TimeSlotsItems0{}
	var startTime, endTime string
	for rows.Next() {
		if err := rows.Scan(&startTime, &endTime); err != nil {
			return err
		}
		existingTimeSlot := &models.TimeSlotsItems0{}
		existingTimeSlot.StartTime = startTime
		existingTimeSlot.EndTime = endTime
		existingTimeSlots = append(existingTimeSlots, existingTimeSlot)
	}

	for _, existingTimeSlot := range existingTimeSlots {
		found := false
		for _, timeSlot := range fixedPrice.TimeSlots {
			t, err := time.Parse("1/2/2006, 3:04:00 PM", timeSlot.StartTime)
			if err != nil {
				return err
			}
			startTime := t.Format("2006-01-02 15:04:00")
			t, err = time.Parse("1/2/2006, 3:04:00 PM", timeSlot.EndTime)
			if err != nil {
				return err
			}
			endTime := t.Format("2006-01-02 15:04:00")
			if startTime == existingTimeSlot.StartTime {
				found = true
				if timeSlot.EndTime != existingTimeSlot.EndTime {
					stmt, err := db.Prepare("UPDATE fixed_price_time_slots SET endTime=?, bookings=? WHERE fixedPriceId=? AND startTime=?")
					if err != nil {
						return err
					}
					defer stmt.Close()
					_, err = stmt.Exec(endTime, timeSlot.Bookings, fixedPriceID, startTime)
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
		startTime := t.Format("2006-01-02 15:04:00")
		t, err = time.Parse("1/2/2006, 3:04:00 PM", timeSlot.EndTime)
		if err != nil {
			return err
		}
		endTime := t.Format("2006-01-02 15:04:00")
		found := false
		for _, existingTimeSlot := range existingTimeSlots {
			if startTime == existingTimeSlot.StartTime {
				found = true
			}
		}
		if !found {
			stmt, err := db.Prepare("INSERT INTO fixed_price_time_slots (fixedPriceId, startTime, endTime, bookings) VALUES (?, ?, ?, ?)")
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(fixedPriceID, startTime, endTime, timeSlot.Bookings)
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

func UpdateTimeSlotByInvoice(invoiceID string, quantity int64) error {

	stmt, err := db.Prepare("UPDATE fixed_price_time_slots fpts INNER JOIN customer_time_slots cts ON fpts.id=cts.timeSlotId SET fpts.booked=fpts.booked-?, cts.active=false WHERE cts.invoiceId=? AND fpts.startTime > CURDATE()")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(quantity, invoiceID)
	if err != nil {
		return err
	}

	return nil
}

func GetTakenTimeSlot(startTime string, fixedPriceID int64) ([]string, []string, error) {

	cuStripeIDs := []string{}
	takenBy := []string{}

	stmt, err := db.Prepare("SELECT takenBy, cuStripeIds, bookings, booked FROM fixed_price_time_slots WHERE startTime=? AND fixedPriceId=?")
	if err != nil {
		return cuStripeIDs, takenBy, err
	}
	defer stmt.Close()

	var t, c sql.NullString
	var bookings, booked int64
	err = stmt.QueryRow(startTime, fixedPriceID).Scan(&t, &c, &bookings, &booked)
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

func updateCustomerSubscriptionTimeSlot(timeSlotID int64, stripeInvoiceID, cuStripeID string) error {

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

func UpdateTakenTimeSlot(stripeInvoiceID, cuStripeID, startTime string, fixedPriceID, quantity, timeSlotID int64) (bool, error) {
	startDate, err := time.Parse("1/2/2006, 3:04:00 PM", startTime)
	if err != nil {
		log.Printf("Failed to parse startTime %s", startTime)
		return false, err
	}

	startTime = startDate.Format("2006-1-2 15:04:00")

	if err := updateCustomerSubscriptionTimeSlot(timeSlotID, stripeInvoiceID, cuStripeID); err != nil {
		log.Printf("Failed to insert customer time slot, %v", err)
		return false, err
	}

	stmt, err := db.Prepare("UPDATE fixed_price_time_slots SET booked=booked+? WHERE startTime=? AND fixedPriceId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(quantity, startTime, fixedPriceID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func insertCustomerSubscriptionTimeSlot(timeSlotID int64, subscriptionID, cuStripeID string) error {

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

func UpdateWeeklyTimeSlot(subscriptionID, cuStripeID, startTime string, fixedPriceID, quantity, timeSlotID int64) (bool, error) {
	startDate, err := time.Parse("1/2/2006, 3:04:00 PM", startTime)
	if err != nil {
		log.Printf("Failed to parse startTime %s", startTime)
		return false, err
	}
	startTime = startDate.Format("2006-1-2 15:04:00")

	if err := insertCustomerSubscriptionTimeSlot(timeSlotID, subscriptionID, cuStripeID); err != nil {
		log.Printf("Failed to insert customer time slot, %v", err)
		return false, err
	}

	stmt, err := db.Prepare("UPDATE fixed_price_time_slots SET booked=booked+? WHERE DAYNAME(startTime)=DAYNAME(?) AND TIME(startTime)=TIME(?) AND fixedPriceId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(quantity, startTime, startTime, fixedPriceID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func UpdateMonthlyTimeSlot(subscriptionID, cuStripeID, startTime string, fixedPriceID, quantity, timeSlotID int64) (bool, error) {
	startDate, err := time.Parse("1/2/2006, 3:04:00 PM", startTime)
	if err != nil {
		log.Printf("Failed to parse startTime %s", startTime)
		return false, err
	}
	startTime = startDate.Format("2006-1-2 15:04:00")

	if err := insertCustomerSubscriptionTimeSlot(timeSlotID, subscriptionID, cuStripeID); err != nil {
		log.Printf("Failed to insert customer time slot, %v", err)
		return false, err
	}

	stmt, err := db.Prepare("UPDATE fixed_price_time_slots SET booked=booked+? WHERE DAYOFMONTH(startTime)=DAYOFMONTH(?) AND TIME(startTime)=TIME(?) AND fixedPriceId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(quantity, startTime, startTime, fixedPriceID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func UpdateYearlyTimeSlot(subscriptionID, cuStripeID, startTime string, fixedPriceID, quantity, timeSlotID int64) (bool, error) {
	startDate, err := time.Parse("1/2/2006, 3:04:00 PM", startTime)
	if err != nil {
		log.Printf("Failed to parse startTime %s", startTime)
		return false, err
	}
	startTime = startDate.Format("2006-1-2 15:04:00")

	if err := insertCustomerSubscriptionTimeSlot(timeSlotID, subscriptionID, cuStripeID); err != nil {
		log.Printf("Failed to insert customer time slot, %v", err)
		return false, err
	}

	stmt, err := db.Prepare("UPDATE fixed_price_time_slots SET booked=booked+? WHERE MONTH(startTime)=MONTH(?) AND DAYOFMONTH(startTime)=DAYOFMONTH(?) AND TIME(startTime)=TIME(?) AND fixedPriceId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(quantity, startTime, startTime, startTime, fixedPriceID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
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

	stmt, err := db.Prepare("SELECT id, startTime, endTime, booked, bookings FROM fixed_price_time_slots WHERE fixedPriceId=?")
	if err != nil {
		return timeSlots, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(fixedPriceID)
	if err != nil {
		return timeSlots, err
	}

	var ID, booked, bookings int64
	var startTime, endTime string
	for rows.Next() {
		if err := rows.Scan(&ID, &startTime, &endTime, &booked, &bookings); err != nil {
			return timeSlots, err
		}
		timeSlot := &models.TimeSlot{}
		timeSlot.ID = ID
		timeSlot.StartTime = startTime
		timeSlot.EndTime = endTime
		timeSlot.Booked = booked
		timeSlot.Bookings = bookings
		timeSlot.Customers = getCustomerTimeSlots(ID)
		timeSlots = append(timeSlots, timeSlot)
	}

	return timeSlots, nil
}

func GetPublicTimeSlots(fixedPriceID int64, subscription bool) ([]*models.TimeSlot, error) {
	timeSlots := []*models.TimeSlot{}

	var sql string
	if subscription {
		sql = "SELECT id, startTime, endTime, booked, bookings FROM fixed_price_time_slots WHERE fixedPriceId=?"
	} else {
		sql = "SELECT id, startTime, endTime, booked, bookings FROM fixed_price_time_slots WHERE fixedPriceId=? AND startTime > CURDATE()"
	}

	stmt, err := db.Prepare(sql)
	if err != nil {
		return timeSlots, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(fixedPriceID)
	if err != nil {
		return timeSlots, err
	}

	var ID, booked, bookings int64
	var startTime, endTime string
	for rows.Next() {
		if err := rows.Scan(&ID, &startTime, &endTime, &booked, &bookings); err != nil {
			return timeSlots, err
		}
		timeSlot := &models.TimeSlot{}
		timeSlot.ID = ID
		timeSlot.StartTime = startTime
		timeSlot.EndTime = endTime
		timeSlot.Booked = booked
		timeSlot.Bookings = bookings
		timeSlots = append(timeSlots, timeSlot)
	}

	return timeSlots, nil
}

func GetAvailableTimeSlots(fixedPriceId int64, subscription bool) (int64, error) {
	availableTimeSlots := int64(0)
	var sqlStmt string
	if subscription {
		sqlStmt = "SELECT COUNT(id) FROM fixed_price_time_slots WHERE fixedPriceId=? AND NOT booked <=> bookings"
	} else {
		sqlStmt = "SELECT COUNT(id) FROM fixed_price_time_slots WHERE fixedPriceId=? AND startTime > CURDATE() AND NOT booked <=> bookings"
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

func GetInvoiceStartTimeEndTime(invoiceID string) (string, string, error) {
	var startTime, endTime string
	stmt, err := db.Prepare("SELECT fpts.startTime, fpts.endTime FROM fixed_price_time_slots fpts INNER JOIN customer_time_slots cts ON fpts.id=cts.timeSlotId WHERE cts.invoiceId=?")
	if err != nil {
		return startTime, endTime, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(invoiceID)
	switch err = row.Scan(&startTime, &endTime); err {
	case sql.ErrNoRows:
		return startTime, endTime, err
	case nil:
		//
	default:
		log.Printf("Unknown %v", err)
	}

	return startTime, endTime, nil
}

func GetSubscriptionTimeSlot(subscriptionID string, fixedPriceID int64) (*models.TimeSlot, error) {
	timeSlot := &models.TimeSlot{}
	stmt, err := db.Prepare("SELECT fpts.startTime, fpts.endTime FROM fixed_price_time_slots fpts INNER JOIN customer_time_slots cts ON fpts.id=cts.timeSlotId WHERE cts.subscriptionId=? AND fpts.fixedPriceId=?")
	if err != nil {
		return timeSlot, err
	}
	defer stmt.Close()

	var startTime, endTime string
	row := stmt.QueryRow(subscriptionID, fixedPriceID)
	switch err = row.Scan(&startTime, &endTime); err {
	case sql.ErrNoRows:
		return timeSlot, err
	case nil:
		timeSlot.StartTime = startTime
		timeSlot.EndTime = endTime
		if err != nil {
			log.Printf("Failed to create endTime, %v", err)
		}
	default:
		log.Printf("Unknown %v", err)
	}

	return timeSlot, nil
}
