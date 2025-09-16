package classe

import (
	"world_of_milousques/classeitem"
	"world_of_milousques/sorts"
)

type Classe struct {
	Nom         string                  `json:"nom"`
	Pvmax       int                     `json:"pv_max"`
	ManaMax     int                     `json:"mana_max"`
	Sorts       []sorts.Sorts           `json:"sorts"`
	ClasseItems []classeitem.Classeitem `json:"classe_items"`
}

func GetClasse(nom string) Classe {
	switch nom {
	case "Guerrier":
		return Classe{
			Nom:     "Guerrier",
			Pvmax:   130,
			ManaMax: 70,
			Sorts: []sorts.Sorts{
				sorts.GetSorts("Fracasser"),
				sorts.GetSorts("Briser"),
			},
			ClasseItems: []classeitem.Classeitem{},
		}
	case "Mage":
		return Classe{
			Nom:     "Mage",
			Pvmax:   70,
			ManaMax: 130,
			Sorts: []sorts.Sorts{
				sorts.GetSorts("Boule de feu"),
				sorts.GetSorts("Explosion"),
			},
			ClasseItems: []classeitem.Classeitem{},
		}
	case "Voleur":
		return Classe{
			Nom:     "Voleur",
			Pvmax:   100,
			ManaMax: 100,
			Sorts: []sorts.Sorts{
				sorts.GetSorts("Coup bas"),
				sorts.GetSorts("Fourberie"),
			},
			ClasseItems: []classeitem.Classeitem{},
		}
	default:
		return Classe{
			Nom:         nom,
			Pvmax:       100,
			ManaMax:     100,
			Sorts:       []sorts.Sorts{},
			ClasseItems: []classeitem.Classeitem{},
		}
	}
}

func GetClassesDisponibles() []Classe {
	classesNoms := []string{"Guerrier", "Mage", "Voleur"}
	var classes []Classe
	for _, nom := range classesNoms {
		classes = append(classes, GetClasse(nom))
	}
	return classes
}
