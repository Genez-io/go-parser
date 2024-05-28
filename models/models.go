package models

type Class struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type Method struct {
	Name      string `json:"name"`
	Comment   string `json:"comment"`
	ClassName string `json:"className"`
}