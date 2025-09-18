// Package craft gère le système de forge et de création d'objets
// Permet de créer des équipements et potions à partir de ressources récoltées
package craft

import (
	"fmt"
	"world_of_milousques/character"
	"world_of_milousques/item"
	"world_of_milousques/ui"
	"world_of_milousques/utils"
)

// Ingredient représente un ingrédient nécessaire pour une recette
type Ingredient struct {
	Item     item.Item
	Quantite int
}

// Recette représente une recette de craft
type Recette struct {
	Nom          string
	Description  string
	Ingredients  []Ingredient
	Produit      item.Item
	QuantiteProduit int
}

// GetRecettesDisponibles retourne toutes les recettes de craft disponibles
func GetRecettesDisponibles() []Recette {
	return []Recette{
		// === ÉQUIPEMENTS EN MÉTAL (10 de chaque ressource) ===
		{
			Nom:         "Casque en Métal",
			Description: "Protection de tête métallique résistante (+10 défense)",
			Ingredients: []Ingredient{
				{Item: item.NewItem("Bois"), Quantite: 10},
				{Item: item.NewItem("Fer"), Quantite: 10},
				{Item: item.NewItem("Blé"), Quantite: 10},
				{Item: item.NewItem("Laitue Vireuse"), Quantite: 10},
				{Item: item.NewItem("Pichon"), Quantite: 10},
			},
			Produit:         item.NewItem("Casque en Métal"),
			QuantiteProduit: 1,
		},
		{
			Nom:         "Torse en Métal",
			Description: "Protection du torse métallique résistante (+10 défense)",
			Ingredients: []Ingredient{
				{Item: item.NewItem("Bois"), Quantite: 10},
				{Item: item.NewItem("Fer"), Quantite: 10},
				{Item: item.NewItem("Blé"), Quantite: 10},
				{Item: item.NewItem("Laitue Vireuse"), Quantite: 10},
				{Item: item.NewItem("Pichon"), Quantite: 10},
			},
			Produit:         item.NewItem("Torse en Métal"),
			QuantiteProduit: 1,
		},
		{
			Nom:         "Jambières en Métal",
			Description: "Protection des jambes métallique résistante (+10 défense)",
			Ingredients: []Ingredient{
				{Item: item.NewItem("Bois"), Quantite: 10},
				{Item: item.NewItem("Fer"), Quantite: 10},
				{Item: item.NewItem("Blé"), Quantite: 10},
				{Item: item.NewItem("Laitue Vireuse"), Quantite: 10},
				{Item: item.NewItem("Pichon"), Quantite: 10},
			},
			Produit:         item.NewItem("Jambières en Métal"),
			QuantiteProduit: 1,
		},
		// === ARMES D'EXPERT (20 de chaque ressource) ===
		{
			Nom:         "Bâton d'Expert",
			Description: "Bâton magique d'expert pour mage (+20 attaque)",
			Ingredients: []Ingredient{
				{Item: item.NewItem("Bois"), Quantite: 20},
				{Item: item.NewItem("Fer"), Quantite: 20},
				{Item: item.NewItem("Blé"), Quantite: 20},
				{Item: item.NewItem("Laitue Vireuse"), Quantite: 20},
				{Item: item.NewItem("Pichon"), Quantite: 20},
			},
			Produit:         item.NewItem("Bâton d'Expert"),
			QuantiteProduit: 1,
		},
		{
			Nom:         "Épée d'Expert",
			Description: "Épée forgée d'expert pour guerrier (+20 attaque)",
			Ingredients: []Ingredient{
				{Item: item.NewItem("Bois"), Quantite: 20},
				{Item: item.NewItem("Fer"), Quantite: 20},
				{Item: item.NewItem("Blé"), Quantite: 20},
				{Item: item.NewItem("Laitue Vireuse"), Quantite: 20},
				{Item: item.NewItem("Pichon"), Quantite: 20},
			},
			Produit:         item.NewItem("Épée d'Expert"),
			QuantiteProduit: 1,
		},
		{
			Nom:         "Dague d'Expert",
			Description: "Dague empoisonnée d'expert pour voleur (+20 attaque)",
			Ingredients: []Ingredient{
				{Item: item.NewItem("Bois"), Quantite: 20},
				{Item: item.NewItem("Fer"), Quantite: 20},
				{Item: item.NewItem("Blé"), Quantite: 20},
				{Item: item.NewItem("Laitue Vireuse"), Quantite: 20},
				{Item: item.NewItem("Pichon"), Quantite: 20},
			},
			Produit:         item.NewItem("Dague d'Expert"),
			QuantiteProduit: 1,
		},
		// === POTIONS ===
		{
			Nom:         "Potion de Vie",
			Description: "Potion qui restaure 50 PV",
			Ingredients: []Ingredient{
				{Item: item.NewItem("Laitue Vireuse"), Quantite: 3},
				{Item: item.NewItem("Pichon"), Quantite: 3},
			},
			Produit:         item.NewItem("Potion de Vie"),
			QuantiteProduit: 1,
		},
		{
			Nom:         "Potion de Mana",
			Description: "Potion qui restaure 50 Mana",
			Ingredients: []Ingredient{
				{Item: item.NewItem("Laitue Vireuse"), Quantite: 3},
				{Item: item.NewItem("Blé"), Quantite: 3},
			},
			Produit:         item.NewItem("Potion de Mana"),
			QuantiteProduit: 1,
		},
	}
}

