package main

import (
	"log"
	"net/http"
)

func main() {
	db, err := LoadDB()
	if err != nil {
		log.Fatal(err)
	}
	template, err := NewTemplate()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		template.Render(w, "products", map[string]any{
			"Products": db.GetProducts(),
		})
	})
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			db.AddProduct(Product{Name: r.FormValue("name")})
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			template.Render(w, "products-form", nil)
		}
	})
	http.HandleFunc("/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			w.WriteHeader(400)
			return
		}
		product, err := db.GetProduct(id)
		if err != nil {
			w.WriteHeader(404)
			return
		}
		if r.Method == "POST" {
			product.Name = r.FormValue("name")
			db.SaveProduct(product)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			template.Render(w, "products-form", map[string]any{
				"Product": product,
			})
		}
	})
	http.HandleFunc("/{id}/remove", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			w.WriteHeader(400)
			return
		}
		product, err := db.GetProduct(id)
		if err != nil {
			w.WriteHeader(404)
			return
		}
		if r.Method == "POST" {
			db.RemoveProduct(product)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			template.Render(w, "products-confirm", nil)
		}
	})
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
