package main

type LEDValues struct {
	Display    int      `json:"display"`
	Brightness int      `json:"brightness"`
	Colors     []string `json:"colors"`
}
