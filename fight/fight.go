// Package fight gère le système de combat tour par tour
// Inclut l'utilisation de sorts, potions et le calcul des bonus d'équipement
package fight

import (
	"fmt"
	"world_of_milousques/character"
	"world_of_milousques/ui"
	"world_of_milousques/utils"
)

type Ennemi struct {
	Nom     string
	Pv      int
	Attaque int
}

func Fight(joueur *character.Character, ennemi *Ennemi) {
	tourCount := 0
	maxTours := 100 // Limite le nombre de tours pour éviter les combats infinis
	
	for joueur.Pdv > 0 && ennemi.Pv > 0 && tourCount < maxTours {
		tourCount++
		
		fmt.Printf("\n=== Tour %d ===\n", tourCount)
		
		ui.AfficherMenuCombat(
			joueur.Nom, joueur.Pdv, joueur.Classe.Pvmax, joueur.Mana, joueur.Classe.ManaMax,
			ennemi.Nom, ennemi.Pv, joueur.Classe.Sorts, joueur.Inventaire.Potions, joueur.Inventaire.PotionsMana,
		)

		// Créer les options disponibles
		options := make([]string, 0)
		
		// Vérifier si le joueur a des sorts utilisables
		sortUtilisable := false
		for _, s := range joueur.Classe.Sorts {
			options = append(options, fmt.Sprintf("%s (Dégâts: %d, Mana: %d)", s.Nom, s.Degats, s.Cout))
			if joueur.Mana >= s.Cout {
				sortUtilisable = true
			}
		}
		
		// Ajouter les options d'utilisation de potions et de fuite
		options = append(options, fmt.Sprintf("Utiliser une potion de vie (+50 PV) (%d disponibles)", joueur.Inventaire.Potions))
		options = append(options, fmt.Sprintf("Utiliser une potion de mana (+50 Mana) (%d disponibles)", joueur.Inventaire.PotionsMana))
		options = append(options, "Fuir le combat")
		
		// Si aucun sort n'est utilisable et pas de potions, proposer la fuite
		if !sortUtilisable && joueur.Inventaire.Potions <= 0 && joueur.Inventaire.PotionsMana <= 0 {
			fmt.Println("\n⚠️  Plus de mana et pas de potions ! Vous devez fuir ou utiliser une attaque de base.")
		}
		
		choix := utils.ScanChoice("Choisis ton action : ", options)

		optionPotionVie := len(joueur.Classe.Sorts) + 1
		optionPotionMana := len(joueur.Classe.Sorts) + 2
		optionFuite := len(joueur.Classe.Sorts) + 3

		if choix == optionFuite {
			fmt.Println("\n🏃 Vous fuyez le combat !")
			break
		} else if choix == optionPotionVie {
			if joueur.Inventaire.Potions > 0 {
				anciensPV := joueur.Pdv
				joueur.Pdv += 50
				if joueur.Pdv > joueur.Classe.Pvmax {
					joueur.Pdv = joueur.Classe.Pvmax
				}
				joueur.Inventaire.Potions--
				pvRecuperes := joueur.Pdv - anciensPV
				fmt.Printf("🧆 Vous utilisez une potion de vie et récupérez %d PV !\n", pvRecuperes)
			} else {
				fmt.Println("⚠️  Vous n'avez pas de potion de vie !")
				continue
			}
		} else if choix == optionPotionMana {
			if joueur.Inventaire.PotionsMana > 0 {
				ancienMana := joueur.Mana
				joueur.Mana += 50
				if joueur.Mana > joueur.Classe.ManaMax {
					joueur.Mana = joueur.Classe.ManaMax
				}
				joueur.Inventaire.PotionsMana--
				manaRecupere := joueur.Mana - ancienMana
				fmt.Printf("🧙 Vous utilisez une potion de mana et récupérez %d Mana !\n", manaRecupere)
			} else {
				fmt.Println("⚠️  Vous n'avez pas de potion de mana !")
				continue
			}
		} else if choix >= 1 && choix <= len(joueur.Classe.Sorts) {
			s := joueur.Classe.Sorts[choix-1]
			if joueur.Mana < s.Cout {
				fmt.Println("⚠️  Pas assez de mana pour lancer ce sort !")
				continue
			}
			joueur.Mana -= s.Cout
			
			// Appliquer bonus d'attaque de l'équipement
			bonusAttaque := joueur.CalculerAttaqueBonus()
			degatsFinaux := s.Degats + bonusAttaque
			ennemi.Pv -= degatsFinaux
			
			if bonusAttaque > 0 {
				fmt.Printf("⚔️  Tu lances %s et infliges %d dégâts (%d base + %d bonus équipement) !\n", s.Nom, degatsFinaux, s.Degats, bonusAttaque)
			} else {
				fmt.Printf("⚔️  Tu lances %s et infliges %d dégâts !\n", s.Nom, degatsFinaux)
			}
		} else {
			fmt.Println("⚠️  Choix invalide, vous perdez votre tour !")
		}

		if ennemi.Pv <= 0 {
			fmt.Printf("🏆 %s est vaincu !\n", ennemi.Nom)
			break
		}

		// Appliquer bonus de défense
		bonusDefense := joueur.CalculerDefenseBonus()
		degatsSubis := ennemi.Attaque - bonusDefense
		if degatsSubis < 1 {
			degatsSubis = 1 // Minimum 1 dégât
		}
		
		joueur.Pdv -= degatsSubis
		
		if bonusDefense > 0 {
			fmt.Printf("🔴 %s t'attaque ! Tu subis %d dégâts (%d - %d défense) !\n", ennemi.Nom, degatsSubis, ennemi.Attaque, bonusDefense)
		} else {
			fmt.Printf("🔴 %s t'attaque et inflige %d dégâts !\n", ennemi.Nom, degatsSubis)
		}
		
		// Petite pause pour la lisibilité
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
	}
	
	// Vérifier si le combat s'est arrêté à cause de la limite de tours
	if tourCount >= maxTours {
		fmt.Printf("\n⚠️  Combat trop long (%d tours) ! Arrêt automatique.\n", maxTours)
		fmt.Println("Le combat se termine par un match nul...")
		return
	}

	if joueur.Pdv > 0 && ennemi.Pv <= 0 {
		
		// Mettre à jour le progrès des quêtes après victoire
		joueur.MettreAJourProgresQuete(ennemi.Nom)
		
		// Gagner XP après victoire (calculer sur les PV originaux si l'ennemi est mort)
		pvOriginaux := ennemi.Pv
		if ennemi.Pv <= 0 {
			// Estimer les PV originaux si l'ennemi est mort
			pvOriginaux = ennemi.Attaque * 3 // Estimation basique
		}
		xpGagne := 25 + (pvOriginaux / 2)
		joueur.GagnerExperience(xpGagne)
	} else if joueur.Pdv <= 0 {
		fmt.Println("💀 Tu as été vaincu... Game Over.")
	} else {
		fmt.Println("🏃 Vous avez fui le combat avec succès !")
	}
}
