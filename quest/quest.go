package quest

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Quest struct {
	Titre       string
	Description string
	EstComplete bool
	Recompense  string
}

func NewQuest(titre, description, recompense string) Quest {
	return Quest{
		Titre:       titre,
		Description: description,
		EstComplete: false,
		Recompense:  recompense,
	}
}

func (q *Quest) Complete() {
	q.EstComplete = true
	fmt.Println("Quête terminée :", q.Titre)
	fmt.Println("Récompense :", q.Recompense)
}

func (q *Quest) Afficher() {
	status := "En cours"
	if q.EstComplete {
		status = "Terminée"
	}
	fmt.Println("Quête :", q.Titre)
	fmt.Println("Description :", q.Description)
	fmt.Println("Statut :", status)
	fmt.Println("Récompense :", q.Recompense)
	fmt.Println("----------------------")
}

func ProposerQuete(q Quest) bool {
	fmt.Println("\nNouvelle quête proposée :", q.Titre)
	fmt.Println(q.Description)
	fmt.Println("Récompense :", q.Recompense)
	fmt.Println("Voulez-vous accepter cette quête ? (O/N)")

	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		input = strings.ToUpper(strings.TrimSpace(input))
		if input == "O" {
			fmt.Println("Quête acceptée :", q.Titre)
			return true
		} else if input == "N" {
			fmt.Println("Quête refusée :", q.Titre)
			return false
		} else {
			fmt.Println("Entrée invalide, tapez O pour accepter ou N pour refuser.")
		}
	}
}