// AfficherForge affiche le menu principal de la forge
func AfficherForge(joueur *character.Character) {
	for {
		fmt.Println("\n🔨 === FORGE D'ASTRAB === 🔨")
		fmt.Println("Maître Forgeron : Bienvenue dans ma forge ! Que puis-je créer pour vous ?")
		
		options := []string{
			"Voir les recettes disponibles",
			"Crafter un objet",
			"Voir mon inventaire",
			"Quitter la forge",
		}
		
		ui.AfficherMenu("Forge", options)
		choix := utils.ScanChoice("Que voulez-vous faire ? ", options)
		
		switch choix {
		case 1:
			afficherRecettes()
		case 2:
			crafterObjet(joueur)
		case 3:
			joueur.Inventaire.Afficher()
			fmt.Println("\nAppuyez sur Entrée pour continuer...")
			fmt.Scanln()
		case 4:
			fmt.Println("Maître Forgeron : Revenez quand vous voulez !")
			return
		}
	}
}

// afficherRecettes affiche toutes les recettes disponibles
func afficherRecettes() {
	recettes := GetRecettesDisponibles()
	
	fmt.Println("\n📜 === RECETTES DISPONIBLES === 📜")
	for i, recette := range recettes {
		fmt.Printf("\n%d. %s\n", i+1, recette.Nom)
		fmt.Printf("   Description: %s\n", recette.Description)
		fmt.Printf("   Produit: %dx %s\n", recette.QuantiteProduit, recette.Produit.Nom)
		fmt.Printf("   Ingrédients requis:\n")
		for _, ingredient := range recette.Ingredients {
			fmt.Printf("     - %dx %s\n", ingredient.Quantite, ingredient.Item.Nom)
		}
	}
	
	fmt.Println("\nAppuyez sur Entrée pour continuer...")
	fmt.Scanln()
}

