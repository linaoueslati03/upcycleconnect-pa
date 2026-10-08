package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-pdf/fpdf"
)

// gererPDFPrestation renvoie le récapitulatif d'une prestation sous forme de fichier PDF.
func gererPDFPrestation(w http.ResponseWriter, r *http.Request) {
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	var prestation Prestation
	err = db.QueryRow(
		"SELECT id, titre, categorie, tarif, statut FROM prestations WHERE id = $1", id,
	).Scan(&prestation.ID, &prestation.Titre, &prestation.Categorie, &prestation.Tarif, &prestation.Statut)

	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "prestation introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	// Les polices intégrées (Arial) ne gèrent pas l'UTF-8 : on convertit chaque texte
	// en cp1252 pour que les accents et le symbole € s'affichent correctement.
	texte := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetTitle("Récapitulatif de la prestation", true)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(0, 12, texte("Récapitulatif de la prestation"), "", 1, "C", false, 0, "")
	pdf.Ln(10)

	ajouterLigne(pdf, texte("Titre"), texte(prestation.Titre))
	ajouterLigne(pdf, texte("Catégorie"), texte(prestation.Categorie))
	ajouterLigne(pdf, texte("Tarif"), texte(fmt.Sprintf("%.2f €", prestation.Tarif)))
	ajouterLigne(pdf, texte("Statut"), texte(prestation.Statut))

	pdf.Ln(15)
	pdf.SetFont("Arial", "I", 10)
	pdf.CellFormat(0, 10, texte("Document généré automatiquement par UpcycleConnect"), "", 1, "C", false, 0, "")

	// Le PDF est d'abord généré en mémoire : en cas d'erreur, on peut encore renvoyer
	// un code 500, ce qui n'est plus possible une fois l'envoi au client commencé.
	var contenu bytes.Buffer
	if err := pdf.Output(&contenu); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur lors de la génération du PDF")
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="prestation_%d.pdf"`, prestation.ID))
	w.Write(contenu.Bytes())
}

// ajouterLigne dessine une ligne du tableau : le libellé en gras à gauche, la valeur à droite.
func ajouterLigne(pdf *fpdf.Fpdf, libelle, valeur string) {
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(50, 10, libelle, "1", 0, "L", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(130, 10, valeur, "1", 1, "L", false, 0, "")
}
