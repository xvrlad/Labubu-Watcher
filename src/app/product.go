package main

type Product struct {
	Handle    string `json:"handle"`
	Variants  []Variant `json:"variants"`
}