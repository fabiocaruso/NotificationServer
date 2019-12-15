package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/unrolled/secure"

	//csrf "github.com/gobuffalo/mw-csrf"
	i18n "github.com/gobuffalo/mw-i18n"
	"github.com/gobuffalo/packr/v2"

	//"log"
	"crypto/ed25519"
	"gopkg.in/couchbase/gocb.v1"
	"github.com/fabiocaruso/NotificationServer/services"
	"fmt"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator
var privateKey ed25519.PrivateKey
var publicKey ed25519.PublicKey
var nsdb *gocb.Cluster
var nsBucket *gocb.Bucket

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	// The Seed string needs to be length 32
	// TODO: get the seed from env variable
	privateKey = ed25519.NewKeyFromSeed([]byte("b{2'*&-kjECuLynMZaE7@f:yzD}$MND?"))
	publicKey = (privateKey.Public()).(ed25519.PublicKey)
	/*var err error
	publicKey, privateKey, err = ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatal("Could not generate public- and privateKey: ", err.Error())
	}*/
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_notification_server_session",
		})

		var connErr error
		nsdb, connErr = connDB()
		if connErr != nil {
			fmt.Println("Error: " + connErr.Error())
		}

		var bucketErr error
		nsBucket, bucketErr = getBucket(nsdb, "NotificationServer")
		if bucketErr != nil {
			fmt.Println("Error: " + bucketErr.Error())
		}

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		//app.Use(csrf.New)

		// Setup and use translations:
		app.Use(translations())

		// API Version 1
		v1 := app.Group("/api/v1")

		// Ressources
		usersr := &UsersResource{}
		udr := &UserDevicesResource{}

		// Services
		for name, instance := range services.Services {
			service := instance.(services.Service)
			v1.GET("/services/" + name + "/{botToken}", service.WebhookHandler)
			v1.POST("/services/" + name + "/{botToken}", service.WebhookHandler)
		}

		v1.GET("/", HomeHandler)

		v1.POST("/auth", authHandler)

		v1.GET("/users", usersr.Show)
		v1.PUT("/users", usersr.Update)
		v1.POST("/users", usersr.Add)
		v1.DELETE("/users", usersr.Delete)
		
		v1.GET("/user/{apikey:[a-z0-9]+}/devices", udr.Show)
		v1.PUT("/user/{apikey:[a-z0-9]+}/devices", udr.Update)
		v1.POST("/user/{apikey:[a-z0-9]+}/devices", udr.Add)
		v1.DELETE("/user/{apikey:[a-z0-9]+}/devices", udr.Delete)

		v1.POST("/user/{apikey:[a-z0-9]+}/sendMessage", sendMessageHandler)
		v1.POST("/user/{apikey:[a-z0-9]+}/sendmessage", sendMessageHandler)

		v1.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(packr.New("app:locales", "../locales"), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
