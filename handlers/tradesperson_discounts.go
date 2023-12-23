package handlers

import (
	"log"
	"redbudway-api/database"
	"redbudway-api/restapi/operations"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/coupon"
	"github.com/stripe/stripe-go/v72/promotioncode"
)

func PostTradespersonTradespersonIDCouponHandler(params operations.PostTradespersonTradespersonIDCouponParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	_coupon := params.Coupon

	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.PostTradespersonTradespersonIDCouponOKBody{}
	payload.Created = false
	response := operations.NewPostTradespersonTradespersonIDCouponOK()
	response.SetPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	couponParams := &stripe.CouponParams{
		Name: stripe.String(_coupon.Name),
	}
	if _coupon.Type == "percent_off" {
		couponParams.PercentOff = stripe.Float64(_coupon.Percent)
	} else if _coupon.Type == "amount_off" {
		amountInt := int64(_coupon.Amount * float64(100.00))
		couponParams.AmountOff = stripe.Int64(amountInt)
		couponParams.Currency = stripe.String("USD")
	}

	if _coupon.LimitedTime {
		tm, err := time.Parse(time.RFC3339, _coupon.RedeemBy)
		if err != nil {
			log.Printf("Failed to parse date time, %v", err)
			return response
		}
		expires := tm.Unix()
		couponParams.RedeemBy = &expires
	}
	if _coupon.MaxRedemption {
		couponParams.MaxRedemptions = &_coupon.MaxRedemptions
	}

	couponParams.Duration = stripe.String(_coupon.Duration)
	if _coupon.Duration != "once" {
		couponParams.DurationInMonths = stripe.Int64(_coupon.Months)
	}

	stripeCoupon, err := coupon.New(couponParams)
	if err != nil {
		log.Printf("Failed to create stripe coupon %v", err)
		return response
	}

	if _coupon.Type == "percent_off" {
		if err := database.CreatePercentOffCoupon(tradespersonID, stripeCoupon, _coupon); err != nil {
			log.Printf("Failed to save coupon %v", err)
			return response
		}
	} else if _coupon.Type == "amount_off" {
		if err := database.CreateAmountOffCoupon(tradespersonID, stripeCoupon, _coupon); err != nil {
			log.Printf("Failed to save coupon %v", err)
			return response
		}
	}

	payload.Created = true
	response.SetPayload(&payload)

	return response
}

func PutTradespersonTradespersonIDCouponCouponIDHandler(params operations.PutTradespersonTradespersonIDCouponCouponIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	couponID := params.CouponID
	_coupon := params.Coupon
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.PutTradespersonTradespersonIDCouponCouponIDOKBody{}
	payload.Updated = false
	response := operations.NewPutTradespersonTradespersonIDCouponCouponIDOK()
	response.SetPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	couponParams := &stripe.CouponParams{
		Name: stripe.String(_coupon.Name),
	}
	stripeCoupon, err := coupon.Update(couponID, couponParams)
	if err != nil {
		log.Printf("Failed to create stripe coupon %v", err)
		return response
	}

	if err := database.UpdateCoupon(tradespersonID, stripeCoupon.ID, _coupon); err != nil {
		log.Printf("Failed to save promo %v", err)
		return response
	}

	payload.Updated = true
	response.SetPayload(&payload)

	return response
}

func GetTradespersonTradespersonIDCouponCouponIDHandler(params operations.GetTradespersonTradespersonIDCouponCouponIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	couponID := params.CouponID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDCouponCouponIDOK()

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	coupon, err := database.GetCoupon(tradespersonID, couponID)
	if err != nil {
		log.Printf("Failed to retrieve coupon, %v", err)
		return response
	}

	response.SetPayload(&coupon)
	return response
}

func PostTradespersonTradespersonIDCouponCouponIDPromoHandler(params operations.PostTradespersonTradespersonIDCouponCouponIDPromoParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	couponID := params.CouponID
	promo := params.Promo
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.PostTradespersonTradespersonIDCouponCouponIDPromoOKBody{}
	payload.Created = false
	response := operations.NewPostTradespersonTradespersonIDCouponCouponIDPromoOK()
	response.SetPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	promoParams := &stripe.PromotionCodeParams{
		Active: stripe.Bool(true),
		Coupon: stripe.String(couponID),
	}

	if promo.Code != "" {
		promoParams.Code = stripe.String(promo.Code)
	}

	stripePromo, err := promotioncode.New(promoParams)
	if err != nil {
		log.Printf("Failed to create stripe promo %v", err)
		return response
	}

	if err := database.CreateCouponPromo(tradespersonID, couponID, stripePromo.ID, stripePromo.Code, stripePromo.Active); err != nil {
		log.Printf("Failed to save promo %v", err)
		return response
	}
	promo.ID = stripePromo.ID
	promo.Code = stripePromo.Code
	promo.Active = stripePromo.Active
	payload.Created = true
	payload.Promo = promo
	response.SetPayload(&payload)

	return response
}

