package character

import (
	"encoding/json"
	"fmt"
	"os"

	"world_of_milousques/classe"
	"world_of_milousques/inventory"
)

type Quete struct {
	Nom       string
	Accomplie bool
	Recompense string
}

type Character struct {
	Nom        string               `json:"nom"`
	Niveau     float64              `json:"niveau"`
	Pdv        int                  `json:"pdv"`
	Mana       int                  `json:"mana"`
	Classe     classe.Classe        `json:"classe"`
	Inventaire inventory.Inventaire `json:"inventaire"`
	Quetes     []Quete              `json:"quetes"`
}

func InitCharacter(nom string, c classe.Classe, niveau float64, pdv int, pdvmax int) Character {
	return Character{
		Nom:        nom,
		Niveau:     niveau,
		Pdv:        pdv,
		Mana:       c.ManaMax,
		Classe:     c,
		Inventaire: inventory.Inventaire{},
		Quetes:     []Quete{},
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

func (c *Character) ProposerEtAjouterQuete(nom string, recompense string) {
	c.Quetes = append(c.Quetes, Quete{Nom: nom, Accomplie: false, Recompense: recompense})
}

func (c *Character) CompleterQuete(nom string) {
	for i := range c.Quetes {
		if c.Quetes[i].Nom == nom {
			c.Quetes[i].Accomplie = true
			fmt.Println("Quête complétée :", nom)
			fmt.Println("Récompense :", c.Quetes[i].Recompense)
			break
		}
	}
}

func (c *Character) AfficherQuetes() {
	if len(c.Quetes) == 0 {
		fmt.Println("Aucune quête.")
		return
	}
	fmt.Println("Quêtes :")
	for _, q := range c.Quetes {
		status := "Non accomplie"
		if q.Accomplie {
			status = "Accomplie"
		}
		fmt.Printf("- %s : %s | Récompense : %s\n", q.Nom, status, q.Recompense)
	}
}
