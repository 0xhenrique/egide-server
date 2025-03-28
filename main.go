package main

import (
    "log"

    "github.com/joho/godotenv"
    "github.com/0xhenrique/egide-server/auth"
    //"github.com/0xhenrique/egide-server/routes"
)

func main() {
    godotenv.Load()
    //db := connectDB()
    //defer db.Close()

    // Initialize authentication and start server
    auth.NewAuth()

    // Other routes
    //routes.InitRoutes(db)

    log.Println("Server running on: 8080")
}
