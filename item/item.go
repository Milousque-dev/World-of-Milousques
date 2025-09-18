// Package item définit tous les objets du jeu : ressources, armes, armures et potions
// Inclut leurs propriétés (attaque, défense, valeur) et les classes requises
package item

// ItemType définit les différents types d'objets disponibles
type ItemType string

const (
	TypeRessource ItemType = "ressource"
	TypeArme      ItemType = "arme"
	TypeCasque    ItemType = "casque"
	TypeTorse     ItemType = "torse"
	TypeJambiere  ItemType = "jambiere"
	TypePotion    ItemType = "potion"
	TypeSpecial   ItemType = "special"
)

type Item struct {
	Nom         string
	Type        ItemType
	Poids       int
	Effet       string
	Valeur      int
	Attaque     int    // Bonus d'attaque pour les armes
	Defense     int    // Bonus de défense pour les armures
	ClasseRequise string // Classe requise pour équiper ("" = toutes)
}

func NewItem(nom string) Item {
	switch nom {
	// === RESSOURCES DE BASE ===
	case "Bois":
		return Item{Nom: "Bois", Type: TypeRessource, Poids: 10, Effet: "Manger du bois vous fera mal aux dents", Valeur: 5}
	case "Fer":
		return Item{Nom: "Fer", Type: TypeRessource, Poids: 15, Effet: "Pas le meilleur matériaux pour fabriquer un lit", Valeur: 10}
	case "Blé":
		return Item{Nom: "Blé", Type: TypeRessource, Poids: 2, Effet: "Le meilleur atout pour rentrer à Ynuv", Valeur: 3}
	case "Laitue Vireuse":
		return Item{Nom: "Laitue Vireuse", Type: TypeRessource, Poids: 1, Effet: "La solution de secours favorite de Yelram Bob !", Valeur: 8}
	case "Pichon":
		return Item{Nom: "Pichon", Type: TypeRessource, Poids: 2, Effet: "Piche qui glisse n'amasse pas de risques !", Valeur: 12}
	
	// === CASQUES ===
	case "Casque en Cuir":
		return Item{Nom: "Casque en Cuir", Type: TypeCasque, Poids: 5, Effet: "Protection de tête en cuir souple", Valeur: 150, Defense: 5}
	case "Casque en Métal":
		return Item{Nom: "Casque en Métal", Type: TypeCasque, Poids: 8, Effet: "Protection de tête métallique résistante", Valeur: 300, Defense: 10}
	
	// === TORSES ===
	case "Torse en Cuir":
		return Item{Nom: "Torse en Cuir", Type: TypeTorse, Poids: 12, Effet: "Protection du torse en cuir souple", Valeur: 150, Defense: 5}
	case "Torse en Métal":
		return Item{Nom: "Torse en Métal", Type: TypeTorse, Poids: 20, Effet: "Protection du torse métallique résistante", Valeur: 300, Defense: 10}
	
	// === JAMBIÈRES ===
	case "Jambières en Cuir":
		return Item{Nom: "Jambières en Cuir", Type: TypeJambiere, Poids: 8, Effet: "Protection des jambes en cuir souple", Valeur: 150, Defense: 5}
	case "Jambières en Métal":
		return Item{Nom: "Jambières en Métal", Type: TypeJambiere, Poids: 15, Effet: "Protection des jambes métallique résistante", Valeur: 300, Defense: 10}
	
	// === ARMES SIMPLES ===
	case "Bâton Simple":
		return Item{Nom: "Bâton Simple", Type: TypeArme, Poids: 8, Effet: "Bâton de bois simple pour mage", Valeur: 250, Attaque: 10, ClasseRequise: "Mage"}
	case "Épée Simple":
		return Item{Nom: "Épée Simple", Type: TypeArme, Poids: 12, Effet: "Épée de fer simple pour guerrier", Valeur: 250, Attaque: 10, ClasseRequise: "Guerrier"}
	case "Dague Simple":
		return Item{Nom: "Dague Simple", Type: TypeArme, Poids: 6, Effet: "Dague acérée simple pour voleur", Valeur: 250, Attaque: 10, ClasseRequise: "Voleur"}
	
	// === ARMES D'EXPERT ===
	case "Bâton d'Expert":
		return Item{Nom: "Bâton d'Expert", Type: TypeArme, Poids: 15, Effet: "Bâton magique d'expert pour mage", Valeur: 500, Attaque: 20, ClasseRequise: "Mage"}
	case "Épée d'Expert":
		return Item{Nom: "Épée d'Expert", Type: TypeArme, Poids: 20, Effet: "Épée forgée d'expert pour guerrier", Valeur: 500, Attaque: 20, ClasseRequise: "Guerrier"}
	case "Dague d'Expert":
		return Item{Nom: "Dague d'Expert", Type: TypeArme, Poids: 10, Effet: "Dague empoisonnée d'expert pour voleur", Valeur: 500, Attaque: 20, ClasseRequise: "Voleur"}
	
	// === POTIONS ===
	case "Potion de Vie":
		return Item{Nom: "Potion de Vie", Type: TypePotion, Poids: 2, Effet: "Restaure 50 PV", Valeur: 50}
	case "Potion de Mana":
		return Item{Nom: "Potion de Mana", Type: TypePotion, Poids: 2, Effet: "Restaure 50 Mana", Valeur: 50}
	
	default:
		return Item{Nom: nom, Type: TypeSpecial, Poids: 10, Effet: "Objet mystérieux aux propriétés inconnues", Valeur: 10}
	}
}
