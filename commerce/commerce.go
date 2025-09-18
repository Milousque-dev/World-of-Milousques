// Package commerce gère le système de marchand et de vente/achat d'objets
// Permet l'achat d'équipements et potions, ainsi que la vente de ressources
package commerce

import (
	"fmt"
	"world_of_milousques/character"
	"world_of_milousques/item"
	"world_of_milousques/ui"
	"world_of_milousques/utils"
)

// Article représente un article vendu par le marchand
type Article struct {
	Item     item.Item
	Prix     int
	Stock    int
	Illimite bool // Si true, stock infini
}

// Marchand représente un marchand avec son inventaire
type Marchand struct {
	Nom       string
	Salut     string
	Articles  []Article
}

// GetMarchandAstrab retourne le marchand principal d'Astrab
func GetMarchandAstrab() Marchand {
	return Marchand{
		Nom:   "Maître Karim le Marchand",
		Salut: "Bienvenue dans ma boutique ! J'ai tout l'équipement qu'il vous faut, aventurier !",
		Articles: []Article{
			// Équipements en cuir (150 or chacun)
			{Item: item.NewItem("Casque en Cuir"), Prix: 150, Stock: 5, Illimite: false},
			{Item: item.NewItem("Torse en Cuir"), Prix: 150, Stock: 5, Illimite: false},
			{Item: item.NewItem("Jambières en Cuir"), Prix: 150, Stock: 5, Illimite: false},
			// Armes simples (250 or chacune)
			{Item: item.NewItem("Bâton Simple"), Prix: 250, Stock: 3, Illimite: false},
			{Item: item.NewItem("Épée Simple"), Prix: 250, Stock: 3, Illimite: false},
			{Item: item.NewItem("Dague Simple"), Prix: 250, Stock: 3, Illimite: false},
			// Potions (stock illimité)
			{Item: item.NewItem("Potion de Vie"), Prix: 50, Stock: 0, Illimite: true},
			{Item: item.NewItem("Potion de Mana"), Prix: 50, Stock: 0, Illimite: true},
		},
	}
}

// AfficherMarchand affiche le menu principal du marchand
func AfficherMarchand(joueur *character.Character) {
	marchand := GetMarchandAstrab()
	
	for {
		fmt.Printf("\n💰 === %s === 💰\n", marchand.Nom)
		fmt.Printf("%s\n", marchand.Salut)
		fmt.Printf("💳 Votre argent : %d pièces d'or\n", joueur.Argent)
		
		options := []string{
			"Voir les articles à vendre",
			"Acheter un article",
			"Vendre mes objets",
			"Voir mon inventaire",
			"Quitter la boutique",
		}
		
		ui.AfficherMenu("Boutique", options)
		choix := utils.ScanChoice("Que voulez-vous faire ? ", options)
		
		switch choix {
		case 1:
			afficherArticles(marchand)
		case 2:
			acheterArticle(joueur, &marchand)
		case 3:
			vendreObjets(joueur)
		case 4:
			joueur.Inventaire.Afficher()
			fmt.Println("\nAppuyez sur Entrée pour continuer...")
			fmt.Scanln()
		case 5:
			fmt.Printf("%s : Merci de votre visite ! Revenez quand vous voulez !\n", marchand.Nom)
			return
		}
	}
}

// afficherArticles affiche tous les articles du marchand
func afficherArticles(marchand Marchand) {
	fmt.Println("\n🛒 === ARTICLES DISPONIBLES === 🛒")
	
	for i, article := range marchand.Articles {
		stockInfo := fmt.Sprintf("(%d en stock)", article.Stock)
		if article.Illimite {
			stockInfo = "(Stock illimité)"
		}
		
		disponible := "✅"
		if article.Stock == 0 && !article.Illimite {
			disponible = "❌ RUPTURE"
			stockInfo = ""
		}
		
		fmt.Printf("%d. %s %s - %d pièces d'or %s\n", 
			i+1, disponible, article.Item.Nom, article.Prix, stockInfo)
		fmt.Printf("   %s\n", article.Item.Effet)
	}
	
	fmt.Println("\nAppuyez sur Entrée pour continuer...")
	fmt.Scanln()
}

