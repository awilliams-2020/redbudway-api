package database

import (
	"database/sql"
	"log"

	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/invoice"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/sub"

	"redbudway-api/internal"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"
)

func getService(priceId string) *operations.GetTradespersonTradespersonIDScheduleOKBodyServicesItems0 {
	service := &operations.GetTradespersonTradespersonIDScheduleOKBodyServicesItems0{}

	stripePrice, _ := price.Get(
		priceId,
		nil,
	)
	stripeProduct, _ := product.Get(
		stripePrice.Product.ID,
		nil,
	)
	service.Price = stripePrice.UnitAmountDecimal
	service.Title = stripeProduct.Name
	return service
}

func getCustomer(subscription bool, subscriptionID, invoiceID sql.NullString) (*models.Customer, int64) {
	_customer := &models.Customer{}
	quantity := int64(0)

	if subscription {
		stripeSubscription, _ := sub.Get(
			subscriptionID.String,
			nil,
		)
		quantity = stripeSubscription.Items.Data[0].Quantity
		stripeCustomer, err := customer.Get(stripeSubscription.Customer.ID, nil)
		if err != nil {
			log.Printf("Failed to get stripe customer %v", err)
		}
		_customer.Name = stripeCustomer.Name

		_customer.Address = &models.Address{
			City:    stripeCustomer.Address.City,
			LineOne: stripeCustomer.Address.Line1,
			LineTwo: stripeCustomer.Address.Line2,
			ZipCode: stripeCustomer.Address.PostalCode,
			State:   stripeCustomer.Address.State,
		}
		_customer.Email = stripeCustomer.Email
		_customer.Phone = stripeCustomer.Phone
	} else {
		stripeInvoice, _ := invoice.Get(
			invoiceID.String,
			nil,
		)
		quantity = stripeInvoice.Lines.Data[0].Quantity
		_customer.Name = *stripeInvoice.CustomerName
		_customer.Address = &models.Address{
			City:    stripeInvoice.CustomerAddress.City,
			LineOne: stripeInvoice.CustomerAddress.Line1,
			LineTwo: stripeInvoice.CustomerAddress.Line2,
			ZipCode: stripeInvoice.CustomerAddress.PostalCode,
			State:   stripeInvoice.CustomerAddress.State,
		}
		_customer.Email = stripeInvoice.CustomerEmail
		_customer.Phone = *stripeInvoice.CustomerPhone
	}
	return _customer, quantity
}

func GetTradespersonSchedule(tradespersonID string, accessToken *string) (*operations.GetTradespersonTradespersonIDScheduleOK, error) {
	response := operations.NewGetTradespersonTradespersonIDScheduleOK()
	stmt, err := db.Prepare("SELECT fpts.startTime, fpts.endTime, cts.subscriptionId, cts.invoiceId, fp.priceId, fp.subscription, fp.subInterval FROM fixed_price_time_slots fpts INNER JOIN customer_time_slots cts ON fpts.id=cts.timeSlotId INNER JOIN fixed_prices fp ON fp.id=fpts.fixedPriceId INNER JOIN tradesperson_account ta ON ta.tradespersonId=fp.tradespersonId WHERE ( (DATE(fpts.startTime) >= CURRENT_DATE() AND fp.subscription=False) || fp.subscription=True ) AND fp.tradespersonId=? AND cts.active=true GROUP BY fpts.id")
	if err != nil {
		return response, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID)
	if err != nil {
		return response, err
	}
	var startTime, subInterval, priceId string
	var endTime string
	var subscription bool
	var subscriptionID, invoiceID sql.NullString
	m := make(map[string]*operations.GetTradespersonTradespersonIDScheduleOKBodyServicesItems0)
	for rows.Next() {
		if err := rows.Scan(&startTime, &endTime, &subscriptionID, &invoiceID, &priceId, &subscription, &subInterval); err != nil {
			return response, err
		}
		service, exist := m[priceId]
		if !exist {
			service = getService(priceId)
			service.Subscription = subscription
			service.Interval = subInterval
			timeSlots := service.TimeSlots
			timeSlot := &operations.GetTradespersonTradespersonIDScheduleOKBodyServicesItems0TimeSlotsItems0{}
			timeSlot.StartTime = startTime
			timeSlot.EndTime = endTime
			customer, quantity := getCustomer(subscription, subscriptionID, invoiceID)
			timeSlot.Customers = append(timeSlot.Customers, customer)
			timeSlot.Quantity = quantity
			timeSlots = append(timeSlots, timeSlot)
			service.TimeSlots = timeSlots
			m[priceId] = service
		} else {
			timeSlots := service.TimeSlots
			timeSlot := &operations.GetTradespersonTradespersonIDScheduleOKBodyServicesItems0TimeSlotsItems0{}
			timeSlot.StartTime = startTime
			timeSlot.EndTime = endTime
			customer, quantity := getCustomer(subscription, subscriptionID, invoiceID)
			timeSlot.Customers = append(timeSlot.Customers, customer)
			timeSlot.Quantity = quantity
			timeSlots = append(timeSlots, timeSlot)
			service.TimeSlots = timeSlots
			m[priceId] = service
		}
	}

	payload := &operations.GetTradespersonTradespersonIDScheduleOKBody{}
	payload.Services = []*operations.GetTradespersonTradespersonIDScheduleOKBodyServicesItems0{}
	for _, service := range m {
		payload.Services = append(payload.Services, service)
	}

	response.SetPayload(payload)

	googleTimeSlots := models.GoogleTimeSlots{}
	if accessToken != nil {
		googleTimeSlots = internal.GetGoogleTimeSlots(*accessToken)
		if err != nil {
			log.Printf("Failed to get google time slots, %s", err)
			return response, nil
		}
	}
	payload.GoogleTimeSlots = googleTimeSlots
	response.SetPayload(payload)

	return response, nil
}
