package character

import (
	"world_of_milousques/inventory"
	"world_of_milousques/classe"
)

type Character struct {
	Nom string
	Niveau float64
	Pdv int
	PdvMax int
	Classe classe.Classe
	Inventaire inventory.Inventaire
}

func InitCharacter(nom string, classe classe.Classe, niveau float64, pdv, pdvmax int) Character {
	return Character{
		Nom: nom,
		Classe: classe,
		Niveau: niveau,
		Pdv: pdv,
		PdvMax: pdvmax,
		Inventaire: inventory.Inventaire{
			Potions: 3,
			Items: []string{},
		},
	}
}
