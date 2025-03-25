package track

import "github.com/go-chi/chi/v5"

type Router struct {
	router *chi.Mux
}

func NewRouter(router *chi.Mux) *Router {
	return &Router{router}
}

func (r *Router) Routes() {
	r.router.Route("/tracks", func(r chi.Router) {
		//r.Get("/", Handler.Find)
	})
}
