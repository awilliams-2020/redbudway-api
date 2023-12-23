package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"redbudway-api/models"
	"redbudway-api/restapi/operations"
	"strconv"

	"github.com/stripe/stripe-go/v72"
)

func CreateAmountOffCoupon(tradespersonId string, stripeCoupon *stripe.Coupon, coupon *models.Coupon) error {
	stmt, err := db.Prepare("INSERT INTO tradesperson_coupons (tradespersonId, couponId, name, discountType, amountOff, duration, months, maxRedemptions, redeemBy, services, subscriptions) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	services, err := json.Marshal(coupon.Services)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(tradespersonId, stripeCoupon.ID, coupon.Name, coupon.Type, stripeCoupon.AmountOff, coupon.Duration, coupon.Months, coupon.MaxRedemptions, stripeCoupon.RedeemBy, string(services), coupon.Subscriptions)
	if err != nil {
		return err
	}

	return nil
}

func CreatePercentOffCoupon(tradespersonId string, stripeCoupon *stripe.Coupon, coupon *models.Coupon) error {
	stmt, err := db.Prepare("INSERT INTO tradesperson_coupons (tradespersonId, couponId, name, discountType, percentOff, duration, months, maxRedemptions, redeemBy, services, subscriptions) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	services, err := json.Marshal(coupon.Services)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(tradespersonId, stripeCoupon.ID, coupon.Name, coupon.Type, stripeCoupon.PercentOff, coupon.Duration, coupon.Months, coupon.MaxRedemptions, stripeCoupon.RedeemBy, string(services), coupon.Subscriptions)
	if err != nil {
		return err
	}

	return nil
}

func UpdateCoupon(tradespersonID, couponID string, coupon operations.PutTradespersonTradespersonIDCouponCouponIDBody) error {
	stmt, err := db.Prepare("UPDATE tradesperson_coupons SET name=?, services=? WHERE tradespersonId=? AND couponId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	services, err := json.Marshal(coupon.Services)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(coupon.Name, string(services), tradespersonID, couponID)
	if err != nil {
		return err
	}

	return nil
}

func GetCoupon(tradespersonId, couponId string) (models.Coupon, error) {

	coupon := models.Coupon{}
	stmt, err := db.Prepare("SELECT name, discountType, IFNULL(amountOff, 0), IFNULL(percentOff, 0), duration, months, maxRedemption, redeemBy, services FROM tradesperson_coupons WHERE tradespersonId=? AND couponId=?")
	if err != nil {
		return coupon, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(tradespersonId, couponId)
	switch err = row.Scan(&coupon.Name, &coupon.Type, &coupon.Amount, &coupon.Percent, &coupon.Duration, &coupon.Months, &coupon.MaxRedemptions, &coupon.RedeemBy, &coupon.Services); err {
	case sql.ErrNoRows:
		//
	case nil:
		//
	default:
		log.Printf("Unknown %v", err)
	}

	return coupon, nil
}

func DeleteCoupon(tradespersonID, couponID string) (bool, error) {
	stmt, err := db.Prepare("DELETE FROM tradesperson_coupons WHERE tradespersonId=? AND couponId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(tradespersonID, couponID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func getCouponPromos(tradespersonID, couponID string) ([]*models.Promo, error) {
	promos := []*models.Promo{}

	stmt, err := db.Prepare("SELECT promoId, code, active FROM tradesperson_promos WHERE tradespersonId=? AND couponId=?")
	if err != nil {
		log.Printf("Failed to create prepare statement, %v", err)
		return promos, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID, couponID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return promos, err
	}
	for rows.Next() {
		promo := &models.Promo{}
		if err := rows.Scan(&promo.ID, &promo.Code, &promo.Active); err != nil {
			log.Printf("Failed to scan row %v", err)
			return promos, err
		}
		promos = append(promos, promo)
	}
	return promos, nil
}

func GetCouponPromo(tradespersonID, promoID string) (*models.Promo, error) {
	promo := &models.Promo{}

	// stmt, err := db.Prepare("SELECT code, active, expires, redemptions, redeemed, amount, services FROM tradesperson_promos WHERE tradespersonId=? AND promoId=?")
	// if err != nil {
	// 	log.Printf("Failed to create prepare statement, %v", err)
	// 	return promo, err
	// }
	// defer stmt.Close()

	// var expires sql.NullString
	// var maxRedemptions sql.NullInt64
	// var amount sql.NullFloat64
	// var services string
	// row := stmt.QueryRow(tradespersonID, promoID)
	// switch err = row.Scan(&promo.Code, &promo.Active, &expires, &maxRedemptions, &promo.Redemptions, &amount, &services); err {
	// case sql.ErrNoRows:
	// 	//
	// case nil:
	// 	promo.Expires = expires.String
	// 	promo.Limited = false
	// 	if promo.Expires != "0" {
	// 		promo.Limited = true
	// 		i, err := strconv.ParseInt(promo.Expires, 10, 64)
	// 		if err != nil {
	// 			log.Printf("Failed to parse string to int, %v", err)
	// 		} else {
	// 			promo.Expires = fmt.Sprintf("%d", i)
	// 		}
	// 	}
	// 	promo.MaxRedemptions = maxRedemptions.Int64
	// 	promo.MaxRedemption = false
	// 	if promo.MaxRedemptions > int64(0) {
	// 		promo.MaxRedemption = true
	// 	}
	// 	promo.Amount = amount.Float64
	// 	promo.Minimum = false
	// 	if promo.Amount > float64(0) {
	// 		promo.Minimum = true
	// 	}
	// 	if err := json.Unmarshal([]byte(services), &promo.Services); err != nil {
	// 		log.Printf("Failed to unmarshal services into promo", err)
	// 	}
	// default:
	// 	log.Printf("Unknown %v", err)
	// }

	return promo, nil
}

func DeleteCouponPromo(tradespersonID, promoID string) (bool, error) {
	stmt, err := db.Prepare("DELETE FROM tradesperson_promos WHERE tradespersonId=? AND promoId=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(tradespersonID, promoID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func CreateCouponPromo(tradespersonId, couponId, promoId, code string, active bool) error {
	stmt, err := db.Prepare("INSERT INTO tradesperson_promos (tradespersonId, couponId, promoId, code, active) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tradespersonId, couponId, promoId, code, active)
	if err != nil {
		return err
	}

	return nil
}

func UpdateCouponPromo(tradespersonId, promoId string, active bool) error {
	stmt, err := db.Prepare("UPDATE tradesperson_promos SET active=? WHERE tradespersonId=? AND promoId=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(active, tradespersonId, promoId)
	if err != nil {
		return err
	}

	return nil
}

func GetCouponsWithPromos(tradespersonID string) ([]*models.Coupon, error) {
	coupons := []*models.Coupon{}

	stmt, err := db.Prepare("SELECT couponId, name, discountType, amountOff, percentOff, duration, months, maxRedemptions, redeemBy, services, timesRedeemed, subscriptions FROM tradesperson_coupons WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create prepare statement, %v", err)
		return coupons, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return coupons, err
	}

	var months, maxRedemptions sql.NullInt64
	var redeemBy sql.NullString
	var amount, percent sql.NullFloat64
	var services string
	for rows.Next() {
		coupon := &models.Coupon{}
		if err := rows.Scan(&coupon.ID, &coupon.Name, &coupon.Type, &amount, &percent, &coupon.Duration, &months, &maxRedemptions, &redeemBy, &services, &coupon.TimesRedeemed, &coupon.Subscriptions); err != nil {
			log.Printf("Failed to scan row %v", err)
			return coupons, err
		}
		strAmount := fmt.Sprintf("%.2f", amount.Float64/float64(100.00))
		floatAmount, err := strconv.ParseFloat(strAmount, 64)
		if err != nil {

			return coupons, err
		}
		coupon.Amount = floatAmount
		coupon.Percent = percent.Float64
		coupon.Months = months.Int64
		coupon.MaxRedemptions = maxRedemptions.Int64
		coupon.RedeemBy = redeemBy.String
		promos, err := getCouponPromos(tradespersonID, coupon.ID)
		if err != nil {

		}
		coupon.Promos = promos
		if err := json.Unmarshal([]byte(services), &coupon.Services); err != nil {
			log.Printf("Failed to unmarshal services into promo", err)
		}
		coupons = append(coupons, coupon)
	}
	return coupons, nil
}

func GetFixedPrices(tradespersonID string) []*operations.GetTradespersonTradespersonIDServicesOKBodyItems0 {

	fixedPrices := []*operations.GetTradespersonTradespersonIDServicesOKBodyItems0{}
	stmt, err := db.Prepare("SELECT priceId, title, subscription FROM fixed_prices WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return fixedPrices
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return fixedPrices
	}

	for rows.Next() {
		service := operations.GetTradespersonTradespersonIDServicesOKBodyItems0{}
		if err := rows.Scan(&service.ID, &service.Title, &service.Subscription); err != nil {
			log.Printf("Failed to scan row %v", err)
			return fixedPrices
		}
		fixedPrices = append(fixedPrices, &service)
	}

	return fixedPrices
}

func GetQuotes(tradespersonID string) []*operations.GetTradespersonTradespersonIDServicesOKBodyItems0 {

	quotes := []*operations.GetTradespersonTradespersonIDServicesOKBodyItems0{}
	stmt, err := db.Prepare("SELECT quote, title FROM quotes WHERE tradespersonId=?")
	if err != nil {
		log.Printf("Failed to create select statement %s", err)
		return quotes
	}
	defer stmt.Close()

	rows, err := stmt.Query(tradespersonID)
	if err != nil {
		log.Printf("Failed to execute select statement %s", err)
		return quotes
	}

	for rows.Next() {
		service := operations.GetTradespersonTradespersonIDServicesOKBodyItems0{}
		service.Subscription = false
		if err := rows.Scan(&service.ID, &service.Title); err != nil {
			log.Printf("Failed to scan row %v", err)
			return quotes
		}
		quotes = append(quotes, &service)
	}

	return quotes
}
