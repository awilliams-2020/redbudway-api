// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"redbudway-api/restapi/operations"
)

//go:generate swagger generate server --target ../../redbudway-api --name RedbudWayAPI --spec ../swagger.yaml --principal interface{}

func configureFlags(api *operations.RedbudWayAPIAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.RedbudWayAPIAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.DeleteTradespersonAccountTradespersonIDHandler == nil {
		api.DeleteTradespersonAccountTradespersonIDHandler = operations.DeleteTradespersonAccountTradespersonIDHandlerFunc(func(params operations.DeleteTradespersonAccountTradespersonIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.DeleteTradespersonAccountTradespersonID has not yet been implemented")
		})
	}
	if api.DeleteTradespersonTradespersonIDBillingCustomerSubscriptionsHandler == nil {
		api.DeleteTradespersonTradespersonIDBillingCustomerSubscriptionsHandler = operations.DeleteTradespersonTradespersonIDBillingCustomerSubscriptionsHandlerFunc(func(params operations.DeleteTradespersonTradespersonIDBillingCustomerSubscriptionsParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.DeleteTradespersonTradespersonIDBillingCustomerSubscriptions has not yet been implemented")
		})
	}
	if api.GetTradespersonAccountTradespersonIDHandler == nil {
		api.GetTradespersonAccountTradespersonIDHandler = operations.GetTradespersonAccountTradespersonIDHandlerFunc(func(params operations.GetTradespersonAccountTradespersonIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTradespersonAccountTradespersonID has not yet been implemented")
		})
	}
	if api.GetTradespersonAccountTradespersonIDSettingsHandler == nil {
		api.GetTradespersonAccountTradespersonIDSettingsHandler = operations.GetTradespersonAccountTradespersonIDSettingsHandlerFunc(func(params operations.GetTradespersonAccountTradespersonIDSettingsParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTradespersonAccountTradespersonIDSettings has not yet been implemented")
		})
	}
	if api.GetTradespersonAccountTradespersonIDStatusHandler == nil {
		api.GetTradespersonAccountTradespersonIDStatusHandler = operations.GetTradespersonAccountTradespersonIDStatusHandlerFunc(func(params operations.GetTradespersonAccountTradespersonIDStatusParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTradespersonAccountTradespersonIDStatus has not yet been implemented")
		})
	}
	if api.GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceHandler == nil {
		api.GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceHandler = operations.GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceHandlerFunc(func(params operations.GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoiceParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTradespersonTradespersonIDBillingCustomerSubscriptionInvoice has not yet been implemented")
		})
	}
	if api.GetTradespersonTradespersonIDBillingCustomerSubscriptionsHandler == nil {
		api.GetTradespersonTradespersonIDBillingCustomerSubscriptionsHandler = operations.GetTradespersonTradespersonIDBillingCustomerSubscriptionsHandlerFunc(func(params operations.GetTradespersonTradespersonIDBillingCustomerSubscriptionsParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTradespersonTradespersonIDBillingCustomerSubscriptions has not yet been implemented")
		})
	}
	if api.GetTradespersonTradespersonIDBillingCustomersHandler == nil {
		api.GetTradespersonTradespersonIDBillingCustomersHandler = operations.GetTradespersonTradespersonIDBillingCustomersHandlerFunc(func(params operations.GetTradespersonTradespersonIDBillingCustomersParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTradespersonTradespersonIDBillingCustomers has not yet been implemented")
		})
	}
	if api.GetTradespersonTradespersonIDBillingInvoicesHandler == nil {
		api.GetTradespersonTradespersonIDBillingInvoicesHandler = operations.GetTradespersonTradespersonIDBillingInvoicesHandlerFunc(func(params operations.GetTradespersonTradespersonIDBillingInvoicesParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTradespersonTradespersonIDBillingInvoices has not yet been implemented")
		})
	}
	if api.GetTradespersonTradespersonIDFixedPriceHandler == nil {
		api.GetTradespersonTradespersonIDFixedPriceHandler = operations.GetTradespersonTradespersonIDFixedPriceHandlerFunc(func(params operations.GetTradespersonTradespersonIDFixedPriceParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTradespersonTradespersonIDFixedPrice has not yet been implemented")
		})
	}
	if api.GetTradespersonTradespersonIDFixedPriceReviewsHandler == nil {
		api.GetTradespersonTradespersonIDFixedPriceReviewsHandler = operations.GetTradespersonTradespersonIDFixedPriceReviewsHandlerFunc(func(params operations.GetTradespersonTradespersonIDFixedPriceReviewsParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTradespersonTradespersonIDFixedPriceReviews has not yet been implemented")
		})
	}
	if api.GetTradespersonTradespersonIDFixedPricesHandler == nil {
		api.GetTradespersonTradespersonIDFixedPricesHandler = operations.GetTradespersonTradespersonIDFixedPricesHandlerFunc(func(params operations.GetTradespersonTradespersonIDFixedPricesParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTradespersonTradespersonIDFixedPrices has not yet been implemented")
		})
	}
	if api.GetTradespersonTradespersonIDLoginLinkHandler == nil {
		api.GetTradespersonTradespersonIDLoginLinkHandler = operations.GetTradespersonTradespersonIDLoginLinkHandlerFunc(func(params operations.GetTradespersonTradespersonIDLoginLinkParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTradespersonTradespersonIDLoginLink has not yet been implemented")
		})
	}
	if api.GetTradespersonTradespersonIDOnboardHandler == nil {
		api.GetTradespersonTradespersonIDOnboardHandler = operations.GetTradespersonTradespersonIDOnboardHandlerFunc(func(params operations.GetTradespersonTradespersonIDOnboardParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTradespersonTradespersonIDOnboard has not yet been implemented")
		})
	}
	if api.GetTradespersonTradespersonIDQuoteHandler == nil {
		api.GetTradespersonTradespersonIDQuoteHandler = operations.GetTradespersonTradespersonIDQuoteHandlerFunc(func(params operations.GetTradespersonTradespersonIDQuoteParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTradespersonTradespersonIDQuote has not yet been implemented")
		})
	}
	if api.GetTradespersonTradespersonIDQuoteReviewsHandler == nil {
		api.GetTradespersonTradespersonIDQuoteReviewsHandler = operations.GetTradespersonTradespersonIDQuoteReviewsHandlerFunc(func(params operations.GetTradespersonTradespersonIDQuoteReviewsParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTradespersonTradespersonIDQuoteReviews has not yet been implemented")
		})
	}
	if api.GetTradespersonTradespersonIDQuotesHandler == nil {
		api.GetTradespersonTradespersonIDQuotesHandler = operations.GetTradespersonTradespersonIDQuotesHandlerFunc(func(params operations.GetTradespersonTradespersonIDQuotesParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTradespersonTradespersonIDQuotes has not yet been implemented")
		})
	}
	if api.GetTradespersonTradespersonIDScheduleHandler == nil {
		api.GetTradespersonTradespersonIDScheduleHandler = operations.GetTradespersonTradespersonIDScheduleHandlerFunc(func(params operations.GetTradespersonTradespersonIDScheduleParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTradespersonTradespersonIDSchedule has not yet been implemented")
		})
	}
	if api.GetTradespersonTradespersonIDTimeSlotsHandler == nil {
		api.GetTradespersonTradespersonIDTimeSlotsHandler = operations.GetTradespersonTradespersonIDTimeSlotsHandlerFunc(func(params operations.GetTradespersonTradespersonIDTimeSlotsParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetTradespersonTradespersonIDTimeSlots has not yet been implemented")
		})
	}
	if api.PostForgotPasswordHandler == nil {
		api.PostForgotPasswordHandler = operations.PostForgotPasswordHandlerFunc(func(params operations.PostForgotPasswordParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostForgotPassword has not yet been implemented")
		})
	}
	if api.PostResetPasswordHandler == nil {
		api.PostResetPasswordHandler = operations.PostResetPasswordHandlerFunc(func(params operations.PostResetPasswordParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostResetPassword has not yet been implemented")
		})
	}
	if api.PostTradespersonAccountHandler == nil {
		api.PostTradespersonAccountHandler = operations.PostTradespersonAccountHandlerFunc(func(params operations.PostTradespersonAccountParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostTradespersonAccount has not yet been implemented")
		})
	}
	if api.PostTradespersonLoginHandler == nil {
		api.PostTradespersonLoginHandler = operations.PostTradespersonLoginHandlerFunc(func(params operations.PostTradespersonLoginParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostTradespersonLogin has not yet been implemented")
		})
	}
	if api.PostTradespersonTradespersonIDBillingCustomerSubscriptionRefundHandler == nil {
		api.PostTradespersonTradespersonIDBillingCustomerSubscriptionRefundHandler = operations.PostTradespersonTradespersonIDBillingCustomerSubscriptionRefundHandlerFunc(func(params operations.PostTradespersonTradespersonIDBillingCustomerSubscriptionRefundParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostTradespersonTradespersonIDBillingCustomerSubscriptionRefund has not yet been implemented")
		})
	}
	if api.PostTradespersonTradespersonIDFixedPriceHandler == nil {
		api.PostTradespersonTradespersonIDFixedPriceHandler = operations.PostTradespersonTradespersonIDFixedPriceHandlerFunc(func(params operations.PostTradespersonTradespersonIDFixedPriceParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostTradespersonTradespersonIDFixedPrice has not yet been implemented")
		})
	}
	if api.PostTradespersonTradespersonIDFixedPriceReviewHandler == nil {
		api.PostTradespersonTradespersonIDFixedPriceReviewHandler = operations.PostTradespersonTradespersonIDFixedPriceReviewHandlerFunc(func(params operations.PostTradespersonTradespersonIDFixedPriceReviewParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostTradespersonTradespersonIDFixedPriceReview has not yet been implemented")
		})
	}
	if api.PostTradespersonTradespersonIDQuoteHandler == nil {
		api.PostTradespersonTradespersonIDQuoteHandler = operations.PostTradespersonTradespersonIDQuoteHandlerFunc(func(params operations.PostTradespersonTradespersonIDQuoteParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostTradespersonTradespersonIDQuote has not yet been implemented")
		})
	}
	if api.PostTradespersonTradespersonIDQuoteReviewHandler == nil {
		api.PostTradespersonTradespersonIDQuoteReviewHandler = operations.PostTradespersonTradespersonIDQuoteReviewHandlerFunc(func(params operations.PostTradespersonTradespersonIDQuoteReviewParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostTradespersonTradespersonIDQuoteReview has not yet been implemented")
		})
	}
	if api.PutTradespersonAccountTradespersonIDHandler == nil {
		api.PutTradespersonAccountTradespersonIDHandler = operations.PutTradespersonAccountTradespersonIDHandlerFunc(func(params operations.PutTradespersonAccountTradespersonIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PutTradespersonAccountTradespersonID has not yet been implemented")
		})
	}
	if api.PutTradespersonAccountTradespersonIDSettingsHandler == nil {
		api.PutTradespersonAccountTradespersonIDSettingsHandler = operations.PutTradespersonAccountTradespersonIDSettingsHandlerFunc(func(params operations.PutTradespersonAccountTradespersonIDSettingsParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PutTradespersonAccountTradespersonIDSettings has not yet been implemented")
		})
	}
	if api.PutTradespersonTradespersonIDFixedPriceHandler == nil {
		api.PutTradespersonTradespersonIDFixedPriceHandler = operations.PutTradespersonTradespersonIDFixedPriceHandlerFunc(func(params operations.PutTradespersonTradespersonIDFixedPriceParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PutTradespersonTradespersonIDFixedPrice has not yet been implemented")
		})
	}
	if api.PutTradespersonTradespersonIDQuoteHandler == nil {
		api.PutTradespersonTradespersonIDQuoteHandler = operations.PutTradespersonTradespersonIDQuoteHandlerFunc(func(params operations.PutTradespersonTradespersonIDQuoteParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PutTradespersonTradespersonIDQuote has not yet been implemented")
		})
	}
	if api.PutTradespersonTradespersonIDUpdatePasswordHandler == nil {
		api.PutTradespersonTradespersonIDUpdatePasswordHandler = operations.PutTradespersonTradespersonIDUpdatePasswordHandlerFunc(func(params operations.PutTradespersonTradespersonIDUpdatePasswordParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PutTradespersonTradespersonIDUpdatePassword has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
