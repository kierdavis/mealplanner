package mphandlers

import (
	"net/http"
)

func AttachHandlers(m *http.ServeMux) {
	m.HandleFunc("/", handleHome)
	m.HandleFunc("/meals/new", handleAddMeal)
}
