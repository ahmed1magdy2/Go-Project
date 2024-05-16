package main

import (
	"bufio"
	//"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	//"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	serverPort = ":9040"
	dbUser     = "root"
	dbPass     = "root"
	dbName     = "DBGO"
	dbHost     = "localhost:3307"
)

var db *sql.DB

type datax struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func main() {
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, dbHost, dbName))
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	go startServer()

	startCLI()
}

func startServer() {
	http.HandleFunc("/fetch", fetchHandler)
	http.HandleFunc("/insert", insertHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/delete", deleteHandler)

	fmt.Println("Server is listening on port", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, nil))
}

func startCLI() {
	// Initialize scanner to read user input
	scanner := bufio.NewScanner(os.Stdin)

	// Main loop for CLI
	for {
		fmt.Println("Choose operation:")
		fmt.Println("1. Insert data")
		fmt.Println("2. Update data")
		fmt.Println("3. Delete data")
		fmt.Println("4. Show All Data")
		fmt.Println("5. Exit")

		// Read user input
		fmt.Print("Enter choice: ")
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			insertData(scanner)
		case "2":
			updateData(scanner)
		case "3":
			deleteData(scanner)
		case "4":
			fetchDataFromCLI()
		case "5":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice.")
		}
	}
}

func fetchDataFromCLI() {
	fmt.Println("Fetching all data from the server...")
	data, err := fetchDataFromDatabase()
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}

	fmt.Println("Data received from server: \n\n")
	for _, item := range data {
		fmt.Printf("ID: %d, Name: %s, Quantity: %d, Price: %.2f\n", item.ID, item.Name, item.Quantity, item.Price)
	}
	fmt.Println("\n\n")
}

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	data, err := fetchDataFromDatabase()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func fetchDataFromDatabase() ([]datax, error) {
	rows, err := db.Query("SELECT ID, Name, Quantity, Price FROM medicines")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []datax
	for rows.Next() {
		var d datax
		if err := rows.Scan(&d.ID, &d.Name, &d.Quantity, &d.Price); err != nil {
			return nil, err
		}
		data = append(data, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	var newData datax

	if err := json.NewDecoder(r.Body).Decode(&newData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := insertDataIntoDatabase(newData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Data inserted successfully")
}

func insertDataIntoDatabase(data datax) error {
	_, err := db.Exec("INSERT INTO medicines(Name, Quantity, Price) VALUES (?, ?, ?)", data.Name, data.Quantity, data.Price)
	return err
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var data datax

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := updateDataInDatabase(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Data updated successfully")
}

func updateDataInDatabase(data datax) error {
	_, err := db.Exec("UPDATE medicines SET Name = ?, Quantity = ?, Price = ? WHERE ID = ?", data.Name, data.Quantity, data.Price, data.ID)
	return err
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := deleteMedicineByID(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Medicine with ID %d deleted successfully", id)
}

func deleteMedicineByID(id int) error {
	_, err := db.Exec("DELETE FROM medicines WHERE ID = ?", id)
	return err
}

func insertData(scanner *bufio.Scanner) {
	// Read data from user
	fmt.Println("Enter data to insert:")
	data := readDataFromUser(scanner)

	// Insert data into database
	if err := insertDataIntoDatabase(data); err != nil {
		fmt.Println("Error inserting data:", err)
		return
	}

	fmt.Println("Data inserted successfully.")
}

func updateData(scanner *bufio.Scanner) {
	// Read ID from user
	fmt.Println("Enter ID of data to update:")
	id := readIntFromUser(scanner)

	// Read new data from user
	fmt.Println("Enter new data:")
	newData := readDataFromUser(scanner)
	newData.ID = id

	// Update data in database
	if err := updateDataInDatabase(newData); err != nil {
		fmt.Println("Error updating data:", err)
		return
	}

	fmt.Println("Data updated successfully.")
}

func deleteData(scanner *bufio.Scanner) {
	// Read ID from user
	fmt.Println("Enter ID of data to delete:")
	id := readIntFromUser(scanner)

	// Delete data from database
	if err := deleteMedicineByID(id); err != nil {
		fmt.Println("Error deleting data:", err)
		return
	}

	fmt.Println("Data deleted successfully.")
}

func readDataFromUser(scanner *bufio.Scanner) datax {
	fmt.Print("Name: ")
	scanner.Scan()
	name := scanner.Text()

	fmt.Print("Quantity: ")
	scanner.Scan()
	quantity, _ := strconv.Atoi(scanner.Text())

	fmt.Print("Price: ")
	scanner.Scan()
	price, _ := strconv.ParseFloat(scanner.Text(), 64)

	return datax{
		Name:     name,
		Quantity: quantity,
		Price:    price,
	}
}

func readIntFromUser(scanner *bufio.Scanner) int {
	scanner.Scan()
	id, _ := strconv.Atoi(scanner.Text())
	return id
}
