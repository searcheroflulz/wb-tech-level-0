package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"html/template"
	"net/http"
	"wb-tech-level-0/internal/storage/cache"
)

type PageData struct {
	Title string
}

func NewIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("internal/templates/index.html")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
			return
		}

		data := PageData{Title: "Order Lookup"}
		err = tmpl.Execute(w, data)
		if err != nil {
			return
		}
	}
}

func NewGetOrder(cache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := chi.URLParam(r, "orderID")

		result, ok := cache.GetOrder(params)

		if !ok {
			errMsg := "Order ID not found: " + params
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": errMsg})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
