package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ScanInt lit un entier depuis l'entrée standard avec validation et gestion d'erreur
func ScanInt(prompt string, min, max int) int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Erreur de lecture : %v. Réessayez.\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Println("Veuillez entrer une valeur.")
			continue
		}

		value, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("'%s' n'est pas un nombre valide. Réessayez.\n", input)
			continue
		}

		if value < min || value > max {
			fmt.Printf("Veuillez entrer un nombre entre %d et %d.\n", min, max)
			continue
		}

		return value
	}
}

// ScanString lit une chaîne de caractères depuis l'entrée standard avec validation
func ScanString(prompt string, minLength int) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Erreur de lecture : %v. Réessayez.\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if len(input) < minLength {
			if minLength == 1 {
				fmt.Println("Veuillez entrer au moins un caractère.")
			} else {
				fmt.Printf("Veuillez entrer au moins %d caractères.\n", minLength)
			}
			continue
		}

		return input
	}
}

// ScanChoice lit un choix parmi des options avec validation flexible (numéro ou texte)
func ScanChoice(prompt string, options []string) int {
	// Sécurité : vérifier que les options ne sont pas vides
	if len(options) == 0 {
		fmt.Println("Erreur : Aucune option disponible.")
		return 1 // Retourner 1 au lieu de 0 pour éviter les erreurs d'index
	}

	reader := bufio.NewReader(os.Stdin)
	attemptsCount := 0
	maxAttempts := 5 // Réduire le nombre d'essais
	
	for {
		attemptsCount++
		if attemptsCount > maxAttempts {
			fmt.Printf("\n⚠️  Trop de tentatives invalides (%d). Sélection automatique de l'option 1.\n", maxAttempts)
			fmt.Println("Appuyez sur Entrée pour continuer...")
			fmt.Scanln() // Pause pour que l'utilisateur voit le message
			return 1
		}
		
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Erreur de lecture : %v. Tentative %d/%d.\n", err, attemptsCount, maxAttempts)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Printf("Veuillez faire un choix. Tentative %d/%d.\n", attemptsCount, maxAttempts)
			continue
		}

		// Essayer d'abord comme un numéro
		if value, err := strconv.Atoi(input); err == nil {
			if value >= 1 && value <= len(options) {
				return value
			} else {
				fmt.Printf("⚠️  Numéro hors limites (%d). Choisissez entre 1 et %d. Tentative %d/%d.\n", value, len(options), attemptsCount, maxAttempts)
				continue
			}
		}

		// Ensuite essayer comme une correspondance textuelle
		inputUpper := strings.ToUpper(input)
		for i, option := range options {
			if strings.Contains(strings.ToUpper(option), inputUpper) {
				return i + 1
			}
		}

		fmt.Printf("⚠️  Choix invalide '%s'. Entrez un numéro (1-%d) ou le début d'une option. Tentative %d/%d.\n", input, len(options), attemptsCount, maxAttempts)
	}
}
