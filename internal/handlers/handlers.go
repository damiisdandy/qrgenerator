package handlers

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	qrcode "github.com/skip2/go-qrcode"
)

type Handlers struct {
	Logger   *slog.Logger
	Template *template.Template
}

func (h *Handlers) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if err := h.Template.ExecuteTemplate(w, "index.html", nil); err != nil {
		h.Logger.Error("Failed to write response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) Generate(w http.ResponseWriter, r *http.Request) {
	text := ""
	if r.Method == "POST" {
		text = r.FormValue("text")
	} else {
		text = r.URL.Query().Get("text")
	}

	if text == "" {
		h.Logger.Error("No text provided")
		http.Error(w, "No text provided", http.StatusBadRequest)
		return
	}

	// Generate QR code
	qr, err := qrcode.Encode(text, qrcode.Highest, 512)
	if err != nil {
		h.Logger.Error("Failed to generate QR code", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set headers for file download
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", "attachment; filename=qrcode.png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(qr)))

	// Write the QR code image to the response
	if _, err := w.Write(qr); err != nil {
		h.Logger.Error("Failed to write QR code to response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
