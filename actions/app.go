package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/markbates/goth/gothic"
	"github.com/unrolled/secure"

	"ensetservice/models"

	"github.com/gobuffalo/buffalo-pop/v2/pop/popmw"
	csrf "github.com/gobuffalo/mw-csrf"
	i18n "github.com/gobuffalo/mw-i18n"
	"github.com/gobuffalo/packr/v2"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

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
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:           ENV,
			SessionName:   "_NOT_PHPSESSID", // I Hate PHP
			CompressFiles: true,
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		// Setup and use translations:
		app.Use(translations())

		app.GET("/", HomeHandler)
		app.Use(SetCurrentUser)
		app.Use(Authorize)
		app.Middleware.Skip(Authorize, HomeHandler)

		//Routes for Admins Auth
		AdminsAuth := app.Group("/auth/admins")
		AdminsAuth.GET("/login", AdminNew)
		AdminsAuth.POST("/login", AdminLogin)
		AdminsAuth.GET("/logout", AdminLogout)
		AdminsAuth.Use(Unauthorized)
		AdminsAuth.Middleware.Skip(Unauthorized, AdminLogout)
		AdminsAuth.Middleware.Skip(Authorize, AdminNew, AdminLogin)

		// //Routes for Admin registration
		// AdminsAuth.GET("/register", AdminRegisterNew)
		// AdminsAuth.POST("/register", AdminRegisterCreate)
		// AdminsAuth.Middleware.Remove(Authorize)

		//Routes for Students Auth
		StudentsAuth := app.Group("/auth/students")
		bah := buffalo.WrapHandlerFunc(gothic.BeginAuthHandler)
		StudentsAuth.GET("/login/{provider}", bah)
		StudentsAuth.GET("/login/{provider}/callback", StudentLogin)
		StudentsAuth.GET("/logout", StudentLogout)
		StudentsAuth.Use(Unauthorized)
		StudentsAuth.Middleware.Skip(Unauthorized, StudentLogout)
		StudentsAuth.Middleware.Skip(Authorize, bah, StudentLogin)

		//Routes for Documents Mgmt
		docRoute := app.Group("/documents")

		docRouteNew := docRoute.Group("/new")
		docRouteNew.GET("", DocumentNew)
		docRouteNew.POST("", DocumentCreate)
		docRouteNew.Use(StudentsOnly)

		docRouteProcessor := docRoute.Group("/process")
		docRouteProcessor.GET("/{docID}", DocumentProcessorNew)
		docRouteProcessor.POST("/{docID}", DocumentProcessorCreate)
		docRouteProcessor.Use(AdminsOnly) // make sure that users have admin account can access

		docRouteDownloader := docRoute.Group("/download")
		docRouteDownloader.GET("/{docID}", DocumentDownloader)
		// docRouteDownloader.Use(AdminsOnly) // make sure that users have admin account can access

		app.ServeFiles("/", assetsBox) // serve files from the public directory

		// Handle Error pages

		// Handle 404 Error
		app.ErrorHandlers[404] = func(status int, err error, c buffalo.Context) error {
			return c.Render(http.StatusNotFound, r.String(fmt.Sprintf("404 Page Not Found : %v", err.Error())))
		}

		// Handle 500 Error
		app.ErrorHandlers[500] = func(status int, err error, c buffalo.Context) error {
			return c.Render(http.StatusNotFound, r.String(fmt.Sprintf("500 Internal Server Error : %v", err.Error())))
		}
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
		// SSLRedirect:           ENV == "production",
		SSLProxyHeaders:    map[string]string{"X-Forwarded-Proto": "https"},
		BrowserXssFilter:   true,
		ContentTypeNosniff: true,
		FrameDeny:          true,
		// STSSeconds:           31536000,
		// STSIncludeSubdomains: true,
		// STSPreload:           true,
		// ContentSecurityPolicy: "script-src $NONCE",
	})
}
