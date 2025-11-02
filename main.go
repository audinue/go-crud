package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"text/template"
)

type Product struct {
	ID   string
	Name string
}

type ProductDB struct {
	Products map[string]Product
	Counter  int
}

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
	file, err := os.Open("products.json")
	var db ProductDB
	if err != nil {
		db = ProductDB{
			Products: map[string]Product{
				"1": {ID: "1", Name: "Apple"},
				"2": {ID: "2", Name: "Banana"},
				"3": {ID: "3", Name: "Cherry"},
			},
			Counter: 3,
		}
	} else {
		json.NewDecoder(file).Decode(&db)
		file.Close()
	}
	mutex := sync.RWMutex{}
	save := func() {
		file, err := os.Create("products.json")
		if err != nil {
			log.Fatal(err)
		}
		err = json.NewEncoder(file).Encode(db)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		mutex.RLock()
		defer mutex.RUnlock()
		render(w, "products", map[string]any{
			"Products": db.Products,
		})
	})
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			mutex.Lock()
			defer mutex.Unlock()
			db.Counter++
			id := strconv.Itoa(db.Counter)
			db.Products[id] = Product{ID: id, Name: r.FormValue("name")}
			save()
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
		mutex.RLock()
		product, ok := db.Products[id]
		mutex.RUnlock()
		if !ok {
			w.WriteHeader(404)
			return
		}
		if r.Method == "POST" {
			mutex.Lock()
			defer mutex.Unlock()
			product.Name = r.FormValue("name")
			db.Products[id] = product
			save()
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
		mutex.RLock()
		_, ok := db.Products[id]
		mutex.RUnlock()
		if !ok {
			w.WriteHeader(404)
			return
		}
		if r.Method == "POST" {
			mutex.Lock()
			defer mutex.Unlock()
			delete(db.Products, id)
			save()
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
