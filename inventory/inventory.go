package inventory

import (
	"fmt"
	"world_of_milousques/item"
)

type Inventaire struct {
	Potions int         `json:"potions"`
	Items   []item.Item `json:"items"`
}

func (inv *Inventaire) AddItem(it item.Item, quantity int) {
	for i := 0; i < quantity; i++ {
		inv.Items = append(inv.Items, it)
	}
}

func (inv *Inventaire) Recolter(ressources []item.Item) {
	if len(ressources) == 0 {
		fmt.Println("Aucune ressource à récolter ici.")
		return
	}

	fmt.Println("Vous récoltez :")
	for _, it := range ressources {
		fmt.Printf("- %s\n", it.Nom)
		inv.AddItem(it, 1)
	}
	fmt.Printf("Votre inventaire contient maintenant %d objets.\n", len(inv.Items))
}

func (inv *Inventaire) Afficher() {
	if len(inv.Items) == 0 {
		fmt.Println("Votre inventaire est vide.")
		return
	}
	fmt.Println("Inventaire :")
	for i, it := range inv.Items {
		fmt.Printf("%d) %s | Poids: %d | Effet: %s | Valeur: %d\n", i+1, it.Nom, it.Poids, it.Effet, it.Valeur)
	}
	fmt.Printf("Total d'objets : %d\n", len(inv.Items))
}
