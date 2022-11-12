package database

import (
	"database/sql"
	"log"
	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"
	"time"
)

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

func ResetTakenTimeSlotByInvoice(invoiceID string) error {
	stmt, err := db.Prepare("UPDATE fixed_price_time_slots SET taken=False, takenBy=NULL, cuStripeId=NULL WHERE takenBy=? AND startTime > CURDATE()")
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

func UpdateTakenTimeSlot(stripeInvoiceID, cuStripeID, startTime string, fixedPriceID int64) (bool, error) {
	startDate, err := time.Parse("1/2/2006, 3:04:00 PM", startTime)
	if err != nil {
		log.Printf("Failed to parse startTime %s", startTime)
		return false, err
	}

	startTime = startDate.Format("2006-1-2 15:04:00")

	stmt, err := db.Prepare("UPDATE fixed_price_time_slots SET taken=true, takenBy=?, cuStripeId=? WHERE startTime=? AND fixedPriceId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(stripeInvoiceID, cuStripeID, startTime, fixedPriceID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func UpdateWeeklyTimeSlot(stripeInvoiceID, cuStripeID, startTime string, fixedPriceID int64) (bool, error) {
	startDate, err := time.Parse("1/2/2006, 3:04:00 PM", startTime)
	if err != nil {
		log.Printf("Failed to parse startTime %s", startTime)
		return false, err
	}
	startTime = startDate.Format("2006-1-2 15:04:00")

	stmt, err := db.Prepare("UPDATE fixed_price_time_slots SET taken=true, takenBy=?, cuStripeId=? WHERE DAYNAME(startTime)=DAYNAME(?) AND TIME(startTime)=TIME(?) AND fixedPriceId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(stripeInvoiceID, cuStripeID, startTime, startTime, fixedPriceID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func UpdateMonthlyTimeSlot(stripeInvoiceID, cuStripeID, startTime string, fixedPriceID int64) (bool, error) {
	startDate, err := time.Parse("1/2/2006, 3:04:00 PM", startTime)
	if err != nil {
		log.Printf("Failed to parse startTime %s", startTime)
		return false, err
	}
	startTime = startDate.Format("2006-1-2 15:04:00")

	stmt, err := db.Prepare("UPDATE fixed_price_time_slots SET taken=true, takenBy=?, cuStripeId=? WHERE DAYOFMONTH(startTime)=DAYOFMONTH(?) AND TIME(startTime)=TIME(?) AND fixedPriceId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(stripeInvoiceID, cuStripeID, startTime, startTime, fixedPriceID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func UpdateYearlyTimeSlot(stripeInvoiceID, cuStripeID, startTime string, fixedPriceID int64) (bool, error) {
	startDate, err := time.Parse("1/2/2006, 3:04:00 PM", startTime)
	if err != nil {
		log.Printf("Failed to parse startTime %s", startTime)
		return false, err
	}
	startTime = startDate.Format("2006-1-2 15:04:00")

	stmt, err := db.Prepare("UPDATE fixed_price_time_slots SET taken=true, takenBy=?, cuStripeId=? WHERE MONTH(startTime)=MONTH(?) AND DAYOFMONTH(startTime)=DAYOFMONTH(?) AND TIME(startTime)=TIME(?) AND fixedPriceId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(stripeInvoiceID, cuStripeID, startTime, startTime, startTime, fixedPriceID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func getScheduleTimeSlots(fixedPriceId int) ([]*operations.GetTradespersonTradespersonIDTimeSlotsOKBodyItems0TimeSlotsItems0, error) {
	timeSlots := []*operations.GetTradespersonTradespersonIDTimeSlotsOKBodyItems0TimeSlotsItems0{}

	stmt, err := db.Prepare("SELECT startTime, segmentSize, taken, takenBy FROM fixed_price_time_slots WHERE fixedPriceId=?")
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

	var taken bool
	var startTime, segmentSize string
	var takenBy sql.NullString
	for rows.Next() {
		if err := rows.Scan(&startTime, &segmentSize, &taken, &takenBy); err != nil {
			return timeSlots, err
		}
		timeSlot := &operations.GetTradespersonTradespersonIDTimeSlotsOKBodyItems0TimeSlotsItems0{}
		timeSlot.StartTime = startTime
		timeSlot.SegmentSize = segmentSize
		timeSlot.Taken = taken
		if takenBy.Valid {
			timeSlot.TakenBy = takenBy.String
		}
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
		service.TimeSlots, err = getScheduleTimeSlots(id)
		if err != nil {
			log.Printf("Failed to get service timeslots %s", err)
		}
		services = append(services, &service)
	}
	response.SetPayload(services)

	return response, nil
}

func GetTimeSlots(fixedPriceID int64) ([]*models.TimeSlot, error) {
	timeSlots := []*models.TimeSlot{}

	stmt, err := db.Prepare("SELECT ts.subscriptionId, ti.invoiceId, fpts.startTime, fpts.segmentSize, fpts.taken, fpts.takenBy, fpts.cuStripeId FROM fixed_price_time_slots fpts LEFT JOIN tradesperson_invoices ti ON fpts.takenBy=ti.invoiceId LEFT JOIN tradesperson_subscriptions ts ON fpts.takenBy=ts.subscriptionId WHERE fpts.fixedPriceId=?")
	if err != nil {
		return timeSlots, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(fixedPriceID)
	if err != nil {
		return timeSlots, err
	}

	var taken bool
	var subscriptionID, invoiceID, cuStripeID, takenBy sql.NullString
	var startTime, segmentSize string
	for rows.Next() {
		if err := rows.Scan(&subscriptionID, &invoiceID, &startTime, &segmentSize, &taken, &takenBy, &cuStripeID); err != nil {
			return timeSlots, err
		}
		timeSlot := &models.TimeSlot{}
		if subscriptionID.Valid {
			timeSlot.SubscriptionID = subscriptionID.String
		}
		if invoiceID.Valid {
			timeSlot.InvoiceID = invoiceID.String
		}
		timeSlot.StartTime = startTime
		timeSlot.SegmentSize = segmentSize
		timeSlot.Taken = taken
		if takenBy.Valid {
			timeSlot.TakenBy = takenBy.String
		}
		if cuStripeID.Valid {
			timeSlot.CuStripeID = cuStripeID.String
		}
		timeSlots = append(timeSlots, timeSlot)
	}

	return timeSlots, nil
}

func GetAvailableTimeSlots(fixedPriceId int64, subscription bool) (int64, error) {
	availableTimeSlots := int64(0)

	var sqlStmt string
	if subscription {
		sqlStmt = "SELECT COUNT(taken) FROM fixed_price_time_slots WHERE fixedPriceId=? AND taken=false"
	} else {
		sqlStmt = "SELECT COUNT(taken) FROM fixed_price_time_slots WHERE fixedPriceId=? AND startTime > CURDATE() AND taken=false"
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
	stmt, err := db.Prepare("SELECT startTime, segmentSize FROM fixed_price_time_slots WHERE takenBy=?")
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
	stmt, err := db.Prepare("SELECT startTime, segmentSize FROM fixed_price_time_slots WHERE takenBy=? AND fixedPriceId=?")
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
