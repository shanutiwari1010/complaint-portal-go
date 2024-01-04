package models

import (
	"goportal/utils"
	// "sync"
)

// var (
// 	userIDCounter      int
// 	complaintIDCounter int
// 	idMutex            sync.Mutex
// )

type User struct {
	ID           string
	SecretCode   string
	Name         string
	EmailAddress string
	Complaints   []Complaint
}

type Complaint struct {
	ID             string
	Title          string
	Summary        string
	UserSecretCode string
	Severity       int
	Resolved       bool
}

func NewUser(name, emailAddress string) *User {
	return &User{
		ID:           utils.GenerateID(),
		SecretCode:   utils.GenerateSecretCode(),
		Name:         name,
		EmailAddress: emailAddress,
		Complaints:   make([]Complaint, 0),
	}
}

func NewComplaint(title, summary string, severity int) *Complaint {
	return &Complaint{
		ID:       utils.GenerateID(),
		Title:    title,
		Summary:  summary,
		Severity: severity,
	}
}
