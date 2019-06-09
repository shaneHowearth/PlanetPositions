package main

import "github.com/go-chi/chi"

func planetRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/Sunrise/{long}/{lat}/{date}", GetSunrise)
	// router.Delete("/{todoID}", DeleteTodo)
	// router.Post("/", CreateTodo)
	// router.Get("/", GetAllTodos)
	return router
}

func GetSunrise() {}
