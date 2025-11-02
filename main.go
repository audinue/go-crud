package main

import (
	"embed"
	"log"
	"net/http"
	"text/template"
)

//go:embed templates
var fs embed.FS

func main() {
	templates, err := template.ParseFS(fs, "templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
	render := func(w http.ResponseWriter, name string, data any) {
		err := templates.ExecuteTemplate(w, name+".html", data)
		if err != nil {
			log.Println(err)
		}
	}
	db, err := LoadDB()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		render(w, "products", map[string]any{
			"Products": db.GetProducts(),
		})
	})
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			db.AddProduct(Product{Name: r.FormValue("name")})
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			render(w, "products-form", nil)
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
			render(w, "products-form", map[string]any{
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
			render(w, "products-confirm", nil)
		}
	})
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
