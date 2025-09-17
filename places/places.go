package places

import (
	"fmt"
	"math/rand"
	"time"

	"world_of_milousques/character"
	"world_of_milousques/fight"
	"world_of_milousques/item"
)

type Scene struct {
	Titre string
	Description string
	Options []string
	Actions []func(*character.Character)
	Ressources []item.Item
}

func GetIntroDialogue() []Scene {
	return []Scene{
		{
			Titre: "Réveil mystérieux",
			Description: "??? : Réveille toi, aventurier...",
			Options: []string{"Qui êtes-vous ?", "Où suis-je ?", "Je suis prêt, on peut y aller !"},
			Actions: []func(*character.Character){
				func(c *character.Character) { fmt.Println("Mathiouw : Je suis Mathiouw, le berger des jeunes âmes.") },
				func(c *character.Character) { fmt.Println("Mathiouw : Un peu long à expliquer, cela viendra plus tard.") },
				func(c *character.Character) { fmt.Println("Mathiouw : On va commencer le tutoriel combat.") },
			},
			Ressources: []item.Item{},
		},
	}
}

func GetTutorielCombat() (string, string, fight.Ennemi) {
	quete := "Vaincre le Chacha Agressif"
	recompense := "1 potion"
	ennemi := fight.Ennemi{
		Nom: "Chacha Agressif",
		Pv: 50,
		Attaque: 15,
	}
	return quete, recompense, ennemi
}

func GetScenesAventure() []Scene {
	rand.Seed(time.Now().UnixNano())

	return []Scene{
		{
			Titre: "Astrab",
			Description: "Vous arrivez à Astrab, le village aux canards",
			Options: []string{"Parler au maire", "Explorer le marché", "Récolter des ressources"},
			Actions: []func(*character.Character){
				func(c *character.Character) { fmt.Println("Vous discutez avec le maire... (à compléter)") },
				func(c *character.Character) { fmt.Println("Vous explorez le marché... (à compléter)") },
				func(c *character.Character) { c.Inventaire.Recolter([]item.Item{item.NewItem("Blé"), item.NewItem("Laitue Vireuse")}) },
			},
			Ressources: []item.Item{item.NewItem("Blé"), item.NewItem("Laitue Vireuse")},
		},
		{
			Titre: "Forêt mystérieuse",
			Description: "La forêt est dense et silencieuse...",
			Options: []string{"Chemin de gauche", "Chemin de droite", "Récolter des ressources"},
			Actions: []func(*character.Character){
				func(c *character.Character) { fmt.Println("Vous prenez le chemin de gauche... (à compléter)") },
				func(c *character.Character) { fmt.Println("Vous prenez le chemin de droite... (à compléter)") },
				func(c *character.Character) { c.Inventaire.Recolter([]item.Item{item.NewItem("Bois"), item.NewItem("Pichon")}) },
			},
			Ressources: []item.Item{item.NewItem("Bois"), item.NewItem("Pichon")},
		},
		{
			Titre: "Mine abandonnée",
			Description: "Une vieille mine semble abandonnée, il pourrait y avoir du minerai...",
			Options: []string{"Entrer dans la mine", "Faire demi-tour", "Récolter des ressources"},
			Actions: []func(*character.Character){
				func(c *character.Character) { fmt.Println("Vous explorez la mine... (à compléter)") },
				func(c *character.Character) { fmt.Println("Vous retournez sur le chemin... (à compléter)") },
				func(c *character.Character) { c.Inventaire.Recolter([]item.Item{item.NewItem("Fer")}) },
			},
			Ressources: []item.Item{item.NewItem("Fer")},
		},
	}
}
