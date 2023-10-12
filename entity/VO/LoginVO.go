package VO

type LoginVO struct {
	Token   string `json:"token"`
	Role    string `json:"role"`
	Success bool   `json:"success"`
}