func GetTradespersonTradespersonIDPromoPromoIDHandler(params operations.GetTradespersonTradespersonIDPromoPromoIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	promoID := params.PromoID
	token := params.HTTPRequest.Header.Get("Authorization")

	response := operations.NewGetTradespersonTradespersonIDPromoPromoIDOK()

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	promo, err := database.GetCouponPromo(tradespersonID, promoID)
	if err != nil {
		log.Printf("Failed to retrieve promo, %v", err)
		return response
	}

	response.SetPayload(promo)

	return response
}

func PutTradespersonTradespersonIDPromoPromoIDHandler(params operations.PutTradespersonTradespersonIDPromoPromoIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	promoID := params.PromoID
	promo := params.Promo
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.PutTradespersonTradespersonIDPromoPromoIDOKBody{}
	payload.Updated = false
	response := operations.NewPutTradespersonTradespersonIDPromoPromoIDOK()
	response.SetPayload(&payload)

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	promoParams := &stripe.PromotionCodeParams{
		Active: stripe.Bool(promo.Active),
	}
	stripePromo, err := promotioncode.Update(promoID, promoParams)
	if err != nil {
		log.Printf("Failed to create stripe promo %v", err)
		return response
	}

	if err := database.UpdateCouponPromo(tradespersonID, stripePromo.ID, stripePromo.Active); err != nil {
		log.Printf("Failed to save promo %v", err)
		return response
	}

	payload.Updated = true
	response.SetPayload(&payload)

	return response
}

func DeleteTradespersonTradespersonIDPromoPromoIDHandler(params operations.DeleteTradespersonTradespersonIDPromoPromoIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	promoID := params.PromoID
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.DeleteTradespersonTradespersonIDPromoPromoIDOKBody{}
	payload.Deleted = false
	response := operations.NewDeleteTradespersonTradespersonIDPromoPromoIDOK()

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	payload.Deleted, err = database.DeleteCouponPromo(tradespersonID, promoID)
	if err != nil {
		log.Printf("Failed to delete tradesperson promo, %v", err)
	}
	response.SetPayload(&payload)

	return response
}

func DeleteTradespersonTradespersonIDCouponCouponIDHandler(params operations.DeleteTradespersonTradespersonIDCouponCouponIDParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	couponID := params.CouponID
	token := params.HTTPRequest.Header.Get("Authorization")

	payload := operations.DeleteTradespersonTradespersonIDCouponCouponIDOKBody{}
	payload.Deleted = false
	response := operations.NewDeleteTradespersonTradespersonIDCouponCouponIDOK()

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	coupon, err := coupon.Del(couponID, nil)
	if err != nil {
		log.Printf("Failed to delete stripe coupon, %v", err)
		return response
	}

	if coupon.Deleted {
		payload.Deleted, err = database.DeleteCoupon(tradespersonID, couponID)
		if err != nil {
			log.Printf("Failed to delete tradesperson promo, %v", err)
		}
		response.SetPayload(&payload)
	}

	return response
}

func GetTradespersonTradespersonIDDiscountsHandler(params operations.GetTradespersonTradespersonIDDiscountsParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")
	response := operations.NewGetTradespersonTradespersonIDDiscountsOK()

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	coupons, err := database.GetCouponsWithPromos(tradespersonID)
	if err != nil {
		log.Printf("Failed to retrieve discounts %v", err)
		return response
	}

	response.SetPayload(coupons)

	return response
}

func GetTradespersonTradespersonIDServicesHandler(params operations.GetTradespersonTradespersonIDServicesParams, principal interface{}) middleware.Responder {
	tradespersonID := params.TradespersonID
	token := params.HTTPRequest.Header.Get("Authorization")
	response := operations.NewGetTradespersonTradespersonIDServicesOK()

	valid, err := ValidateTradespersonAccessToken(tradespersonID, token)
	if err != nil {
		log.Printf("Failed to validate tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	} else if !valid {
		log.Printf("Bad actor tradesperson %s, accessToken %s", tradespersonID, token)
		return response
	}

	fixedPrices := database.GetFixedPrices(tradespersonID)
	quotes := database.GetQuotes(tradespersonID)

	services := append(fixedPrices, quotes...)

	response.SetPayload(services)

	return response
}
