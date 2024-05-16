# Medicine Management

This project is a Go application that provides a command-line interface (CLI) and an HTTP server to manage a medicines database. It allows users to insert, update, delete, and fetch data about medicines. The application uses a MySQL database to store the data and communicates via JSON over HTTP.

## Features

- Insert new medicine data
- Update existing medicine data
- Delete medicine data
- Fetch all medicine data
- CLI for interacting with the application
- HTTP server for remote operations

## Prerequisites

- Go 1.16 or higher
- MySQL

## Setup

### MySQL Setup

1. Install MySQL and start it.
2. Create a database named `DBGO`:
   ```sql
   CREATE DATABASE DBGO;
3. Create a table named medicines in the DBGO database:
    ```sql
    CREATE TABLE medicines (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    Quantity INT NOT NULL,
    Price FLOAT NOT NULL
    );

## Project Setup
1. Clone the repository:
```
git clone https://github.com/ahmed1magdy2/Go-Project.git
cd Go-Project
```
2. Install Go dependencies:
```
go mod init example
go mod tidy
go get github.com/go-sql-driver/mysql
```
3. Build and run the server:
```
go build -o server server.go
./master
```
4. In a new terminal, build and run the CLI:
```
go build -o client client.go
./client
```
## Usage
### server
The server will run on http://127.0.0.1:9040 and provides the following endpoints:

- POST /insert - Insert new medicine data.
- PUT /update - Update existing medicine data.
- DELETE /delete?id={id} - Delete medicine data by ID.
- GET /fetch - Fetch all medicine data.

### Client
The Client provides an interactive interface for performing CRUD operations. When prompted, enter the corresponding number to choose an operation:

1. Insert data - Insert new medicine data.
2. Update data - Update existing medicine data.
3. Delete data - Delete medicine data by ID.
4. Show All data - Fetch and display all medicine data.
5. Exit - Exit the CLI.

## Project Structure
- `server.go` - Contains the server implementation for listening to many clients and database interactions.
- `client.go` - Contains the CLI implementation for interacting with the server.
