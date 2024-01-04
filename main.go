package main

import (
	"fmt"
	"goportal/handlers"
	"net/http"
)

func main() {

	http.HandleFunc("/hello", handlers.Hello)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/submitComplaint", handlers.SubmitComplaintHandler)
	http.HandleFunc("/getAllComplaintsForUser", handlers.GetAllComplaintsForUserHandler)
	http.HandleFunc("/getAllComplaintsForAdmin", handlers.GetAllComplaintsForAdminHandler)
	http.HandleFunc("/viewComplaint", handlers.ViewComplaintHandler)
	http.HandleFunc("/resolveComplaint", handlers.ResolveComplaintHandler)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}