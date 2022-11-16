package database

import (
	"database/sql"
	"log"

	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/invoice"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/sub"

	"redbudway-api/models"
	"redbudway-api/restapi/operations"
)

func getService(priceId string) *operations.GetTradespersonTradespersonIDScheduleOKBodyItems0 {
	service := &operations.GetTradespersonTradespersonIDScheduleOKBodyItems0{}

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

func getCustomer(subscription bool, subscriptionID, invoiceID sql.NullString) *models.Customer {
	_customer := &models.Customer{}

	if subscription {
		stripeSubscription, _ := sub.Get(
			subscriptionID.String,
			nil,
		)
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
	return _customer
}

func GetTradespersonSchedule(tradespersonID string) (*operations.GetTradespersonTradespersonIDScheduleOK, error) {
	response := operations.NewGetTradespersonTradespersonIDScheduleOK()
	stmt, err := db.Prepare("SELECT fpts.startTime, fpts.segmentSize, cts.subscriptionId, cts.invoiceId, fp.priceId, fp.subscription, fp.subInterval FROM fixed_price_time_slots fpts INNER JOIN customer_time_slots cts ON fpts.id=cts.timeSlotId INNER JOIN fixed_prices fp ON fp.id=fpts.fixedPriceId INNER JOIN tradesperson_account ta ON ta.tradespersonId=fp.tradespersonId WHERE ( (MONTH(fpts.startTime) = MONTH(CURRENT_DATE()) AND fp.subscription=False) || fp.subscription=True ) AND fp.tradespersonId=? GROUP BY fpts.id")
	if err != nil {
		return response, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID)
	if err != nil {
		return response, err
	}
	var startTime, subInterval, priceId string
	var segmentSize float64
	var subscription bool
	var subscriptionID, invoiceID sql.NullString
	m := make(map[string]*operations.GetTradespersonTradespersonIDScheduleOKBodyItems0)
	for rows.Next() {
		if err := rows.Scan(&startTime, &segmentSize, &subscriptionID, &invoiceID, &priceId, &subscription, &subInterval); err != nil {
			return response, err
		}
		service, exist := m[priceId]
		if !exist {
			service = getService(priceId)
			service.Subscription = subscription
			service.Interval = subInterval
			timeSlots := service.TimeSlots
			timeSlot := &operations.GetTradespersonTradespersonIDScheduleOKBodyItems0TimeSlotsItems0{}
			timeSlot.StartTime = startTime
			timeSlot.SegmentSize = segmentSize
			timeSlot.Customers = append(timeSlot.Customers, getCustomer(subscription, subscriptionID, invoiceID))
			timeSlots = append(timeSlots, timeSlot)
			service.TimeSlots = timeSlots
			m[priceId] = service
		} else {
			timeSlots := service.TimeSlots
			timeSlot := &operations.GetTradespersonTradespersonIDScheduleOKBodyItems0TimeSlotsItems0{}
			timeSlot.StartTime = startTime
			timeSlot.SegmentSize = segmentSize
			timeSlot.Customers = append(timeSlot.Customers, getCustomer(subscription, subscriptionID, invoiceID))
			timeSlots = append(timeSlots, timeSlot)
			service.TimeSlots = timeSlots
			m[priceId] = service
		}
	}

	payload := []*operations.GetTradespersonTradespersonIDScheduleOKBodyItems0{}
	for _, service := range m {
		payload = append(payload, service)
	}

	response.Payload = payload

	return response, nil
}
