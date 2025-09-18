package ui

import (
	"fmt"
	"strings"
	"world_of_milousques/sorts"
)

// calculerLargeurAffichage calcule la largeur réelle d'affichage d'une chaîne
// en tenant compte des caractères Unicode qui peuvent prendre plus d'espace
func calculerLargeurAffichage(s string) int {
	// Sur Windows, les emojis peuvent causer des problèmes d'alignement
	// On utilise une approche conservatrice : compter chaque emoji comme 2 espaces
	largeur := 0
	for _, r := range s {
		if estEmoji(r) {
			// Les emojis prennent souvent 2 colonnes dans les terminaux Windows
			largeur += 2
		} else {
			largeur += 1
		}
	}
	return largeur
}

// estEmoji détermine si un caractère est un emoji
func estEmoji(r rune) bool {
	// Plages Unicode communes pour les emojis
	return (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols and Pictographs
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map
		(r >= 0x1F1E0 && r <= 0x1F1FF) || // Regional indicators
		(r >= 0x2600 && r <= 0x26FF) ||   // Misc symbols
		(r >= 0x2700 && r <= 0x27BF) ||   // Dingbats
		(r >= 0xFE0F && r <= 0xFE0F)     // Variation Selector
}

// centrerTexte centre un texte dans une largeur donnée (en largeur d'affichage)
func centrerTexte(texte string, largeur int) string {
	w := calculerLargeurAffichage(texte)
	if w >= largeur {
		return texte
	}
	paddingTotal := largeur - w
	paddingGauche := paddingTotal / 2
	paddingDroite := paddingTotal - paddingGauche
	return strings.Repeat(" ", paddingGauche) + texte + strings.Repeat(" ", paddingDroite)
}

// alignerGauche aligne un texte à gauche dans une largeur donnée (en largeur d'affichage)
func alignerGauche(texte string, largeur int) string {
	w := calculerLargeurAffichage(texte)
	if w >= largeur {
		return texte
	}
	return texte + strings.Repeat(" ", largeur-w)
}

// AfficherMenuSimple affiche un menu simple avec bordures Unicode
// Utilise la largeur d'affichage pour des alignements fiables
// Limite la largeur pour éviter le wrapping dans les terminaux Windows
func AfficherMenuSimple(titre string, options []string) {
	// Calculer la largeur maximale en largeur d'affichage
	largeurContenu := calculerLargeurAffichage(titre)
	for i, opt := range options {
		ligne := fmt.Sprintf(" %d) %s", i+1, opt) // même préfixe que l'affichage réel
		lw := calculerLargeurAffichage(ligne)
		if lw > largeurContenu {
			largeurContenu = lw
		}
	}
	
	// Largeur minimale/maximum pour éviter le wrapping
	if largeurContenu < 30 {
		largeurContenu = 30
	}
	if largeurContenu > 50 {
		largeurContenu = 50
	}
	
	// Largeur totale = contenu + 4 espaces de marge
	largeurTotale := largeurContenu + 4
	
	// Ligne supérieure
	ligneBordure := "\u250C" + strings.Repeat("\u2500", largeurTotale) + "\u2510"
	fmt.Println(ligneBordure)
	
	// Titre centré
	titreCentre := centrerTexte(titre, largeurTotale)
	ligneTitre := "\u2502" + titreCentre + "\u2502"
	fmt.Println(ligneTitre)
	
	// Ligne de séparation
	ligneSeparation := "\u251C" + strings.Repeat("\u2500", largeurTotale) + "\u2524"
	fmt.Println(ligneSeparation)
	
	// Options
	for i, opt := range options {
		ligne := fmt.Sprintf(" %d) %s", i+1, opt)
		ligneAlignee := alignerGauche(ligne, largeurTotale)
		ligneOption := "\u2502" + ligneAlignee + "\u2502"
		fmt.Println(ligneOption)
	}
	
	// Ligne inférieure
	ligneBordureInf := "\u2514" + strings.Repeat("\u2500", largeurTotale) + "\u2518"
	fmt.Println(ligneBordureInf)
}

func AfficherMenu(titre string, options []string) {
// Calculer la largeur maximale nécessaire
	largeurContenu := calculerLargeurAffichage(titre)
	for i, opt := range options {
		ligne := fmt.Sprintf(" %d) %s", i+1, opt) // même préfixe que l'affichage
		if calculerLargeurAffichage(ligne) > largeurContenu {
			largeurContenu = calculerLargeurAffichage(ligne)
		}
	}
	
	// Largeur minimale pour un bel affichage, mais limitée pour éviter le wrapping
	if largeurContenu < 30 {
		largeurContenu = 30
	}
	// Largeur maximale pour éviter le wrapping dans les terminaux Windows (plus conservative)
	if largeurContenu > 50 {
		largeurContenu = 50
	}
	
	// Ajouter de la marge (2 espaces de chaque côté)
	largeurTotale := largeurContenu + 4
	
	// Ligne supérieure avec Unicode
	ligneBordure := "\u250C" + strings.Repeat("\u2500", largeurTotale) + "\u2510"
	fmt.Println(ligneBordure)
	
	// Titre centré avec Unicode
	titreCentre := centrerTexte(titre, largeurTotale)
	ligneTitre := "\u2502" + titreCentre + "\u2502"
	fmt.Println(ligneTitre)
	
	// Ligne de séparation avec Unicode
	ligneSeparation := "\u251C" + strings.Repeat("\u2500", largeurTotale) + "\u2524"
	fmt.Println(ligneSeparation)
	
	// Options avec Unicode
	for i, opt := range options {
		ligne := fmt.Sprintf(" %d) %s", i+1, opt)
		ligneAlignee := alignerGauche(ligne, largeurTotale)
		ligneOption := "\u2502" + ligneAlignee + "\u2502"
		fmt.Println(ligneOption)
	}
	
	// Ligne inférieure avec Unicode
	ligneBordureInf := "\u2514" + strings.Repeat("\u2500", largeurTotale) + "\u2518"
	fmt.Println(ligneBordureInf)
}

func AfficherMenuCombat(joueurNom string, joueurPv, joueurPvMax, joueurMana, joueurManaMax int,
	ennemiNom string, ennemiPv int, sortsList []sorts.Sorts, potions, potionsMana int) {

	lignes := []string{}
	lignes = append(lignes, fmt.Sprintf("%s : PV %d/%d | Mana %d/%d", joueurNom, joueurPv, joueurPvMax, joueurMana, joueurManaMax))
	lignes = append(lignes, fmt.Sprintf("Ennemi %s : PV %d", ennemiNom, ennemiPv))
	lignes = append(lignes, "") // Ligne vide pour séparation

	for i, s := range sortsList {
		lignes = append(lignes, fmt.Sprintf("%d) %s - Dégâts: %d, Mana: %d", i+1, s.Nom, s.Degats, s.Cout))
	}

	lignes = append(lignes, fmt.Sprintf("%d) Utiliser une potion de vie (+50 PV) [%d disponibles]", len(sortsList)+1, potions))
	lignes = append(lignes, fmt.Sprintf("%d) Utiliser une potion de mana (+50 Mana) [%d disponibles]", len(sortsList)+2, potionsMana))
	lignes = append(lignes, fmt.Sprintf("%d) Fuir le combat", len(sortsList)+3))

	// Calculer la largeur maximale
	largeurContenu := 0
	for _, ligne := range lignes {
		if calculerLargeurAffichage(ligne) > largeurContenu {
			largeurContenu = calculerLargeurAffichage(ligne)
		}
	}
	
	// Ajouter de la marge (2 espaces de chaque côté)
	largeurTotale := largeurContenu + 4

	// Ligne supérieure avec Unicode
	ligneBordure := "\u250C" + strings.Repeat("\u2500", largeurTotale) + "\u2510"
	fmt.Println(ligneBordure)
	
	// Contenu
	for _, ligne := range lignes {
		if ligne == "" {
			// Ligne de séparation avec Unicode
			ligneSeparation := "\u251C" + strings.Repeat("\u2500", largeurTotale) + "\u2524"
			fmt.Println(ligneSeparation)
		} else {
			ligneAlignee := alignerGauche(ligne, largeurTotale)
			ligneOption := "\u2502" + ligneAlignee + "\u2502"
			fmt.Println(ligneOption)
		}
	}
	
	// Ligne inférieure avec Unicode
	ligneBordureInf := "\u2514" + strings.Repeat("\u2500", largeurTotale) + "\u2518"
	fmt.Println(ligneBordureInf)
}
