package main

import (
	"encoding/json"
	"errors"
	"maps"
	"os"
	"strconv"
	"sync"
)

type Product struct {
	ID   string
	Name string
}

type ProductDB struct {
	Products map[string]Product
	Counter  int
	mutex    sync.RWMutex
}

func LoadProductDB() (*ProductDB, error) {
	var productDb ProductDB
	file, err := os.Open("products.json")
	if err != nil {
		productDb = ProductDB{
			Products: map[string]Product{
				"1": {ID: "1", Name: "Apple"},
				"2": {ID: "2", Name: "Banana"},
				"3": {ID: "3", Name: "Cherry"},
			},
			Counter: 3,
		}
	} else {
		err = json.NewDecoder(file).Decode(&productDb)
		if err != nil {
			return nil, err
		}
		file.Close()
	}
	return &productDb, nil
}

func (d *ProductDB) save() error {
	file, err := os.Create("products.json")
	if err != nil {
		return err
	}
	err = json.NewEncoder(file).Encode(d)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func (d *ProductDB) All() map[string]Product {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return maps.Clone(d.Products)
}

func (d *ProductDB) Add(product Product) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.Counter++
	product.ID = strconv.Itoa(d.Counter)
	d.Products[product.ID] = product
	return d.save()
}

func (d *ProductDB) Get(id string) (Product, error) {
	d.mutex.RLock()
	product, ok := d.Products[id]
	d.mutex.RUnlock()
	if !ok {
		return Product{}, errors.New("product not found")
	}
	return product, nil
}

func (d *ProductDB) Edit(product Product) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.Products[product.ID] = product
	return d.save()
}

func (d *ProductDB) Remove(product Product) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	delete(d.Products, product.ID)
	return d.save()
}
