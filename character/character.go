package character

import (
	"encoding/json"
	"os"
	"fmt"

	"world_of_milousques/inventory"
	"world_of_milousques/classe"
)

type Character struct {
	Nom        string             `json:"nom"`
	Niveau     float64            `json:"niveau"`
	Pdv        int                `json:"pdv"`
	PdvMax     int                `json:"pdv_max"`
	Classe     classe.Classe      `json:"classe"`
	Inventaire inventory.Inventaire `json:"inventaire"`
}

func InitCharacter(nom string, c classe.Classe, niveau float64, pdv int, pdvmax int) Character {
	return Character{
		Nom:        nom,
		Classe:     c,
		Niveau:     niveau,
		Pdv:        pdv,
		PdvMax:     pdvmax,
		Inventaire: inventory.Inventaire{},
	}
}

func (c *Character) Sauvegarder() error {
	filename := "saves/" + c.Nom + ".json"

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(c)
	if err != nil {
		return err
	}

	fmt.Println("Personnage sauvegardé dans", filename)
	return nil
}

func Charger(nom string) (*Character, error) {
	filename := "saves/" + nom + ".json"

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var c Character
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		return nil, err
	}

	fmt.Println("Personnage chargé depuis", filename)
	return &c, nil
}
