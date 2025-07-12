package models

type Product struct {
	Handle    string `json:"handle"`
	Variants  []Variant `json:"variants"`
}