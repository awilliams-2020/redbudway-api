package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

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
	strPrice := fmt.Sprintf("%.2f", stripePrice.UnitAmountDecimal/float64(100.00))
	floatPrice, err := strconv.ParseFloat(strPrice, 64)
	if err != nil {
		log.Printf("Failed to parse float, %v", err)
	}
	service.Price = floatPrice
	service.Title = stripeProduct.Name
	return service
}

func getCustomer(subscription bool, subscriptionID, invoiceID sql.NullString) (*models.Customer, int64, string, string) {
	_customer := &models.Customer{}
	quantity := int64(0)
	status := ""
	stripeID := ""

	if subscription {
		stripeSubscription, _ := sub.Get(
			subscriptionID.String,
			nil,
		)
		status = string(stripeSubscription.Status)
		quantity = stripeSubscription.Items.Data[0].Quantity
		stripeID = stripeSubscription.Customer.ID
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
		status = string(stripeInvoice.Status)
		quantity = stripeInvoice.Lines.Data[0].Quantity
		stripeID = stripeInvoice.Customer.ID
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
	return _customer, quantity, status, stripeID
}

func GetTradespersonSchedule(tradespersonID string, accessToken *string) (*operations.GetTradespersonTradespersonIDScheduleOK, error) {
	response := operations.NewGetTradespersonTradespersonIDScheduleOK()
	stmt, err := db.Prepare("SELECT fpts.startTime, fpts.endTime, cts.subscriptionId, cts.invoiceId, fp.priceId, fp.subscription, fp.subInterval FROM fixed_price_time_slots fpts INNER JOIN customer_time_slots cts ON fpts.id=cts.timeSlotId INNER JOIN fixed_prices fp ON fp.id=fpts.fixedPriceId INNER JOIN tradesperson_account ta ON ta.tradespersonId=fp.tradespersonId WHERE fp.tradespersonId=? AND (fpts.startTime > CURDATE() AND fp.subscription=False || fp.subscription=True ) AND cts.active=true GROUP BY fpts.id")
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
			customer := &operations.GetTradespersonTradespersonIDScheduleOKBodyServicesItems0TimeSlotsItems0CustomersItems0{}
			info, quantity, status, stripeID := getCustomer(subscription, subscriptionID, invoiceID)
			customer.Info = info
			customer.Quantity = quantity
			customer.Status = status
			customer.StripeID = stripeID
			if invoiceID.Valid {
				customer.InvoiceID = invoiceID.String
			}
			if subscriptionID.Valid {
				customer.SubscriptionID = subscriptionID.String
			}
			timeSlot.Customers = append(timeSlot.Customers, customer)
			timeSlots = append(timeSlots, timeSlot)
			service.TimeSlots = timeSlots
			m[priceId] = service
		} else {
			timeSlots := service.TimeSlots
			timeSlot := &operations.GetTradespersonTradespersonIDScheduleOKBodyServicesItems0TimeSlotsItems0{}
			timeSlot.StartTime = startTime
			timeSlot.EndTime = endTime
			customer := &operations.GetTradespersonTradespersonIDScheduleOKBodyServicesItems0TimeSlotsItems0CustomersItems0{}
			info, quantity, status, stripeID := getCustomer(subscription, subscriptionID, invoiceID)
			customer.Info = info
			customer.Quantity = quantity
			customer.Status = status
			customer.StripeID = stripeID
			if invoiceID.Valid {
				customer.InvoiceID = invoiceID.String
			}
			if subscriptionID.Valid {
				customer.SubscriptionID = subscriptionID.String
			}
			timeSlot.Customers = append(timeSlot.Customers, customer)
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
