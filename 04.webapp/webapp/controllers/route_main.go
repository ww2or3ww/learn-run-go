package controllers

import (
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("show hello")
	generateHTML(w, "world", "layout", "public_navbar", "hello")
}

func world(w http.ResponseWriter, r *http.Request) {
	log.Println("show world")
	generateHTML(w, "!!!", "layout", "public_navbar", "world")
}
