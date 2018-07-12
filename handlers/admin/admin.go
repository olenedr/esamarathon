package admin

import (
	"net/http"

	"github.com/dannyvankooten/grender"

	"github.com/esamarathon/website/config"
	"github.com/esamarathon/website/models/social"
	"github.com/esamarathon/website/models/user"
	"github.com/esamarathon/website/viewmodels"

	"github.com/esamarathon/website/middleware"
	"github.com/gorilla/mux"
)

// AdminRoutes adds the admin routes to the router
func AdminRoutes(base string, router *mux.Router) {
	requireAuth := middleware.AuthMiddleware
	router.HandleFunc(base, CSP(requireAuth(adminIndex))).Methods("GET", "POST")
	router.HandleFunc(base+"/toggle", CSP(requireAuth(toggleLivemode))).Methods("GET")
	router.HandleFunc(base+"/front", CSP(requireAuth(updateFront))).Methods("POST")
	router.HandleFunc(base+"/user", CSP(requireAuth(userIndex))).Methods("GET")
	router.HandleFunc(base+"/user", CSP(requireAuth(userStore))).Methods("POST")
	router.HandleFunc(base+"/user/{id}/delete", CSP(requireAuth(deleteUser))).Methods("GET")

	router.HandleFunc(base+"/article", CSP(requireAuth(articleIndex))).Methods("GET")
	router.HandleFunc(base+"/article/create", CSP(requireAuth(articleCreate))).Methods("GET")
	router.HandleFunc(base+"/article/create", CSP(requireAuth(articleStore))).Methods("POST")
	router.HandleFunc(base+"/article/{id}", CSP(requireAuth(articleEdit))).Methods("GET")
	router.HandleFunc(base+"/article/{id}", CSP(requireAuth(articleUpdate))).Methods("POST")
	router.HandleFunc(base+"/article/{id}/delete", CSP(requireAuth(articleDelete))).Methods("GET")

	router.HandleFunc(base+"/page", CSP(requireAuth(pageIndex))).Methods("GET")
	router.HandleFunc(base+"/page/create", CSP(requireAuth(pageCreate))).Methods("GET")
	router.HandleFunc(base+"/page/create", CSP(requireAuth(pageStore))).Methods("POST")
	router.HandleFunc(base+"/page/{id}", CSP(requireAuth(pageEdit))).Methods("GET")
	router.HandleFunc(base+"/page/{id}", CSP(requireAuth(pageUpdate))).Methods("POST")
	router.HandleFunc(base+"/page/{id}/delete", CSP(requireAuth(pageDelete))).Methods("GET")

	router.HandleFunc(base+"/menu", CSP(requireAuth(menuIndex))).Methods("GET")
	router.HandleFunc(base+"/menu/{id}", CSP(requireAuth(menuUpdate))).Methods("POST")

	router.HandleFunc(base+"/schedule", CSP(requireAuth(scheduleIndex))).Methods("GET")
	router.HandleFunc(base+"/schedule/create", CSP(requireAuth(scheduleCreate))).Methods("POST")
	router.HandleFunc(base+"/schedule/{id}", CSP(requireAuth(scheduleUpdate))).Methods("POST")
	router.HandleFunc(base+"/schedule/{id}/delete", CSP(requireAuth(scheduleDelete))).Methods("POST")

	router.HandleFunc(base+"/social/{id}", CSP(requireAuth(socialUpdate))).Methods("POST")
}

// Initiates a renderer for the admin views
var adminRenderer = grender.New(grender.Options{
	TemplatesGlob: "templates_admin/*.html",
	PartialsGlob:  "templates_admin/partials/*.html",
})

/*
*	Admin Index routes
 */
func adminIndex(w http.ResponseWriter, r *http.Request) {
	if !social.IsStored() {
		s := social.Get()
		err := s.Insert()
		if len(err) > 0 {
			user.SetFlashMessage(w, r, "alert", "Couldn't get data from DB. There might be connection issues or the table might not exist!")
		}
	}
	adminRenderer.HTML(w, http.StatusOK, "index.html", viewmodels.AdminIndex(w, r))
}

// Toggles the stream on the frontpage
func toggleLivemode(w http.ResponseWriter, r *http.Request) {
	config.ToggleLiveMode()
	http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
}

// Update the text on the front row based on the input data
func updateFront(w http.ResponseWriter, r *http.Request) {
	// Parse input data
	r.ParseForm()
	title := r.Form.Get("title")
	body := r.Form.Get("body")

	// If title or body is empty
	if title == "" || body == "" {
		// Set flash message and redirect
		user.SetFlashMessage(w, r, "alert", "Not enough input data, please fill inn Title and Content")
		http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
		return
	}

	// Update frontpage with new input
	viewmodels.UpdateFrontpage(title, body)

	// Set flaash and redirect back
	user.SetFlashMessage(w, r, "success", "The frontpage has been updated!")
	http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
}

func socialUpdate(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	si, err := social.Find(id)
	if err != nil {
		user.SetFlashMessage(w, r, "alert", "Couldn't find the social item you wanted to update")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	si.Title = r.Form.Get("title")
	si.Link = r.Form.Get("link")
	si.Image = r.Form.Get("image")
	si.ImageAlt = r.Form.Get("imagealt")
	if r.Form.Get("new_tab") == "true" {
		si.NewTab = true
	} else {
		si.NewTab = false
	}
	err = si.Update()
	if err != nil {
		user.SetFlashMessage(w, r, "alert", "Something went wrong while trying to update")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	user.SetFlashMessage(w, r, "success", "The social link was updated")
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