// acheterArticle permet d'acheter un article
func acheterArticle(joueur *character.Character, marchand *Marchand) {
	fmt.Println("\n💳 === ACHAT D'ARTICLE === 💳")
	fmt.Printf("Votre argent : %d pièces d'or\n", joueur.Argent)
	
	// Créer les options du menu
	options := make([]string, 0)
	for i, article := range marchand.Articles {
		disponible := "✅"
		prixInfo := fmt.Sprintf("- %d or", article.Prix)
		
		if article.Stock == 0 && !article.Illimite {
			disponible = "❌"
			prixInfo = "- RUPTURE"
		} else if joueur.Argent < article.Prix {
			disponible = "💸"
			prixInfo = fmt.Sprintf("- %d or (trop cher)", article.Prix)
		}
		
		options = append(options, fmt.Sprintf("%s %s %s", 
			disponible, article.Item.Nom, prixInfo))
		_ = i // éviter unused variable
	}
	options = append(options, "Retour")
	
	ui.AfficherMenu("Choisir un article à acheter", options)
	choix := utils.ScanChoice("Quel article voulez-vous acheter ? ", options)
	
	if choix == len(options) {
		return // Retour
	}
	
	articleChoisi := &marchand.Articles[choix-1]
	
	// Vérifications
	if articleChoisi.Stock == 0 && !articleChoisi.Illimite {
		fmt.Println("❌ Cet article n'est plus en stock !")
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}
	
	if joueur.Argent < articleChoisi.Prix {
		fmt.Printf("💸 Vous n'avez pas assez d'argent ! Il vous faut %d pièces d'or.\n", 
			articleChoisi.Prix)
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}
	
	// Confirmation
	fmt.Printf("\n🛒 Acheter : %s\n", articleChoisi.Item.Nom)
	fmt.Printf("Prix : %d pièces d'or\n", articleChoisi.Prix)
	fmt.Printf("Argent restant : %d pièces d'or\n", joueur.Argent-articleChoisi.Prix)
	
	options = []string{"Confirmer l'achat", "Annuler"}
	ui.AfficherMenu("Confirmation", options)
	confirmation := utils.ScanChoice("Êtes-vous sûr ? ", options)
	
		if confirmation == 1 {
			// Effectuer l'achat
			joueur.Argent -= articleChoisi.Prix
			
			// Cas spéciaux pour les potions
			if articleChoisi.Item.Nom == "Potion de Vie" {
				joueur.Inventaire.Potions++
			} else if articleChoisi.Item.Nom == "Potion de Mana" {
				joueur.Inventaire.PotionsMana++
			} else {
				// Cas normal pour les objets
				joueur.Inventaire.Items = append(joueur.Inventaire.Items, articleChoisi.Item)
			}
			
			// Diminuer le stock si pas illimité
			if !articleChoisi.Illimite {
				articleChoisi.Stock--
			}
			
			fmt.Printf("✅ %s acheté avec succès !\n", articleChoisi.Item.Nom)
			fmt.Printf("Argent restant : %d pièces d'or\n", joueur.Argent)
			fmt.Println("\nAppuyez sur Entrée pour continuer...")
			fmt.Scanln()
		}
}

// vendreObjets permet de vendre des objets de l'inventaire
func vendreObjets(joueur *character.Character) {
	if len(joueur.Inventaire.Items) == 0 {
		fmt.Println("❌ Votre inventaire est vide !")
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}
	
	fmt.Println("\n💰 === VENTE D'OBJETS === 💰")
	fmt.Printf("Votre argent actuel : %d pièces d'or\n", joueur.Argent)
	
	// Grouper les objets identiques
	objetsGroupes := make(map[string]GroupeObjet)
	for _, objet := range joueur.Inventaire.Items {
		if groupe, existe := objetsGroupes[objet.Nom]; existe {
			groupe.Quantite++
			objetsGroupes[objet.Nom] = groupe
		} else {
			objetsGroupes[objet.Nom] = GroupeObjet{
				Item:     objet,
				Quantite: 1,
				PrixVente: objet.Valeur / 2, // Le marchand achète à 50% du prix
			}
		}
	}
	
	// Créer les options
	options := make([]string, 0)
	groupes := make([]GroupeObjet, 0)
	
	for _, groupe := range objetsGroupes {
		options = append(options, fmt.Sprintf("%s (%dx) - %d or chacun", 
			groupe.Item.Nom, groupe.Quantite, groupe.PrixVente))
		groupes = append(groupes, groupe)
	}
	options = append(options, "Retour")
	
	ui.AfficherMenu("Choisir un objet à vendre", options)
	choix := utils.ScanChoice("Quel objet voulez-vous vendre ? ", options)
	
	if choix == len(options) {
		return // Retour
	}
	
	groupeChoisi := groupes[choix-1]
	
	// Demander la quantité si plusieurs exemplaires
	quantiteAVendre := 1
	if groupeChoisi.Quantite > 1 {
		quantiteAVendre = utils.ScanInt(
			fmt.Sprintf("Combien voulez-vous en vendre ? (max %d) : ", groupeChoisi.Quantite),
			1, groupeChoisi.Quantite)
	}
	
	prixTotal := groupeChoisi.PrixVente * quantiteAVendre
	
	// Confirmation
	fmt.Printf("\n💰 Vendre : %dx %s\n", quantiteAVendre, groupeChoisi.Item.Nom)
	fmt.Printf("Prix total : %d pièces d'or\n", prixTotal)
	fmt.Printf("Argent après vente : %d pièces d'or\n", joueur.Argent+prixTotal)
	
	options = []string{"Confirmer la vente", "Annuler"}
	ui.AfficherMenu("Confirmation", options)
	confirmation := utils.ScanChoice("Êtes-vous sûr ? ", options)
	
	if confirmation == 1 {
		// Effectuer la vente
		retirerObjets(joueur, groupeChoisi.Item.Nom, quantiteAVendre)
		joueur.Argent += prixTotal
		
		fmt.Printf("✅ %dx %s vendu avec succès !\n", quantiteAVendre, groupeChoisi.Item.Nom)
		fmt.Printf("Vous avez gagné : %d pièces d'or\n", prixTotal)
		fmt.Printf("Argent total : %d pièces d'or\n", joueur.Argent)
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
	}
}

// GroupeObjet représente un groupe d'objets identiques
type GroupeObjet struct {
	Item      item.Item
	Quantite  int
	PrixVente int
}

// retirerObjets retire une quantité spécifiée d'un objet de l'inventaire
func retirerObjets(joueur *character.Character, nomObjet string, quantite int) {
	retirees := 0
	nouvelInventaire := make([]item.Item, 0)
	
	for _, objet := range joueur.Inventaire.Items {
		if objet.Nom == nomObjet && retirees < quantite {
			retirees++
			// Ne pas ajouter cet objet au nouvel inventaire (= le retirer)
		} else {
			nouvelInventaire = append(nouvelInventaire, objet)
		}
	}
	
	joueur.Inventaire.Items = nouvelInventaire
}