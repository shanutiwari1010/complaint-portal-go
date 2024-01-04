package handlers

import (
	"encoding/json"
	"fmt"
	"goportal/models"
	"goportal/utils"
	"net/http"
	"sync"
)

var (
	userMutex      sync.Mutex
	complaintMutex sync.Mutex
	users          = make(map[string]models.User)
	complaints     = make(map[string]models.Complaint)
)

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for /login route
	var request struct {
		SecretCode string `json:"secret_code"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.SecretCode == "" {
		http.Error(w, "Secret Code is required", http.StatusBadRequest)
		return
	}

	userMutex.Lock()
	user, found := users[request.SecretCode]
	userMutex.Unlock()

	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var newUser models.User

	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if newUser.Name == "" || newUser.EmailAddress == "" {
		http.Error(w, "Name and email address are required", http.StatusBadRequest)
		return
	}

	userMutex.Lock()
	defer userMutex.Unlock()

	for _, existingUser := range users {
		if existingUser.EmailAddress == newUser.EmailAddress {
			http.Error(w, "Email address is already registered", http.StatusConflict)
			return
		}
	}

	newUser.ID = utils.GenerateID()
	newUser.SecretCode = utils.GenerateSecretCode()

	users[newUser.SecretCode] = newUser

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}

func SubmitComplaintHandler(w http.ResponseWriter, r *http.Request) {
	var newComplaint models.Complaint

	err := json.NewDecoder(r.Body).Decode(&newComplaint)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if newComplaint.Title == "" || newComplaint.Summary == "" || newComplaint.Severity < 1 {
		http.Error(w, "Invalid complaint details. Title, summary, and severity are required.", http.StatusBadRequest)
		return
	}

	newComplaint.ID = utils.GenerateID()

	complaintMutex.Lock()
	defer complaintMutex.Unlock()

	complaints[newComplaint.ID] = newComplaint
}

func GetAllComplaintsForUserHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SecretCode string `json:"secret_code"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.SecretCode == "" {
		http.Error(w, "Secret code is required", http.StatusBadRequest)
		return
	}

	userMutex.Lock()
	user, found := users[request.SecretCode]
	userMutex.Unlock()

	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.Complaints)
}

func GetAllComplaintsForAdminHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for /getAllComplaintsForAdmin route
	isAdmin := true

	if !isAdmin {
		http.Error(w, "Unauthorized: Only admins can access this route", http.StatusUnauthorized)
		return
	}

	// Retrieve all complaints
	complaintMutex.Lock()
	allComplaints := make([]models.Complaint, 0, len(complaints))
	for _, complaint := range complaints {
		allComplaints = append(allComplaints, complaint)
	}
	complaintMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allComplaints)
}

func ViewComplaintHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for /viewComplaint route
	var request struct {
		SecretCode  string `json:"secret_code"`
		ComplaintID string `json:"complaint_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.SecretCode == "" || request.ComplaintID == "" {
		http.Error(w, "Secret code and complaint ID are required", http.StatusBadRequest)
		return
	}

	complaintMutex.Lock()
	defer complaintMutex.Unlock()

	complaint, found := complaints[request.ComplaintID]

	// Check if the complaint is found
	if !found {
		http.Error(w, "Complaint not found", http.StatusNotFound)
		return
	}

	if complaint.UserSecretCode != request.SecretCode {
		http.Error(w, "Unauthorized to view this complaint", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(complaint)
}

func ResolveComplaintHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var request struct {
		SecretCode  string `json:"secret_code"`
		ComplaintID string `json:"complaint_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if both secret code and complaint ID are provided
	if request.SecretCode == "" || request.ComplaintID == "" {
		http.Error(w, "Secret code and complaint ID are required", http.StatusBadRequest)
		return
	}

	// Lock the complaints map for concurrent access
	complaintMutex.Lock()
	defer complaintMutex.Unlock()

	// Look up the complaint in the complaints map
	complaint, found := complaints[request.ComplaintID]
	if !found {
		http.Error(w, "Complaint not found", http.StatusNotFound)
		return
	}

	// Check if the user has permission to resolve the complaint
	if complaint.UserSecretCode != request.SecretCode {
		http.Error(w, "Unauthorized to resolve this complaint", http.StatusUnauthorized)
		return
	}

	// Check if the complaint is already resolved
	if complaint.Resolved {
		http.Error(w, "Complaint is already resolved", http.StatusConflict)
		return
	}

	// Mark the complaint as resolved
	complaint.Resolved = true

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Complaint resolved successfully"))
}
