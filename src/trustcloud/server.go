package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/render"
	"net/http"
	"trustcloud/api"
	"trustcloud/legacy"
	"trustcloud/partner"
	"trustcloud/util"
)

var m *martini.Martini

/*
Production : proxied through Apache
See:   /etc/apache2/sites-available/api.trustcloud.com - on production server
*/
func init() {

	m = martini.New()

	m.Use(render.Renderer())

	m.Use(martini.Recovery())
	m.Use(martini.Logger())

	m.Use(martini.Static("public"))

	m.Map(util.WebLog)

	m.Use(CheckForwarder)

	m.Use(MapEncoder)

	//	allowCORSHandler := cors.Allow(&cors.Options{
	//		AllowAllOrigins: true,
	//		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE"},
	//	})

	allowCORSHandler := cors.Allow(&cors.Options{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin"},
	})

	m.Use(allowCORSHandler)

	// Inject database

	m.MapTo(api.DataBase, (*api.DB)(nil))
	m.MapTo(legacy.DataBase, (*legacy.DB)(nil))

	r := martini.NewRouter()

	// Services

	MapServices(r)
	MapPartnerServices(r)
	MapLegacyServices(r)

	// Add the router action
	m.Action(r.Handle)
}

func MapServices(r martini.Router) {
	r.Get("/provider/:id", api.GetProvider)

	r.Get("/providers", api.GetProviders)

	r.Post("/project", api.AddProject)

	r.Get("/project/:id", api.GetProject)

	r.Get("/projects", api.GetProjects)

	// bit lazy
	r.Get("/convert_id/:id", api.DecodeProjectId)

	//		r.Put("/update_project/:id/:status", api.UpdateProject)
	r.Get("/update_project/:id/:status/:webserver", api.UpdateProject)

	r.Get("/update_job/:id/:status", api.UpdateJob)

	r.Post("/note", api.AddNote)

	r.Get("/notes/:id", api.GetNotes)

	r.Get("/trustcheck/:id", api.DoTrustcheck)

	r.Post("/buy_guarantee", api.BuyGuarantee)

	r.Post("/complete_payment", api.CompletePayment)

	r.Post("/admin_login", api.AdminLogin)

	r.Delete("/note/:project_id/:id", api.DeleteNote)

	r.Get("/", func() string {
		return "Welcome to TrustCloud API Services..."
	})
}

func MapPartnerServices(r martini.Router) {

	// Create a new Provider
	r.Post("/partner_provider", partner.AddProvider)

	// Get a Provider Card
	r.Get("/partner_provider/:id", partner.GetProviderCard)

	// Update a Provider
	r.Put("/partner_provider/:id", partner.UpdateProvider)
}

func MapLegacyServices(r martini.Router) {
	r.Get("/legacy_admin_login/:id/:password", legacy.AdminLogin)

	r.Get("/canned_responses", legacy.GetCannedResponses)

	r.Post("/add_canned_response", legacy.AddCannedResponse)

	r.Delete("/canned_response/:id", legacy.DeleteCannedResponse)

	r.Put("/canned_response/:id", legacy.UpdateCannedResponse)

	r.Get("/transaction/:id", legacy.GetTransaction)

	r.Post("/buy_trustcheck", legacy.BuyTrustcheck)

	r.Post("/complete_trustcheck_payment", legacy.CompleteTrustcheckPayment)

	r.Get("/legacy_transaction/:id", legacy.GetLegacyTransaction)

	// render trustcards
	r.Get("/render_card/:id/:size", legacy.RenderCard)
}

func CheckForwarder(c martini.Context, w http.ResponseWriter, r *http.Request) {

	if util.IsSet(util.Configuration.Environment.Http.Forwarder) {
		if r.Header.Get("X-Forwarded-Server") != util.Configuration.Environment.Http.Forwarder {
			util.ErrorLog.Println("Request from invalid origin")

			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

func MapEncoder(c martini.Context, w http.ResponseWriter, r *http.Request) {
	c.MapTo(util.JsonEncoder{}, (*util.Encoder)(nil))
	w.Header().Set("Content-Type", "application/json")
}

func main() {
	if !util.IsSet(util.Configuration.Environment.Http.Forwarder) {
		go func() {
			// Listen on https (in DEV)
			if err := http.ListenAndServeTLS(":4735", "cert/cert.pem", "cert/key.pem", m); err != nil {
				util.ErrorLog.Println("Exited with error... ", err)
			}
		}()
	}

	port := util.Configuration.General.Env.ListenPort

	util.InfoLog.Println("Started TrustCloud API Services Server...! [listening on port ", port, "]\n")

	error := http.ListenAndServe(":"+port, m)

	util.ErrorLog.Println("Exited with error... ", error)
}
