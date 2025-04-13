package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"new_project/internal/models"
)

//type Animal struct {
//	Id          int       `json:"id"`
//	Name        string    `json:"name"`
//	Species     string    `json:"species"`
//	Type        string    `json:"type"`
//	Habitat     string    `json:"habitat"`
//	Image       string    `json:"image"`
//	Description string    `json:"description"`
//	DietType    string    `json:"dietType"`
//	Lifespan    string    `json:"lifespan"`
//	FunFact     string    `json:"funFact"`
//	CreatedAt   time.Time `json:"created_at"`
//	UpdatedAt   time.Time `json:"updated_at"`
//}

type Animal models.Animal

type AnimalsResponse struct {
	Animals     []Animal `json:"animals"`
	TotalRows   int      `json:"total_rows"`
	CurrentPage int      `json:"current_page"`
	TotalPages  int      `json:"total_pages"`
}

func (s *service) GetAnimalById(id string) (*Animal, error) {
	row := s.db.QueryRowContext(context.Background(), `
		SELECT id, name, species, type, habitat, image, description,diet_type,lifespan,fun_fact
		FROM animals
		WHERE id = $1`, id)

	var animalResp Animal
	err := row.Scan(&animalResp.Id, &animalResp.Name, &animalResp.Species, &animalResp.Type, &animalResp.Habitat, &animalResp.Image, &animalResp.Description, &animalResp.DietType, &animalResp.Lifespan, &animalResp.FunFact)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("animals not found")
		}
		return nil, err
	}

	return &animalResp, nil
}

// AddShortenedUrl - Add shortened url into the database.
func (s *service) InsertAnimal(animalResp *Animal) error {
	_, err := s.db.ExecContext(context.Background(), "INSERT INTO animals (name, species, type, habitat, image, description,diet_type,lifespan,fun_fact) VALUES ($1, $2, $3, $4 , $5 ,$6,$7,$8,$9 )",
		&animalResp.Name, &animalResp.Species, &animalResp.Type, &animalResp.Habitat, &animalResp.Image, &animalResp.Description, &animalResp.DietType, &animalResp.Lifespan, &animalResp.FunFact)
	return err
}

// DeleteByID deletes a record by its ID.
func (s *service) DeleteByID(id int) error {
	_, err := s.db.ExecContext(context.Background(), "DELETE FROM animals WHERE id = ?", id)
	return err
}

func (s *service) GetAllAnimals(page, pageSize int) (*AnimalsResponse, error) {
	// Validate and set default pagination values
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get total count first
	var total int
	err := s.db.QueryRowContext(context.Background(), `
        SELECT COUNT(*)
        FROM animals
    `).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error counting animals: %w", err)
	}

	// Calculate total pages
	pages := (total + pageSize - 1) / pageSize

	// Get paginated results
	rows, err := s.db.QueryContext(context.Background(), `
        SELECT id, name, species, type, habitat, image, description ,diet_type,lifespan,fun_fact
        FROM animals
        ORDER BY id
        LIMIT $1 OFFSET $2
    `, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying animals: %w", err)
	}
	defer rows.Close()

	var animals []Animal
	for rows.Next() {
		var animal Animal
		err := rows.Scan(
			&animal.Id,
			&animal.Name,
			&animal.Species,
			&animal.Type,
			&animal.Habitat,
			&animal.Image,
			&animal.Description,
			&animal.DietType,
			&animal.Lifespan,
			&animal.FunFact,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning animal row: %w", err)
		}
		animals = append(animals, animal)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating animal rows: %w", err)
	}

	return &AnimalsResponse{
		Animals:     animals,
		TotalRows:   total,
		CurrentPage: page,
		TotalPages:  pages,
	}, nil
}
