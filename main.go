package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Produk struct {
	ID    int	`json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

type Kategori struct {
	ID   int    `json:"id"`
	Nama string `json:"nama"`
	Deskripsi string `json:"deskripsi"`
}


var produk = []Produk{
	{ID: 1, Nama: "Indomie Rebus", Harga: 3500, Stok: 10},
	{ID: 2, Nama: "Kecap Bango", Harga: 20000, Stok: 15},
	{ID: 3, Nama: "Susu Ultra", Harga: 15000, Stok: 8},
}

var kategori = []Kategori{
	{ID: 1, Nama: "Makanan", Deskripsi: "Produk makanan siap saji dan bahan makanan"},
	{ID: 2, Nama: "Minuman", Deskripsi: "Berbagai jenis minuman segar dan sehat"},
	{ID: 3, Nama: "Bumbu Dapur", Deskripsi: "Bumbu-bumbu dapur untuk memasak"},
}

func GetprodukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(p)
			return
		}
	}
	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

// PUT /api/produk/{id}
func updateProdukByID(w http.ResponseWriter, r *http.Request) {
	// GET id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	// Ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	// GET data dari request
	var updatedProduk Produk
	if err := json.NewDecoder(r.Body).Decode(&updatedProduk); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	// Loop produk, Cari ID yang sesuai request
	for i, p := range produk {
		if p.ID == id {
			updatedProduk.ID = id
			produk[i] = updatedProduk
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map [string]string{"message": "Produk berhasil ditambahkan"})
			return
		}
	}
	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

func deleteProduk(w http.ResponseWriter, r *http.Request) {
	// GET id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
		
	// ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}
	// Loop produk cari ID, dapat index yang dihapus
	for i, p := range produk {
		if p.ID == id {
			// bikin slice baru dengan data sebelu dan sesudah index
			produk = append(produk[:i], produk[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Produk berhasil dihapus"})
			return
		}
	}
	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}
	
func getKategori(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Kategori ID", http.StatusBadRequest)
		return
	}

	for _, p := range kategori {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(kategori)
			return
		}
	}
	http.Error(w, "Kategori tidak ditemukan", http.StatusNotFound)
}

// PUT /api/kategori/{id}
func updateKategori(w http.ResponseWriter, r *http.Request) {
	// GET id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
	// Ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Kategori ID", http.StatusBadRequest)
		return
	}

	// GET data dari request
	var updatedKategori Kategori
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Gagal membaca request body", http.StatusBadRequest)
		return
	}
	if len(bytes.TrimSpace(body)) == 0 {
		http.Error(w, "Request body kosong", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &updatedKategori); err != nil {
		// Fallback: support form data (application/x-www-form-urlencoded)
		r.Body = io.NopCloser(bytes.NewReader(body))
		if err := r.ParseForm(); err == nil {
			updatedKategori.Nama = r.FormValue("nama")
			updatedKategori.Deskripsi = r.FormValue("deskripsi")
			if updatedKategori.Nama != "" || updatedKategori.Deskripsi != "" {
				goto kategoriOK
			}
		}
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

kategoriOK:

	// Loop produk, Cari ID yang sesuai request
	for i, p := range kategori {
		if p.ID == id {
			updatedKategori.ID = id
			kategori[i] = updatedKategori
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(updatedKategori)
			return
		}
	}
	http.Error(w, "Kategori tidak ditemukan", http.StatusNotFound)
}

func deleteKategori(w http.ResponseWriter, r *http.Request) {
	// GET id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
		
	// ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Kategori ID", http.StatusBadRequest)
		return
	}
	// Loop produk cari ID, dapat index yang dihapus
	for i, p := range kategori {
		if p.ID == id {
			// bikin slice baru dengan data sebelu dan sesudah index
			kategori = append(kategori[:i], kategori[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Kategori berhasil dihapus"})
			return
		}
	}
	http.Error(w, "Kategori tidak ditemukan", http.StatusNotFound)
}


func main() {
	
	// GET localhost:8080/api/produk/{id}
	// PUT localhost:8080/api/produk/{id}
	// DELETE localhost:8080/api/produk/{id}
	

	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			GetprodukByID(w, r)
		} else if r.Method == "PUT" {
			updateProdukByID(w, r)
		} else if r.Method == "DELETE" {
			deleteProduk(w, r)
		}
		
	})

	//Get localhost:8080/api/produk
	//POST localhost:8080/api/produk
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(produk)
			return

		} else if r.Method == "POST" {
			//Baca dari Request Body
			//Masukan data ke variabel produk
			var produkBaru Produk
			if err := json.NewDecoder(r.Body).Decode(&produkBaru); err != nil {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
				return
			}

			//Masukan ke dalam variabel produk
			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)
			w.Header().Set("Content-Type", "application/json")	
			w.WriteHeader(http.StatusCreated) //201
			_ = json.NewEncoder(w).Encode(produkBaru)
			return
		}

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	// Get localhost:8080/api/produk
	// POST localhost:8080/api/produk
	http.HandleFunc("/api/kategori", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(kategori)
			return

		} else if r.Method == "POST" {
			//Baca dari Request Body
			//Masukan data ke variabel kategori
			var kategoriBaru Kategori
			if err := json.NewDecoder(r.Body).Decode(&kategoriBaru); err != nil {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
				return
			}

			//Masukan ke dalam variabel kategori
			kategoriBaru.ID = len(kategori) + 1
			kategori = append(kategori, kategoriBaru)
			w.Header().Set("Content-Type", "application/json")	
			w.WriteHeader(http.StatusCreated) //201
			_ = json.NewEncoder(w).Encode(kategoriBaru)
			return
		}
		
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}	)

	// GET localhost:8080/api/kategori/{id}
	// PUT localhost:8080/api/kategori/{id}
	// DELETE localhost:8080/api/kategori/{id}
	http.HandleFunc("/api/kategori/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getKategori(w, r)
		} else if r.Method == "PUT" {
			updateKategori(w, r)
		} else if r.Method == "DELETE" {
			deleteKategori(w, r)
		}

	})

	

	
	//Localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "message": "Service API is running smoothly"})
	})

	fmt.Println("Starting server on Localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server failed to start:", err)
	}
}
