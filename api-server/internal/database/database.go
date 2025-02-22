package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"api-server/internal/models"
)

var Db *sql.DB

// Inicializa la base de datos SQLite y crea la tabla de votos si no existe.
func InitDB() {
	var err error
	Db, err = sql.Open("sqlite3", "./votes.db")
	if err != nil {
		log.Fatalf("No se pudo abrir la base de datos SQLite: %v", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS votes (
		hash TEXT PRIMARY KEY,
		data JSON
	);
	`
	_, err = Db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("No se pudo crear la tabla de votos: %v", err)
	}

	fmt.Println("¡Base de datos SQLite inicializada con éxito!")
}

// Cierra la conexión con la base de datos.
func CloseDB() {
	if err := Db.Close(); err != nil {
		log.Fatalf("No se pudo cerrar la base de datos: %v", err)
	}
}

// Guarda un voto en la base de datos con un hash enlazado al anterior.
func StoreVote(vote models.Vote) error {
	lastVoteBlock, err := GetLastVote()
	prevHash := ""

	if err == nil {
		prevHash = lastVoteBlock.Hash
	}

	voteBlock := models.VoteBlock{
		Vote:      vote.Vote,
		Timestamp: vote.Timestamp,
		PrevHash:  prevHash,
	}
	voteBlock.Hash = voteBlock.ComputeHash()

	voteData, err := json.Marshal(voteBlock)
	if err != nil {
		return fmt.Errorf("fallo al serializar el bloque de voto: %w", err)
	}

	insertQuery := `INSERT INTO votes (hash, data) VALUES (?, ?)`
	_, err = Db.Exec(insertQuery, voteBlock.Hash, voteData)
	if err != nil {
		return fmt.Errorf("fallo al almacenar el bloque de voto: %w", err)
	}

	return nil
}

// Recupera un bloque de voto por su hash.
func GetVote(hash string) (*models.VoteBlock, error) {
	var voteData []byte

	query := `SELECT data FROM votes WHERE hash = ?`
	err := Db.QueryRow(query, hash).Scan(&voteData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bloque de voto no encontrado")
		}
		return nil, fmt.Errorf("fallo al recuperar el bloque de voto: %w", err)
	}

	var voteBlock models.VoteBlock
	err = json.Unmarshal(voteData, &voteBlock)
	if err != nil {
		return nil, fmt.Errorf("fallo al deserializar el bloque de voto: %w", err)
	}

	return &voteBlock, nil
}

// Recupera todos los bloques de votos de la base de datos.
func GetAllVotes() ([]models.VoteBlock, error) {
	query := `SELECT data FROM votes`
	rows, err := Db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("fallo al recuperar los votos: %w", err)
	}
	defer rows.Close()

	var voteBlocks []models.VoteBlock
	for rows.Next() {
		var voteData []byte
		if err := rows.Scan(&voteData); err != nil {
			return nil, fmt.Errorf("fallo al escanear los datos del voto: %w", err)
		}

		var voteBlock models.VoteBlock
		if err := json.Unmarshal(voteData, &voteBlock); err != nil {
			return nil, fmt.Errorf("fallo al deserializar el bloque de voto: %w", err)
		}

		voteBlocks = append(voteBlocks, voteBlock)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre las filas: %w", err)
	}

	return voteBlocks, nil
}

// Recupera el último bloque de voto almacenado.
func GetLastVote() (*models.VoteBlock, error) {
	query := `SELECT data FROM votes ORDER BY rowid DESC LIMIT 1`
	var voteData []byte
	err := Db.QueryRow(query).Scan(&voteData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no se encontraron votos")
		}
		return nil, fmt.Errorf("fallo al recuperar el último bloque de voto: %w", err)
	}

	var voteBlock models.VoteBlock
	err = json.Unmarshal(voteData, &voteBlock)
	if err != nil {
		return nil, fmt.Errorf("fallo al deserializar el bloque de voto: %w", err)
	}

	return &voteBlock, nil
}
