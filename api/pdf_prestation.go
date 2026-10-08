package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-pdf/fpdf"
)

func gererPDFPrestation(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	idPrestation, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, "ID de prestation invalide", http.StatusBadRequest)
		return
	}

	var prestation Prestation

	err = db.QueryRow(`
		SELECT id, titre, categorie, tarif, statut
		FROM prestations
		WHERE id = $1
	`, idPrestation).Scan(
		&prestation.ID,
		&prestation.Titre,
		&prestation.Categorie,
		&prestation.Tarif,
		&prestation.Statut,
	)

	if err != nil {
		http.Error(
			w,
			"Prestation introuvable",
			http.StatusNotFound,
		)
		return
	}

	pdf := fpdf.New("P", "mm", "A4", "")

	pdf.AddPage()

	pdf.SetTitle(
		"Récapitulatif de la prestation",
		true,
	)

	pdf.SetFont("Arial", "B", 18)

	pdf.Cell(
		0,
		12,
		"Recapitulatif de la prestation",
		"",
		1,
		"C",
		"",
		0,
	)

	pdf.Ln(10)

	// Titre

	pdf.SetFont("Arial", "B", 12)

	pdf.Cell(
		50,
		10,
		"Titre",
		"1",
		0,
		"L",
		false,
		0,
	)

	pdf.SetFont("Arial", "", 12)

	pdf.Cell(
		130,
		10,
		prestation.Titre,
		"1",
		1,
		"L",
		false,
		0,
	)

	// Catégorie

	pdf.SetFont("Arial", "B", 12)

	pdf.Cell(
		50,
		10,
		"Categorie",
		"1",
		0,
		"L",
		false,
		0,
	)

	pdf.SetFont("Arial", "", 12)

	pdf.Cell(
		130,
		10,
		prestation.Categorie,
		"1",
		1,
		"L",
		false,
		0,
	)

	// Tarif

	pdf.SetFont("Arial", "B", 12)

	pdf.Cell(
		50,
		10,
		"Tarif",
		"1",
		0,
		"L",
		false,
		0,
	)

	pdf.SetFont("Arial", "", 12)

	pdf.Cell(
		130,
		10,
		fmt.Sprintf("%.2f EUR", prestation.Tarif),
		"1",
		1,
		"L",
		false,
		0,
	)

	// Statut

	pdf.SetFont("Arial", "B", 12)

	pdf.Cell(
		50,
		10,
		"Statut",
		"1",
		0,
		"L",
		false,
		0,
	)

	pdf.SetFont("Arial", "", 12)

	pdf.Cell(
		130,
		10,
		prestation.Statut,
		"1",
		1,
		"L",
		false,
		0,
	)

	pdf.Ln(15)

	pdf.SetFont("Arial", "I", 10)

	pdf.Cell(
		0,
		10,
		"Document genere automatiquement par UpcycleConnect",
		"",
		1,
		"C",
		"",
		0,
	)

	w.Header().Set(
		"Content-Type",
		"application/pdf",
	)

	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf(
			`attachment; filename="prestation_%d.pdf"`,
			prestation.ID,
		),
	)

	err = pdf.Output(w)

	if err != nil {
		http.Error(
			w,
			"Erreur lors de la generation du PDF",
			http.StatusInternalServerError,
		)

		return
	}
}
