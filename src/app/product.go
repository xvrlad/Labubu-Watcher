package main

type Product struct {
	Handle    string `json:"handle"`;
	Variants  []struct {
		Available bool `json:"available"`;
	} `json:"variants"`;
}