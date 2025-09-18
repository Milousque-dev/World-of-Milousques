// Package banque gère le système de coffre-fort personnel pour chaque joueur
// Permet de stocker et récupérer des objets avec une capacité limitée
package banque

import (
	"encoding/json"
	"fmt"
	"os"
	"world_of_milousques/character"
	"world_of_milousques/item"
	"world_of_milousques/ui"
	"world_of_milousques/utils"
)

// Banque représente le coffre-fort d'un joueur
type Banque struct {
	Proprietaire string      `json:"proprietaire"`
	Objets       []item.Item `json:"objets"`
	MaxCapacite  int         `json:"max_capacite"`
}

// NewBanque crée une nouvelle banque pour un joueur
func NewBanque(proprietaire string) *Banque {
	return &Banque{
		Proprietaire: proprietaire,
		Objets:       []item.Item{},
		MaxCapacite:  200, // La banque peut contenir 200 objets
	}
}

// ChargerBanque charge la banque d'un joueur depuis un fichier
func ChargerBanque(proprietaire string) (*Banque, error) {
	filename := "saves/banque_" + proprietaire + ".json"
	
	// Si le fichier n'existe pas, créer une nouvelle banque
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return NewBanque(proprietaire), nil
	}
	
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	var banque Banque
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&banque)
	if err != nil {
		return nil, err
	}
	
	return &banque, nil
}

// Sauvegarder sauvegarde la banque dans un fichier
func (b *Banque) Sauvegarder() error {
	filename := "saves/banque_" + b.Proprietaire + ".json"
	
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(b)
}

// AjouterObjet ajoute un objet à la banque
func (b *Banque) AjouterObjet(objet item.Item) bool {
	if len(b.Objets) >= b.MaxCapacite {
		return false // Banque pleine
	}
	
	b.Objets = append(b.Objets, objet)
	return true
}

// RetirerObjet retire un objet de la banque par index
func (b *Banque) RetirerObjet(index int) (item.Item, bool) {
	if index < 0 || index >= len(b.Objets) {
		return item.Item{}, false
	}
	
	objet := b.Objets[index]
	
	// Créer une nouvelle slice sans l'objet retiré
	nouveauxObjets := make([]item.Item, 0)
	for i, obj := range b.Objets {
		if i != index {
			nouveauxObjets = append(nouveauxObjets, obj)
		}
	}
	b.Objets = nouveauxObjets
	
	return objet, true
}

// AfficherBanque gère l'interface de la banque
func AfficherBanque(joueur *character.Character) {
	banque, err := ChargerBanque(joueur.Nom)
	if err != nil {
		fmt.Printf("Erreur lors du chargement de votre coffre : %v\n", err)
		return
	}
	
	for {
		fmt.Println("\n🏦 === BANQUE ROYALE D'ASTRAB === 🏦")
		fmt.Printf("Banquier Salomon : Bienvenue %s ! Votre coffre-fort vous attend.\n", joueur.Nom)
		fmt.Printf("💰 Capacité du coffre : %d/%d objets\n", len(banque.Objets), banque.MaxCapacite)
		fmt.Printf("🎒 Votre inventaire : %d/100 objets\n", len(joueur.Inventaire.Items))
		
		options := []string{
			"🏦 Déposer des objets",
			"📤 Retirer des objets", 
			"📋 Voir le contenu du coffre",
			"🎒 Voir mon inventaire",
			"🚪 Quitter la banque",
		}
		
		ui.AfficherMenu("Services bancaires", options)
		choix := utils.ScanChoice("Que souhaitez-vous faire ? ", options)
		
		switch choix {
		case 1:
			deposerObjets(joueur, banque)
		case 2:
			retirerObjets(joueur, banque)
		case 3:
			afficherContenuBanque(banque)
		case 4:
			joueur.Inventaire.Afficher()
			fmt.Println("\nAppuyez sur Entrée pour continuer...")
			fmt.Scanln()
		case 5:
			// Sauvegarder avant de quitter
			if err := banque.Sauvegarder(); err != nil {
				fmt.Printf("Erreur lors de la sauvegarde : %v\n", err)
			} else {
				fmt.Println("Banquier Salomon : Vos biens sont en sécurité ! À bientôt !")
			}
			return
		}
	}
}

