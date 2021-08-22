package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"log"
	"net/http"
	"database/sql" 
	"strconv"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql" 
)

type Product struct {
	Id int `json:"id"`
	Name string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductHandler struct {
	sync.Mutex
	db *sql.DB 
}

func newProductHandler(db *sql.DB) *ProductHandler { // Kinda like a constructor
	return &ProductHandler{
		db: db, 
	}
}

// Implements Handler Interface (see below)
func (ph *ProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ph.get(w, r)
	case "POST":
		ph.post(w, r)
	case "PUT", "PATCH":
		ph.put(w, r)
	case "DELETE":
		ph.delete(w, r)
	default:
		respondWithJSON(w, http.StatusMethodNotAllowed, "invalid method") // http has built in status code variables. StatusMethodNotAllowed is 405
	}
}

func (ph *ProductHandler) get(w http.ResponseWriter, r *http.Request) {
	defer ph.Unlock()
	ph.Lock()
	id, err := idFromUrl(r)
	if err != nil {
		// Return all products 
		results, err := ph.db.Query("select * from products")
		if err != nil {
			log.Fatal("Error when fetching products from table", err.Error())
		}
		defer results.Close() 

		var products []Product 
		for results.Next() {	
			var product Product
			err = results.Scan(&product.Id, &product.Name, &product.Price) 
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return 
			}
			products = append(products, product)
		}
		respondWithJSON(w, http.StatusOK, products)
		return 
	}
	
	// returning product with id 
	var product Product 
	qry := fmt.Sprintf("Select * from products where id = %v", id)
	err = ph.db.QueryRow(qry).Scan(&product.Id, &product.Name, &product.Price) 
	// Unmarshall supports &product but not Scan WHY???
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return 
	}
	respondWithJSON(w, http.StatusOK, product)

}