// crafterObjet permet au joueur de crafter un objet
func crafterObjet(joueur *character.Character) {
	recettes := GetRecettesDisponibles()
	
	fmt.Println("\n⚒️  === CRÉATION D'OBJET === ⚒️")
	
	// Créer les options du menu avec les recettes
	options := make([]string, 0)
	for i, recette := range recettes {
		disponible := "✅"
		if !peutCrafter(joueur, recette) {
			disponible = "❌"
		}
		options = append(options, fmt.Sprintf("%s %s (%dx %s)", 
			disponible, recette.Nom, recette.QuantiteProduit, recette.Produit.Nom))
		_ = i // éviter unused variable
	}
	options = append(options, "Retour")
	
	ui.AfficherMenu("Choisir une recette à crafter", options)
	choix := utils.ScanChoice("Quelle recette voulez-vous utiliser ? ", options)
	
	if choix == len(options) {
		return // Retour
	}
	
	recetteChoisie := recettes[choix-1]
	
	if !peutCrafter(joueur, recetteChoisie) {
		fmt.Println("\n❌ Vous n'avez pas les ingrédients nécessaires pour cette recette !")
		fmt.Println("\nIngrédients requis :")
		for _, ingredient := range recetteChoisie.Ingredients {
			quantitePossedee := compterItem(joueur, ingredient.Item.Nom)
			fmt.Printf("  - %s : %d/%d %s\n", 
				ingredient.Item.Nom, 
				quantitePossedee, 
				ingredient.Quantite,
				getStatusIcon(quantitePossedee >= ingredient.Quantite))
		}
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}
	
	// Confirmation
	fmt.Printf("\n🔨 Crafter : %s\n", recetteChoisie.Nom)
	fmt.Printf("Produit : %dx %s\n", recetteChoisie.QuantiteProduit, recetteChoisie.Produit.Nom)
	
	options = []string{"Confirmer le craft", "Annuler"}
	ui.AfficherMenu("Confirmation", options)
	confirmation := utils.ScanChoice("Êtes-vous sûr ? ", options)
	
	if confirmation == 1 {
		// Effectuer le craft
		retirerIngredients(joueur, recetteChoisie)
		ajouterProduit(joueur, recetteChoisie)
		
		fmt.Printf("\n✅ %s créé avec succès !\n", recetteChoisie.Nom)
		fmt.Printf("Vous avez reçu : %dx %s\n", recetteChoisie.QuantiteProduit, recetteChoisie.Produit.Nom)
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
	}
}

// peutCrafter vérifie si le joueur a les ingrédients nécessaires
func peutCrafter(joueur *character.Character, recette Recette) bool {
	for _, ingredient := range recette.Ingredients {
		if compterItem(joueur, ingredient.Item.Nom) < ingredient.Quantite {
			return false
		}
	}
	return true
}

// compterItem compte combien d'exemplaires d'un item le joueur possède
func compterItem(joueur *character.Character, nomItem string) int {
	count := 0
	for _, item := range joueur.Inventaire.Items {
		if item.Nom == nomItem {
			count++
		}
	}
	return count
}

// getStatusIcon retourne l'icône de statut pour les ingrédients
func getStatusIcon(disponible bool) string {
	if disponible {
		return "✅"
	}
	return "❌"
}

// retirerIngredients retire les ingrédients de l'inventaire du joueur
func retirerIngredients(joueur *character.Character, recette Recette) {
	for _, ingredient := range recette.Ingredients {
		retirees := 0
		nouvelInventaire := make([]item.Item, 0)
		
		for _, item := range joueur.Inventaire.Items {
			if item.Nom == ingredient.Item.Nom && retirees < ingredient.Quantite {
				retirees++
				// Ne pas ajouter cet item au nouvel inventaire (= le retirer)
			} else {
				nouvelInventaire = append(nouvelInventaire, item)
			}
		}
		
		joueur.Inventaire.Items = nouvelInventaire
	}
}

// ajouterProduit ajoute le produit crafté à l'inventaire
func ajouterProduit(joueur *character.Character, recette Recette) {
	// Cas spéciaux pour les potions
	if recette.Produit.Nom == "Potion de Vie" {
		joueur.Inventaire.Potions += recette.QuantiteProduit
		return
	}
	if recette.Produit.Nom == "Potion de Mana" {
		joueur.Inventaire.PotionsMana += recette.QuantiteProduit
		return
	}
	
	// Cas normal pour les objets
	for i := 0; i < recette.QuantiteProduit; i++ {
		joueur.Inventaire.Items = append(joueur.Inventaire.Items, recette.Produit)
	}
}
