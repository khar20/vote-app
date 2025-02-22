package config

import (
	"log"

	"api-server/internal/utils"
)

var (
	PrivateKey utils.PrivateKey
	PublicKey  utils.PublicKey
)

// LoadConfig carga la configuración desde el archivo .env y genera las claves Paillier.
func LoadConfig() {
	// Generar claves Paillier
	publicKey, privateKey, err := utils.GeneratePaillierKeys(1024)
	if err != nil {
		log.Fatal("Error al generar las claves Paillier")
	}
	PrivateKey = *privateKey
	PublicKey = *publicKey

	log.Println("Configuración cargada correctamente")
}
