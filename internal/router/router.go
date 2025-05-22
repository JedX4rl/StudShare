package router

import (
	"StudShare/internal/router/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *handler.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.URLFormat)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(h.TokenManager))
			r.Post("/logout", h.Logout)
		})
	}) //Готово

	r.Route("/users", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(h.TokenManager))
			r.Get("/", h.GetProfile)
			r.Put("/update", h.UpdateProfile)
		})
		r.Get("/profile", h.GetProfileByID)
	}) //Готово

	r.Route("/listings", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(h.TokenManager))
			r.Post("/", h.CreateListing)   //good
			r.Put("/", h.UpdateListing)    //good
			r.Delete("/", h.DeleteListing) //good
		})

		r.Get("/all", h.GetAllListings)   //good
		r.Get("/", h.GetListingByID)      //good
		r.Get("/near", h.GetListingsNear) //good
	})

	r.Route("/reviews", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(h.TokenManager))
			r.Post("/", h.AddReview)
			r.Delete("/", h.DeleteReview)
		})

		r.Get("/user", h.GetReviewsForUser)
	})

	r.Route("/files", func(r chi.Router) {
		r.Use(AuthMiddleware(h.TokenManager))
		r.Post("/", h.UploadFile)
		r.Delete("/{key}", h.DeleteFile)
	})

	r.Route("/drafts", func(r chi.Router) {
		r.Use(AuthMiddleware(h.TokenManager))
		r.Post("/", h.Create)
		r.Get("/all", h.GetAll)
		r.Get("/", h.GetByID)
		r.Delete("/", h.DeleteByID)
	})

	return r
}
