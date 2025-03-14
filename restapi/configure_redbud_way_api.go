// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"
	"os"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/stripe/stripe-go/v72"

	"redbudway-api/database"
	"redbudway-api/handlers"
	"redbudway-api/internal"
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

	// Applies when the "Authorization" header is set
	api.BearerAuth = internal.ValidateToken

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()

	database.Init()

	stripe.Key = os.Getenv("STRIPE_KEY")

	api.DeleteCustomerCustomerIDHandler = operations.DeleteCustomerCustomerIDHandlerFunc(handlers.DeleteCustomerCustomerIDHandler)

	api.DeleteTradespersonTradespersonIDHandler = operations.DeleteTradespersonTradespersonIDHandlerFunc(handlers.DeleteTradespersonTradespersonIDHandler)

	api.DeleteTradespersonTradespersonIDBillingInvoiceInvoiceIDHandler = operations.DeleteTradespersonTradespersonIDBillingInvoiceInvoiceIDHandlerFunc(handlers.DeleteTradespersonTradespersonIDBillingInvoiceInvoiceIDHandler)

	api.DeleteTradespersonTradespersonIDGoogleTokenHandler = operations.DeleteTradespersonTradespersonIDGoogleTokenHandlerFunc(handlers.DeleteTradespersonTradespersonIDGoogleTokenHandler)

	api.DeleteTradespersonTradespersonIDPromoPromoIDHandler = operations.DeleteTradespersonTradespersonIDPromoPromoIDHandlerFunc(handlers.DeleteTradespersonTradespersonIDPromoPromoIDHandler)

	api.DeleteTradespersonTradespersonIDCouponCouponIDHandler = operations.DeleteTradespersonTradespersonIDCouponCouponIDHandlerFunc(handlers.DeleteTradespersonTradespersonIDCouponCouponIDHandler)

	api.GetAdminAdminIDTradespeopleHandler = operations.GetAdminAdminIDTradespeopleHandlerFunc(handlers.GetAdminAdminIDTradespeopleHandler)

	api.GetAdminAdminIDAccessTokenHandler = operations.GetAdminAdminIDAccessTokenHandlerFunc(handlers.GetAdminAdminIDAccessTokenHandler)

	api.GetCustomerCustomerIDAccessTokenHandler = operations.GetCustomerCustomerIDAccessTokenHandlerFunc(handlers.GetCustomerCustomerIDAccessTokenHandler)

	api.GetCustomerCustomerIDBillingLinkHandler = operations.GetCustomerCustomerIDBillingLinkHandlerFunc(handlers.GetCustomerCustomerIDBillingLinkHandler)

	api.GetCustomerCustomerIDPaymentDefaultHandler = operations.GetCustomerCustomerIDPaymentDefaultHandlerFunc(handlers.GetCustomerCustomerIDPaymentDefaultHandler)

	api.GetCustomerCustomerIDFixedPricePriceIDReviewHandler = operations.GetCustomerCustomerIDFixedPricePriceIDReviewHandlerFunc(handlers.GetCustomerCustomerIDFixedPricePriceIDReviewHandler)

	api.GetCustomerCustomerIDQuoteQuoteIDReviewHandler = operations.GetCustomerCustomerIDQuoteQuoteIDReviewHandlerFunc(handlers.GetCustomerCustomerIDQuoteQuoteIDReviewHandler)

	api.GetCustomerCustomerIDSubscriptionPriceIDReviewHandler = operations.GetCustomerCustomerIDSubscriptionPriceIDReviewHandlerFunc(handlers.GetCustomerCustomerIDSubscriptionPriceIDReviewHandler)

	api.GetCustomerCustomerIDReverifyHandler = operations.GetCustomerCustomerIDReverifyHandlerFunc(handlers.GetCustomerCustomerIDReverifyHandler)

	api.GetCustomerCustomerIDQuotesHandler = operations.GetCustomerCustomerIDQuotesHandlerFunc(handlers.GetCustomerCustomerIDQuotesHandler)

	api.GetCustomerCustomerIDQuoteQuoteIDHandler = operations.GetCustomerCustomerIDQuoteQuoteIDHandlerFunc(handlers.GetCustomerCustomerIDQuoteQuoteIDHandler)

	api.GetCustomerCustomerIDPromoHandler = operations.GetCustomerCustomerIDPromoHandlerFunc(handlers.GetCustomerCustomerIDPromoHandler)

	api.GetFixedPricePriceIDHandler = operations.GetFixedPricePriceIDHandlerFunc(handlers.GetFixedPricePriceIDHandler)

	api.GetFixedPricePriceIDReviewsHandler = operations.GetFixedPricePriceIDReviewsHandlerFunc(handlers.GetFixedPricePriceIDReviewsHandler)

	api.GetFixedPricesHandler = operations.GetFixedPricesHandlerFunc(handlers.GetFixedPricesHandler)

	api.GetFixedPricePagesHandler = operations.GetFixedPricePagesHandlerFunc(handlers.GetFixedPricePagesHandler)

	api.GetForgotPasswordHandler = operations.GetForgotPasswordHandlerFunc(handlers.GetForgotPasswordHandler)

	api.GetAddressHandler = operations.GetAddressHandlerFunc(handlers.GetAddressHandler)

	api.GetLocationHandler = operations.GetLocationHandlerFunc(handlers.GetLocationHandler)

	api.GetQuoteQuoteIDHandler = operations.GetQuoteQuoteIDHandlerFunc(handlers.GetQuoteQuoteIDHandler)

	api.GetQuoteQuoteIDReviewsHandler = operations.GetQuoteQuoteIDReviewsHandlerFunc(handlers.GetQuoteQuoteIDReviewsHandler)

	api.GetQuotesHandler = operations.GetQuotesHandlerFunc(handlers.GetQuotesHandler)

	api.GetQuotePagesHandler = operations.GetQuotePagesHandlerFunc(handlers.GetQuotePagesHandler)

	api.GetProfileVanityOrIDHandler = operations.GetProfileVanityOrIDHandlerFunc(handlers.GetProfileVanityOrIDHandler)

	api.GetProfileVanityOrIDFixedPricesHandler = operations.GetProfileVanityOrIDFixedPricesHandlerFunc(handlers.GetProfileVanityOrIDFixedPricesHandler)

	api.GetProfileVanityOrIDQuotesHandler = operations.GetProfileVanityOrIDQuotesHandlerFunc(handlers.GetProfileVanityOrIDQuotesHandler)

	api.GetTradespersonTradespersonIDAccessTokenHandler = operations.GetTradespersonTradespersonIDAccessTokenHandlerFunc(handlers.GetTradespersonTradespersonIDAccessTokenHandler)

	api.GetTradespersonTradespersonIDHandler = operations.GetTradespersonTradespersonIDHandlerFunc(handlers.GetTradespersonTradespersonIDHandler)

	api.GetTradespersonTradespersonIDServicesHandler = operations.GetTradespersonTradespersonIDServicesHandlerFunc(handlers.GetTradespersonTradespersonIDServicesHandler)

	api.GetTradespersonTradespersonIDCouponCouponIDHandler = operations.GetTradespersonTradespersonIDCouponCouponIDHandlerFunc(handlers.GetTradespersonTradespersonIDCouponCouponIDHandler)

	api.GetTradespersonTradespersonIDProfileHandler = operations.GetTradespersonTradespersonIDProfileHandlerFunc(handlers.GetTradespersonTradespersonIDProfileHandler)

	api.GetTradespersonTradespersonIDBrandingHandler = operations.GetTradespersonTradespersonIDBrandingHandlerFunc(handlers.GetTradespersonTradespersonIDBrandingHandler)

	api.GetTradespersonTradespersonIDSyncHandler = operations.GetTradespersonTradespersonIDSyncHandlerFunc(handlers.GetTradespersonTradespersonIDSyncHandler)

	api.GetTradespersonTradespersonIDLoginLinkHandler = operations.GetTradespersonTradespersonIDLoginLinkHandlerFunc(handlers.GetTradespersonTradespersonIDLoginLinkHandler)

	api.GetTradespersonTradespersonIDSettingsHandler = operations.GetTradespersonTradespersonIDSettingsHandlerFunc(handlers.GetTradespersonTradespersonIDSettingsHandler)

	api.GetTradespersonTradespersonIDBillingSubscriptionsHandler = operations.GetTradespersonTradespersonIDBillingSubscriptionsHandlerFunc(handlers.GetTradespersonTradespersonIDBillingSubscriptionsHandler)

	api.GetTradespersonTradespersonIDBillingSubscriptionPagesHandler = operations.GetTradespersonTradespersonIDBillingSubscriptionPagesHandlerFunc(handlers.GetTradespersonTradespersonIDBillingSubscriptionPagesHandler)

	api.GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionsHandler = operations.GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionsHandlerFunc(handlers.GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionsHandler)

	api.GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDHandler = operations.GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDHandlerFunc(handlers.GetTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDHandler)

	api.GetTradespersonTradespersonIDBillingCustomersHandler = operations.GetTradespersonTradespersonIDBillingCustomersHandlerFunc(handlers.GetTradespersonTradespersonIDBillingCustomersHandler)

	api.GetTradespersonTradespersonIDBillingInvoiceInvoiceIDHandler = operations.GetTradespersonTradespersonIDBillingInvoiceInvoiceIDHandlerFunc(handlers.GetTradespersonTradespersonIDBillingInvoiceInvoiceIDHandler)

	api.GetTradespersonTradespersonIDBillingInvoicePagesHandler = operations.GetTradespersonTradespersonIDBillingInvoicePagesHandlerFunc(handlers.GetTradespersonTradespersonIDBillingInvoicePagesHandler)

	api.GetTradespersonTradespersonIDBillingInvoicesHandler = operations.GetTradespersonTradespersonIDBillingInvoicesHandlerFunc(handlers.GetTradespersonTradespersonIDBillingInvoicesHandler)

	api.GetTradespersonTradespersonIDBillingQuotesHandler = operations.GetTradespersonTradespersonIDBillingQuotesHandlerFunc(handlers.GetTradespersonTradespersonIDBillingQuotesHandler)

	api.GetTradespersonTradespersonIDBillingQuotePagesHandler = operations.GetTradespersonTradespersonIDBillingQuotePagesHandlerFunc(handlers.GetTradespersonTradespersonIDBillingQuotePagesHandler)

	api.GetTradespersonTradespersonIDBillingManualInvoicesHandler = operations.GetTradespersonTradespersonIDBillingManualInvoicesHandlerFunc(handlers.GetTradespersonTradespersonIDBillingManualInvoicesHandler)

	api.GetTradespersonTradespersonIDBillingManualInvoicePagesHandler = operations.GetTradespersonTradespersonIDBillingManualInvoicePagesHandlerFunc(handlers.GetTradespersonTradespersonIDBillingManualInvoicePagesHandler)

	api.GetTradespersonTradespersonIDBillingManualInvoiceInvoiceIDHandler = operations.GetTradespersonTradespersonIDBillingManualInvoiceInvoiceIDHandlerFunc(handlers.GetTradespersonTradespersonIDBillingManualInvoiceInvoiceIDHandler)

	api.GetTradespersonTradespersonIDBillingQuoteQuoteIDHandler = operations.GetTradespersonTradespersonIDBillingQuoteQuoteIDHandlerFunc(handlers.GetTradespersonTradespersonIDBillingQuoteQuoteIDHandler)

	api.GetTradespersonTradespersonIDBillingQuoteQuoteIDPdfHandler = operations.GetTradespersonTradespersonIDBillingQuoteQuoteIDPdfHandlerFunc(handlers.GetTradespersonTradespersonIDBillingQuoteQuoteIDPdfHandler)

	api.GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDHandler = operations.GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDHandlerFunc(handlers.GetTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDHandler)

	api.GetTradespersonTradespersonIDFixedPricePriceIDHandler = operations.GetTradespersonTradespersonIDFixedPricePriceIDHandlerFunc(handlers.GetTradespersonTradespersonIDFixedPricePriceIDHandler)

	api.GetTradespersonTradespersonIDFixedPriceReviewsHandler = operations.GetTradespersonTradespersonIDFixedPriceReviewsHandlerFunc(handlers.GetTradespersonTradespersonIDFixedPriceReviewsHandler)

	api.GetTradespersonTradespersonIDFixedPricesHandler = operations.GetTradespersonTradespersonIDFixedPricesHandlerFunc(handlers.GetTradespersonTradespersonIDFixedPricesHandler)

	api.GetTradespersonTradespersonIDFixedPricePagesHandler = operations.GetTradespersonTradespersonIDFixedPricePagesHandlerFunc(handlers.GetTradespersonTradespersonIDFixedPricePagesHandler)

	api.GetTradespersonTradespersonIDOnboardHandler = operations.GetTradespersonTradespersonIDOnboardHandlerFunc(handlers.GetTradespersonTradespersonIDOnboardHandler)

	api.GetTradespersonTradespersonIDQuoteQuoteIDHandler = operations.GetTradespersonTradespersonIDQuoteQuoteIDHandlerFunc(handlers.GetTradespersonTradespersonIDQuoteQuoteIDHandler)

	api.GetTradespersonTradespersonIDQuoteReviewsHandler = operations.GetTradespersonTradespersonIDQuoteReviewsHandlerFunc(handlers.GetTradespersonTradespersonIDQuoteReviewsHandler)

	api.GetTradespersonTradespersonIDQuotesHandler = operations.GetTradespersonTradespersonIDQuotesHandlerFunc(handlers.GetTradespersonTradespersonIDQuotesHandler)

	api.GetTradespersonTradespersonIDQuotePagesHandler = operations.GetTradespersonTradespersonIDQuotePagesHandlerFunc(handlers.GetTradespersonTradespersonIDQuotePagesHandler)

	api.GetTradespersonTradespersonIDScheduleHandler = operations.GetTradespersonTradespersonIDScheduleHandlerFunc(handlers.GetTradespersonTradespersonIDScheduleHandler)

	api.GetTradespersonTradespersonIDTimeSlotsHandler = operations.GetTradespersonTradespersonIDTimeSlotsHandlerFunc(handlers.GetTradespersonTradespersonIDTimeSlotsHandler)

	api.GetTradespersonTradespersonIDDiscountsHandler = operations.GetTradespersonTradespersonIDDiscountsHandlerFunc(handlers.GetTradespersonTradespersonIDDiscountsHandler)

	api.GetTradespersonTradespersonIDPromoPromoIDHandler = operations.GetTradespersonTradespersonIDPromoPromoIDHandlerFunc(handlers.GetTradespersonTradespersonIDPromoPromoIDHandler)

	api.GetCustomerCustomerIDVerifyHandler = operations.GetCustomerCustomerIDVerifyHandlerFunc(handlers.GetCustomerCustomerIDVerifyHandler)

	api.PostCustomerHandler = operations.PostCustomerHandlerFunc(handlers.PostCustomerHandler)

	api.PostCustomerLoginHandler = operations.PostCustomerLoginHandlerFunc(handlers.PostCustomerLoginHandler)

	api.PostCustomerCustomerIDLogoutHandler = operations.PostCustomerCustomerIDLogoutHandlerFunc(handlers.PostCustomerCustomerIDLogoutHandler)

	api.PostCustomerCustomerIDFixedPricePriceIDReviewHandler = operations.PostCustomerCustomerIDFixedPricePriceIDReviewHandlerFunc(handlers.PostCustomerCustomerIDFixedPricePriceIDReviewHandler)

	api.PostCustomerCustomerIDQuoteQuoteIDReviewHandler = operations.PostCustomerCustomerIDQuoteQuoteIDReviewHandlerFunc(handlers.PostCustomerCustomerIDQuoteQuoteIDReviewHandler)

	api.PostCustomerCustomerIDAccessTokenHandler = operations.PostCustomerCustomerIDAccessTokenHandlerFunc(handlers.PostCustomerCustomerIDAccessTokenHandler)

	api.PostCustomerCustomerIDFixedPricePriceIDBookHandler = operations.PostCustomerCustomerIDFixedPricePriceIDBookHandlerFunc(handlers.PostCustomerCustomerIDFixedPricePriceIDBookHandler)

	api.PostCustomerCustomerIDSubscriptionPriceIDBookHandler = operations.PostCustomerCustomerIDSubscriptionPriceIDBookHandlerFunc(handlers.PostCustomerCustomerIDSubscriptionPriceIDBookHandler)

	api.PostCustomerCustomerIDQuoteQuoteIDRequestHandler = operations.PostCustomerCustomerIDQuoteQuoteIDRequestHandlerFunc(handlers.PostCustomerCustomerIDQuoteQuoteIDRequestHandler)

	api.PostCustomerCustomerIDQuoteQuoteIDAcceptHandler = operations.PostCustomerCustomerIDQuoteQuoteIDAcceptHandlerFunc(handlers.PostCustomerCustomerIDQuoteQuoteIDAcceptHandler)

	api.PostCustomerCustomerIDVerifyHandler = operations.PostCustomerCustomerIDVerifyHandlerFunc(handlers.PostCustomerCustomerIDVerifyHandler)

	api.PostResetPasswordHandler = operations.PostResetPasswordHandlerFunc(handlers.PostResetPasswordHandler)

	api.PostTradespersonHandler = operations.PostTradespersonHandlerFunc(handlers.PostTradespersonHandler)

	api.PostTradespersonLoginHandler = operations.PostTradespersonLoginHandlerFunc(handlers.PostTradespersonLoginHandler)

	api.PostTradespersonTradespersonIDLogoutHandler = operations.PostTradespersonTradespersonIDLogoutHandlerFunc(handlers.PostTradespersonTradespersonIDLogoutHandler)

	api.PostTradespersonTradespersonIDAccessTokenHandler = operations.PostTradespersonTradespersonIDAccessTokenHandlerFunc(handlers.PostTradespersonTradespersonIDAccessTokenHandler)

	api.PostTradespersonTradespersonIDGoogleTokenHandler = operations.PostTradespersonTradespersonIDGoogleTokenHandlerFunc(handlers.PostTradespersonTradespersonIDGoogleTokenHandler)

	api.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalizeHandler = operations.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalizeHandlerFunc(handlers.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDFinalizeHandler)

	api.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDVoidHandler = operations.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDVoidHandlerFunc(handlers.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDVoidHandler)

	api.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDRefundHandler = operations.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDRefundHandlerFunc(handlers.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDRefundHandler)

	api.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDUncollectibleHandler = operations.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDUncollectibleHandlerFunc(handlers.PostTradespersonTradespersonIDBillingInvoiceInvoiceIDUncollectibleHandler)

	api.PostTradespersonTradespersonIDBillingManualInvoiceHandler = operations.PostTradespersonTradespersonIDBillingManualInvoiceHandlerFunc(handlers.PostTradespersonTradespersonIDBillingManualInvoiceHandler)

	api.PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDVoidHandler = operations.PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDVoidHandlerFunc(handlers.PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDVoidHandler)

	api.PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleHandler = operations.PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleHandlerFunc(handlers.PostTradespersonTradespersonIDBillingManualInvoiceInvoiceIDUncollectibleHandler)

	api.PostTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionsCancelHandler = operations.PostTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionsCancelHandlerFunc(handlers.PostTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionsCancelHandler)

	api.PostTradespersonTradespersonIDBillingQuoteQuoteIDCancelHandler = operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDCancelHandlerFunc(handlers.PostTradespersonTradespersonIDBillingQuoteQuoteIDCancelHandler)

	api.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDFinalizeHandler = operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDFinalizeHandlerFunc(handlers.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDFinalizeHandler)

	api.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleHandler = operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleHandlerFunc(handlers.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDUncollectibleHandler)

	api.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundHandler = operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundHandlerFunc(handlers.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDRefundHandler)

	api.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDVoidHandler = operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDVoidHandlerFunc(handlers.PostTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDVoidHandler)

	api.PostTradespersonTradespersonIDBillingQuoteQuoteIDFinalizeHandler = operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDFinalizeHandlerFunc(handlers.PostTradespersonTradespersonIDBillingQuoteQuoteIDFinalizeHandler)

	api.PostTradespersonTradespersonIDBillingQuoteQuoteIDReviseHandler = operations.PostTradespersonTradespersonIDBillingQuoteQuoteIDReviseHandlerFunc(handlers.PostTradespersonTradespersonIDBillingQuoteQuoteIDReviseHandler)

	api.PostTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDRefundHandler = operations.PostTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDRefundHandlerFunc(handlers.PostTradespersonTradespersonIDBillingCustomerStripeIDSubscriptionSubscriptionIDInvoiceInvoiceIDRefundHandler)

	api.PostTradespersonTradespersonIDCouponHandler = operations.PostTradespersonTradespersonIDCouponHandlerFunc(handlers.PostTradespersonTradespersonIDCouponHandler)

	api.PostTradespersonTradespersonIDCouponCouponIDPromoHandler = operations.PostTradespersonTradespersonIDCouponCouponIDPromoHandlerFunc(handlers.PostTradespersonTradespersonIDCouponCouponIDPromoHandler)

	api.PostTradespersonTradespersonIDEmailHandler = operations.PostTradespersonTradespersonIDEmailHandlerFunc(handlers.PostTradespersonTradespersonIDEmailHandler)

	api.PostTradespersonTradespersonIDFixedPriceHandler = operations.PostTradespersonTradespersonIDFixedPriceHandlerFunc(handlers.PostTradespersonTradespersonIDFixedPriceHandler)

	api.PostTradespersonTradespersonIDFixedPriceReviewHandler = operations.PostTradespersonTradespersonIDFixedPriceReviewHandlerFunc(handlers.PostTradespersonTradespersonIDFixedPriceReviewHandler)

	api.PostTradespersonTradespersonIDQuoteHandler = operations.PostTradespersonTradespersonIDQuoteHandlerFunc(handlers.PostTradespersonTradespersonIDQuoteHandler)

	api.PostTradespersonTradespersonIDQuoteReviewHandler = operations.PostTradespersonTradespersonIDQuoteReviewHandlerFunc(handlers.PostTradespersonTradespersonIDQuoteReviewHandler)

	api.PutCustomerCustomerIDHandler = operations.PutCustomerCustomerIDHandlerFunc(handlers.PutCustomerCustomerIDHandler)

	api.PutTradespersonTradespersonIDHandler = operations.PutTradespersonTradespersonIDHandlerFunc(handlers.PutTradespersonTradespersonIDHandler)

	api.PutTradespersonTradespersonIDProfileHandler = operations.PutTradespersonTradespersonIDProfileHandlerFunc(handlers.PutTradespersonTradespersonIDProfileHandler)

	api.PutTradespersonTradespersonIDBrandingHandler = operations.PutTradespersonTradespersonIDBrandingHandlerFunc(handlers.PutTradespersonTradespersonIDBrandingHandler)

	api.PutTradespersonTradespersonIDGoogleTokenHandler = operations.PutTradespersonTradespersonIDGoogleTokenHandlerFunc(handlers.PutTradespersonTradespersonIDGoogleTokenHandler)

	api.PutTradespersonTradespersonIDBillingInvoiceInvoiceIDHandler = operations.PutTradespersonTradespersonIDBillingInvoiceInvoiceIDHandlerFunc(handlers.PutTradespersonTradespersonIDBillingInvoiceInvoiceIDHandler)

	api.PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDHandler = operations.PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDHandlerFunc(handlers.PutTradespersonTradespersonIDBillingQuoteQuoteIDInvoiceInvoiceIDHandler)

	api.PutTradespersonTradespersonIDSettingsHandler = operations.PutTradespersonTradespersonIDSettingsHandlerFunc(handlers.PutTradespersonTradespersonIDSettingsHandler)

	api.PutTradespersonTradespersonIDTimeZoneHandler = operations.PutTradespersonTradespersonIDTimeZoneHandlerFunc(handlers.PutTradespersonTradespersonIDTimeZoneHandler)

	api.PutTradespersonTradespersonIDFixedPricePriceIDHandler = operations.PutTradespersonTradespersonIDFixedPricePriceIDHandlerFunc(handlers.PutTradespersonTradespersonIDFixedPricePriceIDHandler)

	api.PutTradespersonTradespersonIDQuoteQuoteIDHandler = operations.PutTradespersonTradespersonIDQuoteQuoteIDHandlerFunc(handlers.PutTradespersonTradespersonIDQuoteQuoteIDHandler)

	api.PutTradespersonTradespersonIDBillingQuoteQuoteIDHandler = operations.PutTradespersonTradespersonIDBillingQuoteQuoteIDHandlerFunc(handlers.PutTradespersonTradespersonIDBillingQuoteQuoteIDHandler)

	api.PutTradespersonTradespersonIDPromoPromoIDHandler = operations.PutTradespersonTradespersonIDPromoPromoIDHandlerFunc(handlers.PutTradespersonTradespersonIDPromoPromoIDHandler)

	api.PutTradespersonTradespersonIDCouponCouponIDHandler = operations.PutTradespersonTradespersonIDCouponCouponIDHandlerFunc(handlers.PutTradespersonTradespersonIDCouponCouponIDHandler)

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
