package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const PORT = "8080"
const BASE_URL = "http://localhost:" + PORT

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

	http.HandleFunc("/products/client", func(w http.ResponseWriter, r *http.Request) {
		request, err := http.NewRequest(r.Method, BASE_URL+"/products", r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(response.StatusCode)
		w.Write(data)
		return
	})

	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		fmt.Println(err.Error())
	}
}
