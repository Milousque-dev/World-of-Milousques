package classe

import (
	"world_of_milousques/sorts"
	"world_of_milousques/classeitem"
)

type Classe struct {
	Nom string `json:"nom"`
	Pvmax int `json:"pv_max"`
	Sorts sorts.Sorts `json:"sorts"`
	Classeitem classeitem.Classeitem `json:"classeitem"`
}

func GetClasse(nom string) Classe {
	switch nom {
	case "Guerrier":
		return Classe{
			Nom: "Guerrier",
			Pvmax: 130,
			Sorts: sorts.Sorts{},
			Classeitem: classeitem.Classeitem{},
		}
	case "Mage":
		return Classe{
			Nom: "Mage",
			Pvmax: 70,
			Sorts: sorts.Sorts{},
			Classeitem: classeitem.Classeitem{},
		}
	case "Voleur":
		return Classe{
			Nom: "Voleur",
			Pvmax: 100,
			Sorts: sorts.Sorts{},
			Classeitem: classeitem.Classeitem{},
		}
	default:
		return Classe{
			Nom: nom,
			Pvmax: 100,
			Sorts: sorts.Sorts{},
			Classeitem: classeitem.Classeitem{},
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