// deposerObjets gère le dépôt d'objets dans la banque
func deposerObjets(joueur *character.Character, banque *Banque) {
	if len(joueur.Inventaire.Items) == 0 {
		fmt.Println("❌ Votre inventaire est vide !")
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}
	
	if len(banque.Objets) >= banque.MaxCapacite {
		fmt.Println("❌ Votre coffre est plein ! Retirez d'abord des objets.")
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}
	
	fmt.Println("\n💰 === DÉPOSER DES OBJETS === 💰")
	fmt.Printf("Espace disponible dans le coffre : %d objets\n\n", banque.MaxCapacite-len(banque.Objets))
	
	// Grouper les objets identiques
	objetsGroupes := make(map[string]GroupeObjet)
	for i, objet := range joueur.Inventaire.Items {
		if groupe, existe := objetsGroupes[objet.Nom]; existe {
			groupe.Quantite++
			groupe.Indices = append(groupe.Indices, i)
			objetsGroupes[objet.Nom] = groupe
		} else {
			objetsGroupes[objet.Nom] = GroupeObjet{
				Item:     objet,
				Quantite: 1,
				Indices:  []int{i},
			}
		}
	}
	
	// Créer les options
	options := make([]string, 0)
	groupes := make([]GroupeObjet, 0)
	
	for _, groupe := range objetsGroupes {
		options = append(options, fmt.Sprintf("%s (%dx)", groupe.Item.Nom, groupe.Quantite))
		groupes = append(groupes, groupe)
	}
	options = append(options, "Retour")
	
	ui.AfficherMenu("Choisir un objet à déposer", options)
	choix := utils.ScanChoice("Quel objet voulez-vous déposer ? ", options)
	
	if choix == len(options) {
		return // Retour
	}
	
	groupeChoisi := groupes[choix-1]
	
	// Demander la quantité si plusieurs exemplaires
	quantiteADeposer := 1
	maxDeposable := min(groupeChoisi.Quantite, banque.MaxCapacite-len(banque.Objets))
	
	if groupeChoisi.Quantite > 1 && maxDeposable > 1 {
		quantiteADeposer = utils.ScanInt(
			fmt.Sprintf("Combien voulez-vous en déposer ? (max %d) : ", maxDeposable),
			1, maxDeposable)
	}
	
	// Effectuer le dépôt
	for i := 0; i < quantiteADeposer; i++ {
		if !banque.AjouterObjet(groupeChoisi.Item) {
			fmt.Println("❌ Le coffre est plein !")
			break
		}
	}
	
	// Retirer les objets de l'inventaire du joueur
	retirerObjetsInventaire(joueur, groupeChoisi.Item.Nom, quantiteADeposer)
	
	fmt.Printf("✅ %dx %s déposé avec succès dans votre coffre !\n", quantiteADeposer, groupeChoisi.Item.Nom)
	fmt.Println("\nAppuyez sur Entrée pour continuer...")
	fmt.Scanln()
}

// retirerObjets gère le retrait d'objets de la banque
func retirerObjets(joueur *character.Character, banque *Banque) {
	if len(banque.Objets) == 0 {
		fmt.Println("❌ Votre coffre est vide !")
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}
	
	if len(joueur.Inventaire.Items) >= 100 {
		fmt.Println("❌ Votre inventaire est plein ! Videz d'abord votre inventaire.")
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}
	
	fmt.Println("\n📤 === RETIRER DES OBJETS === 📤")
	fmt.Printf("Espace disponible dans l'inventaire : %d objets\n\n", 100-len(joueur.Inventaire.Items))
	
	// Afficher le contenu du coffre
	afficherContenuBanque(banque)
	
	if len(banque.Objets) == 0 {
		return
	}
	
	choix := utils.ScanInt("Quel objet voulez-vous retirer ? (numéro) : ", 1, len(banque.Objets))
	
	objet, success := banque.RetirerObjet(choix - 1)
	if success {
		joueur.Inventaire.Items = append(joueur.Inventaire.Items, objet)
		fmt.Printf("✅ %s retiré avec succès de votre coffre !\n", objet.Nom)
	} else {
		fmt.Println("❌ Erreur lors du retrait de l'objet.")
	}
	
	fmt.Println("\nAppuyez sur Entrée pour continuer...")
	fmt.Scanln()
}

// afficherContenuBanque affiche le contenu du coffre
func afficherContenuBanque(banque *Banque) {
	fmt.Printf("\n📋 === CONTENU DU COFFRE === 📋\n")
	
	if len(banque.Objets) == 0 {
		fmt.Println("Votre coffre est vide.")
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}
	
	for i, objet := range banque.Objets {
		fmt.Printf("%d. %s | Poids: %d | Effet: %s | Valeur: %d or\n", 
			i+1, objet.Nom, objet.Poids, objet.Effet, objet.Valeur)
	}
	
	fmt.Printf("\nTotal : %d/%d objets\n", len(banque.Objets), banque.MaxCapacite)
	
	if len(banque.Objets) < 20 { // Si pas trop d'objets, pas besoin d'appuyer sur Entrée
		return
	}
	
	fmt.Println("\nAppuyez sur Entrée pour continuer...")
	fmt.Scanln()
}

// GroupeObjet représente un groupe d'objets identiques avec leurs indices
type GroupeObjet struct {
	Item     item.Item
	Quantite int
	Indices  []int
}

// retirerObjetsInventaire retire une quantité spécifiée d'un objet de l'inventaire
func retirerObjetsInventaire(joueur *character.Character, nomObjet string, quantite int) {
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

// min retourne le minimum entre deux entiers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}