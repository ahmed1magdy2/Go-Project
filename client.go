package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

const serverURL = "http://127.0.0.1:9040"

type datax struct {
	ID       int
	Name     string
	Quantity int
	Price    float64
}

func main() {
	// Create a client with a persistent connection
	client := &http.Client{}

	// Initialize scanner to read user input
	scanner := bufio.NewScanner(os.Stdin)

	// Main loop for CLI
	for {
		fmt.Println("Choose operation:")
		fmt.Println("1. Insert data")
		fmt.Println("2. Update data")
		fmt.Println("3. Delete data")
		fmt.Println("4. Show All data")
		fmt.Println("5. Exit")

		// Read user input
		fmt.Print("Enter choice: ")
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			insertData(client, scanner)
		case "2":
			updateData(client, scanner)
		case "3":
			deleteData(client, scanner)
		case "4":
			fetchData(client)
		case "5":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice.")
		}
	}
}

func insertData(client *http.Client, scanner *bufio.Scanner) {
	// Read data from user
	fmt.Println("Enter data to insert:")
	data := readDataFromUser(scanner)

	// Send insert request
	err := sendDataToServer(client, http.MethodPost, "/insert", data)
	if err != nil {
		fmt.Println("Error inserting data:", err)
		return
	}

	fmt.Println("Data inserted successfully.")
}

func updateData(client *http.Client, scanner *bufio.Scanner) {
	// Read data from user
	fmt.Println("Enter ID of data to update:")
	id := readIntFromUser(scanner)

	fmt.Println("Enter new data:")
	newData := readDataFromUser(scanner)
	newData.ID = id

	// Send update request
	err := sendDataToServer(client, http.MethodPut, "/update", newData)
	if err != nil {
		fmt.Println("Error updating data:", err)
	}
}

// Function to delete data
func deleteData(client *http.Client, scanner *bufio.Scanner) {
	// Read ID from user
	fmt.Println("Enter ID of data to delete:")
	id := readIntFromUser(scanner)

	// Validate ID
	if id <= 0 {
		fmt.Println("Invalid ID. Please enter a positive integer.")
		return
	}

	// Send DELETE request
	endpoint := fmt.Sprintf("/delete?id=%d", id)
	err := sendDataToServer(client, http.MethodDelete, endpoint, datax{})
	if err != nil {
		fmt.Println("Error deleting data:", err)
		return
	}

	fmt.Println("Data deleted successfully.")
}

func fetchData(client *http.Client) {
	// Fetch data from the server
	data, err := fetchDataFromServer("/fetch")
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}

	// Display data received from the server
	fmt.Println("Data received from server: \n\n")
	for _, item := range data {
		fmt.Printf("ID: %d, Name: %s, Quantity: %d, Price: %.2f\n", item.ID, item.Name, item.Quantity, item.Price)
	}
	fmt.Println("\n\n")
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

func sendDataToServer(client *http.Client, method, endpoint string, data datax) error {
	url := serverURL + endpoint

	// Marshal data into JSON format
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Create a new HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-200 status code: %s", resp.Status)
	}

	return nil
}

func fetchDataFromServer(endpoint string) ([]datax, error) {
	url := serverURL + endpoint

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode JSON response into datax slice
	var data []datax
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
