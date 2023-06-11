package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Product struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from server"))
	})

	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			json, err := json.Marshal([]Product{
				{Id: "1", Name: "Book", Price: 5000},
				{Id: "2", Name: "Pencil", Price: 3000},
				{Id: "3", Name: "Pen", Price: 2500},
			})

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			w.Write(json)
			return
		} else if r.Method == http.MethodPost {
			data, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			var product Product
			if err := json.Unmarshal(data, &product); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			data, err = json.Marshal(product)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			w.Write(data)
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	})

	http.ListenAndServe(":8080", nil)
}
