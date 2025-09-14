package sorts

type Sorts struct {
	Nom string
	Degats int
	Cout int
}

func InitSorts(nom string, degats int, cout int) Sorts {
	return Sorts{
		Nom: nom,
		Degats: degats,
		Cout: cout,
	}
}

func GetSorts(nom string) Sorts{
	switch nom {
	case "Boule de feu":
		return Sorts{
			Nom: "Boule de feu",
			Degats: 30,
			Cout: 20,
		}
	case "Coup bas":
		return Sorts{
			Nom: "Coup bas",
			Degats: 25,
			Cout: 15,
		}
	case "Fracasser":
		return Sorts{
			Nom: "Fracasser",
			Degats: 20,
			Cout: 10,
		}
	default:
		return Sorts{
			Nom: nom,
			Degats: 25,
			Cout: 15,
		}
	}
}
