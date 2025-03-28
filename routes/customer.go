package routes

import (
    "database/sql"
    "encoding/json"
    "net/http"
)

type Customer struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func GetCustomers(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        rows, _ := db.Query("SELECT id, name, email FROM customers")
        defer rows.Close()

        var customers []Customer
        for rows.Next() {
            var c Customer
            rows.Scan(&c.ID, &c.Name, &c.Email)
            customers = append(customers, c)
        }
        json.NewEncoder(w).Encode(customers)
    }
}