func (ph *ProductHandler) post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()             
	body, err := ioutil.ReadAll(r.Body)
	// body here is a []byte

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error()) 
		return // always return after responding.
	}
	contentType := r.Header.Get("content-type")
	if contentType != "application/json" {
		respondWithError(w, http.StatusUnsupportedMediaType, "content type 'application/json' required")
		return
	}
	var product Product
	err = json.Unmarshal(body, &product) // converts a []byte json string to a Product
	
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer ph.Unlock() // defered functions go in a STACK. execution is Last-In-First-Out
	ph.Lock()

	// Adding product to sql table 
	create_product, err := ph.db.Prepare("INSERT INTO products (name, price) VALUES (?, ?)")
	if err != nil {
		fmt.Println("Errrr", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return 
	} 
	_, err = create_product.Exec(product.Name, product.Price)
	if err != nil {
		fmt.Println("Error::", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return 
	} 
	fmt.Println("Product created", product)
	respondWithJSON(w, http.StatusCreated, product)
}

func (ph *ProductHandler) put(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id, err := idFromUrl(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	body, err := ioutil.ReadAll(r.Body) // Reads all from a Reader (r.Body implements Reader interface)
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return                                                          
	}
	contentType := r.Header.Get("content-type")
	if contentType != "application/json" {
		respondWithError(w, http.StatusUnsupportedMediaType, "content type 'application/json' required")
		return
	}
	var product Product
	err = json.Unmarshal(body, &product)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("Updating product", product)

	// Updating product 
	defer ph.Unlock()
	ph.Lock()

	// Check if product exists
	var origP Product 
	qry := fmt.Sprintf("Select * from products where id = %v", id)
	err = ph.db.QueryRow(qry).Scan(&origP.Id, &origP.Name, &origP.Price) 
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return 
	}

	if product.Name != "" {
		origP.Name = product.Name
	}
	if product.Price != 0.0 {
		origP.Price = product.Price
	}

	// Updating Product
	update_product, err := ph.db.Prepare(`UPDATE products SET name = ?, price = ? WHERE id = ?`)
	if err != nil {
		fmt.Println("Errrr", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return 
	} 
	
	_, err = update_product.Exec(origP.Name, origP.Price, origP.Id)
	if err != nil {
		fmt.Println("Error::", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return 
	} 
	
	fmt.Println("Product updated", origP)
	respondWithJSON(w, http.StatusCreated, origP)
}

func (ph *ProductHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := idFromUrl(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "not found")
		return
	}

	defer ph.Unlock()
	ph.Lock()

	// Check if product exists
	var origP Product 
	qry := fmt.Sprintf("Select * from products where id = %v", id)
	err = ph.db.QueryRow(qry).Scan(&origP.Id, &origP.Name, &origP.Price) 
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return 
	}

	// deleting product 
	delete_product, err := ph.db.Prepare(`DELETE FROM products WHERE id = ?`)
	if err != nil {
		fmt.Println("Errrr", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return 
	} 
	
	_, err = delete_product.Exec(origP.Id)
	if err != nil {
		fmt.Println("Error::", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return 
	} 
	
	fmt.Println("Product deleted", origP)
	respondWithJSON(w, http.StatusNoContent, "")   // It is convention to not return anything for DELETE request
}

func respondWithJSON(w http.ResponseWriter, code int, data interface{}) { // empty interface is like Any in JS, try avoid using this.
	response, err := json.Marshal(data) // returns a []byte (json in string)
	if err != nil {
		fmt.Println("Error:", err)
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code) // writes status code (eg. 200 OK)
	w.Write(response)
	// As soon as w.Write() is executed, the Server will send the response
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

// IDs here are simply indexes for simpliciy
func idFromUrl(r *http.Request) (int, error) {
	parts := strings.Split(r.URL.String(), "/") // r.URL.String() returns a string representation of the URL
	if len(parts) != 3 {
		return -1, errors.New("not found") // errors.New returns an error with given message (obviously)
	}
	id, err := strconv.Atoi(parts[len(parts)-1]) // strconv.Atoi parses a string to in int
	if err != nil {
		return -1, errors.New("not found")
	}
	return id, nil

}

func initDB() *sql.DB {
	DB_URL := os.Getenv("DB_URL") 
	//"root:admin@tcp(127.0.0.1:3306)/test"

	URLS := strings.Split(DB_URL, "/")
	
	CONNECTION_URL := URLS[0]
	DATABASE_NAME := URLS[1]

	fmt.Println("Connection:", CONNECTION_URL)
	fmt.Println("Database:", DATABASE_NAME)

	// Add sql driver by: "go get github.com/go-sql-driver/mysql"
	fmt.Println("Drivers:", sql.Drivers())

	db, err := sql.Open("mysql", CONNECTION_URL + "/")
	if err != nil {
		log.Fatal("Unable to open connection to DB", err.Error())
	} else {
		fmt.Println("Connected To DB...")
	}

	_,err = db.Exec("USE " + DATABASE_NAME)
	if err != nil {
		fmt.Println("DATABASE NOT EXISTS", DATABASE_NAME, err.Error())
		_,err = db.Exec("CREATE DATABASE "+ DATABASE_NAME)
		if err != nil {
			log.Fatal("Error: Creating Database", err.Error())
		} else {
			fmt.Println("Database Created Successfully...")
			_,err = db.Exec("USE " + DATABASE_NAME)
			if err != nil {
				log.Fatal("Error: Using Newly Created Database...")
			}
		}
	} else {
		fmt.Println("Database found...")
	}
	return db 
}

func checkTable(table string, db *sql.DB) error {
	results, err := db.Query("select * from products")
	if err != nil {
		fmt.Println("Table", table, "does not exists")
		create_table, err := db.Prepare("CREATE TABLE products(id int NOT NULL AUTO_INCREMENT, name varchar(30), price int, PRIMARY KEY (id))")
		if err != nil {
			fmt.Println("Error: Creating Table Statement")
			return err 
		} else {
			_,err := create_table.Exec() 
			if err != nil {
				fmt.Println("Error: Creating Table:")
				return err 
			} else {
				fmt.Println("Table Created successfully")
				return nil 
			}
		}
	} else {
		fmt.Println("Table", table, "Found...")
		results.Close()
		return nil  
	}
} 

func RunServer() {
	db := initDB()
	err := checkTable("products", db)
	if err != nil {
		log.Fatal("Error: Table products:", err.Error())
	}

	defer db.Close() // TO make sure DB connection is closed properly

	// var db *sql.DB
	port := ":8080"
	ph := newProductHandler(db)
	http.Handle("/products", ph)
	http.Handle("/products/", ph) // "/products" and "/products/" are handled differently

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World") // Fprintf takes a Writer and outputs to it.
	})

	fmt.Println("Listening on Port", port)
	launchErr := http.ListenAndServe(port, nil)
	log.Fatal(launchErr)
}
