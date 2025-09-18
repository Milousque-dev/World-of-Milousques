package inventory

import (
	"fmt"
	"world_of_milousques/item"
)

type Inventaire struct {
	Potions     int         `json:"potions"`
	PotionsMana int         `json:"potions_mana"`
	Items       []item.Item `json:"items"`
}

func (inv *Inventaire) AddItem(it item.Item, quantity int) bool {
	espaceDisponible := 100 - len(inv.Items)
	if quantity > espaceDisponible {
		fmt.Printf("❌ Inventaire plein ! Vous ne pouvez ajouter que %d objets sur les %d demandés.\n", espaceDisponible, quantity)
		quantity = espaceDisponible
	}
	
	if quantity <= 0 {
		return false
	}
	
	for i := 0; i < quantity; i++ {
		inv.Items = append(inv.Items, it)
	}
	return true
}

func (inv *Inventaire) Recolter(ressources []item.Item) {
	if len(ressources) == 0 {
		fmt.Println("Aucune ressource à récolter ici.")
		return
	}

	espaceDisponible := 100 - len(inv.Items)
	if len(ressources) > espaceDisponible {
		fmt.Printf("⚠️  Votre inventaire ne peut contenir que %d objets supplémentaires.\n", espaceDisponible)
		fmt.Printf("Vous ne pouvez récolter que les %d premiers objets.\n", espaceDisponible)
		ressources = ressources[:espaceDisponible]
	}

	if len(ressources) == 0 {
		fmt.Println("❌ Inventaire plein ! Impossible de récolter quoi que ce soit.")
		return
	}

	fmt.Println("Vous récoltez :")
	for _, it := range ressources {
		fmt.Printf("- %s\n", it.Nom)
		inv.AddItem(it, 1)
	}
	fmt.Printf("✅ Votre inventaire contient maintenant %d/100 objets.\n", len(inv.Items))
}

func (inv *Inventaire) Afficher() {
	if len(inv.Items) == 0 {
		fmt.Printf("🎒 Votre inventaire est vide (0/100 objets).\n")
		return
	}
	fmt.Printf("🎒 === INVENTAIRE (%d/100 objets) === 🎒\n", len(inv.Items))
	for i, it := range inv.Items {
		fmt.Printf("%d) %s | Poids: %d | Effet: %s | Valeur: %d\n", i+1, it.Nom, it.Poids, it.Effet, it.Valeur)
	}
	
	if len(inv.Items) >= 90 {
		fmt.Printf("⚠️  Attention ! Votre inventaire est presque plein (%d/100).\n", len(inv.Items))
	} else {
		fmt.Printf("💼 Espace disponible : %d objets\n", 100-len(inv.Items))
	}
}
