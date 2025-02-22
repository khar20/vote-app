package handlers

import (
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"api-server/internal/database"
	"api-server/internal/models"
	"api-server/pkg/config"
)

var mu sync.Mutex

// Obtiene todos los votos almacenados en la base de datos.
func GetVotesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		voteBlocks, err := database.GetAllVotes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron recuperar los votos", "detalles": err.Error()})
			return
		}

		c.JSON(http.StatusOK, voteBlocks)
	}
}

// Permite enviar un voto y almacenarlo en la base de datos.
func PostVoteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var vote models.Vote

		if err := c.ShouldBindJSON(&vote); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de solicitud inválido", "detalles": err.Error()})
			return
		}

		if _, err := time.Parse(time.RFC3339, vote.Timestamp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de marca de tiempo inválido", "detalles": err.Error()})
			return
		}

		mu.Lock()
		defer mu.Unlock()

		if err := database.StoreVote(vote); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo almacenar el voto", "detalles": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"mensaje": "Voto enviado con éxito"})
	}
}

// Realiza el conteo de votos utilizando criptografía homomórfica.
func TallyVotesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()

		voteBlocks, err := database.GetAllVotes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron recuperar los votos", "detalles": err.Error()})
			return
		}

		if len(voteBlocks) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No hay votos para contar"})
			return
		}

		sum := make([]*big.Int, len(voteBlocks[0].Vote))
		for i := range sum {
			sum[i] = big.NewInt(1)
		}

		for _, voteBlock := range voteBlocks {
			for i, v := range voteBlock.Vote {
				encryptedValue := new(big.Int)
				if _, success := encryptedValue.SetString(v, 10); !success {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Formato de voto cifrado inválido"})
					return
				}
				sum[i].Mul(sum[i], encryptedValue)
			}
		}

		result := make([]*big.Int, len(sum))
		for i, s := range sum {
			decrypted := config.PrivateKey.Decrypt(s)
			result[i] = decrypted
		}

		c.JSON(http.StatusOK, gin.H{"resultado": result})
	}
}

// Devuelve la clave pública utilizada para cifrar los votos.
func GetPublicKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		publicKey := config.PublicKey
		c.JSON(http.StatusOK, gin.H{
			"N": publicKey.N.String(),
			"G": publicKey.G.String(),
		})
	}
}

// Verifica la integridad de la cadena de bloques de votos.
func CheckChainHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		blocks, err := database.GetAllVotes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar los bloques", "detalles": err.Error()})
			return
		}

		for i := 1; i < len(blocks); i++ {
			if blocks[i].PrevHash != blocks[i-1].Hash {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":            "Fallo en la verificación de integridad de la cadena de bloques",
					"numero_de_bloque": i,
					"detalles":         "Hashes no coinciden entre bloques",
				})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"mensaje": "Todos los bloques están correctamente enlazados"})
	}
}
