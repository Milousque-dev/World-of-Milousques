package places

import (
	"fmt"

	"world_of_milousques/character"
	"world_of_milousques/fight"
	"world_of_milousques/item"
	"world_of_milousques/ui"
	"world_of_milousques/utils"
)

type Scene struct {
	Titre       string
	Description string
	Options     []string
	Actions     []func(*character.Character)
	Ressources  []item.Item
}

func GetIntroDialogue() []Scene {
	return []Scene{
		{
			Titre: "Réveil mystérieux",
			Description: "??? : Réveille toi, aventurier...",
			Options: []string{"Qui êtes-vous ?", "Où suis-je ?", "Que sont les Milousques ?", "On peut y aller !"},
			Actions: []func(*character.Character){
				// Option 1 : Qui êtes-vous ?
				func(c *character.Character) {
					fmt.Println("\nMathiouw : Je suis Mathiouw, Le berger des jeunes âmes. Mon but est de faire de toi un aventurier assez puissant pour partir en quête des milousques.")
				},
				// Option 2 : Où suis-je ?
				func(c *character.Character) {
					fmt.Println("\nMathiouw : Tu es à Astrab, le lieu d'apparition des chasseurs de milousques comme toi.")
				},
				// Option 3 : Que sont les Milousques ?
				func(c *character.Character) {
					fmt.Println("\nMathiouw : Les milousques sont de puissantes chimères qui donne à ceux capable de les dompter un pouvoir incommensurable !")
				},
				// Option 4 : On peut y aller !
				func(c *character.Character) {
					// Cette action sera gérée différemment dans main.go
					fmt.Println("\nMathiouw : Parfait ! Commençons par un petit test de tes capacités...")
				},
			},
			Ressources: []item.Item{},
		},
	}
}

func GetTutorielCombat() (string, string, *fight.Ennemi) {
	quete := "Vaincre le Chacha Agressif"
	recompense := "1 potion"
	ennemi := &fight.Ennemi{
		Nom: "Chacha Agressif",
		Pv: 50,
		Attaque: 15,
	}
	return quete, recompense, ennemi
}

// ProposerQueteTutoriel propose la quête du tutoriel avec option de refus
func ProposerQueteTutoriel(c *character.Character) bool {
	fmt.Println("\n🎒 === QUÊTE PROPOSÉE === 🎒")
	fmt.Println("Mathiouw : Alors, acceptes-tu de m'aider à vaincre le Chacha Agressif ?")
	fmt.Println("Récompense : 1 potion de soin")
	
	options := []string{"Accepter la quête", "Refuser la quête"}
	ui.AfficherMenu("Décision", options)
	choix := utils.ScanChoice("Votre décision : ", options)
	
	if choix == 1 {
		fmt.Println("\nMathiouw : Parfait ! Je savais que je pouvais compter sur toi.")
		c.ProposerEtAjouterQueteAvecPNJ("Vaincre le Chacha Agressif", "1 potion", "Mathiouw")
		return true
	} else {
		fmt.Println("\nMathiouw : Le chacha a pris la quête à ta place, si tu le bats il empochera la récompense de sa défaite. Qu'il est malin ce chacha !")
		fmt.Println("\n(La quête continue quand même pour le tutoriel)")
		return false
	}
}

