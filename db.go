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
}

type DB struct {
	productDB ProductDB
	mutex     sync.RWMutex
}

func LoadDB() (*DB, error) {
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
	return &DB{productDB: productDb}, nil
}

func (d *DB) save() error {
	file, err := os.Create("products.json")
	if err != nil {
		return err
	}
	err = json.NewEncoder(file).Encode(d.productDB)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func (d *DB) GetProducts() map[string]Product {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return maps.Clone(d.productDB.Products)
}

func (d *DB) AddProduct(product Product) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.productDB.Counter++
	product.ID = strconv.Itoa(d.productDB.Counter)
	d.productDB.Products[product.ID] = product
	d.save()
}

func (d *DB) GetProduct(id string) (Product, error) {
	d.mutex.RLock()
	product, ok := d.productDB.Products[id]
	d.mutex.RUnlock()
	if !ok {
		return Product{}, errors.New("product not found")
	}
	return product, nil
}

func (d *DB) SaveProduct(product Product) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.productDB.Products[product.ID] = product
	d.save()
}

func (d *DB) RemoveProduct(product Product) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	delete(d.productDB.Products, product.ID)
	d.save()
}
