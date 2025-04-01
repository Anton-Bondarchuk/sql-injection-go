package storage

import (
	"context"
	"fmt"
	"sql-injection-go/internal/storage"

	"sql-injection-go/internal/domain/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Conn *pgxpool.Pool
}

func New(ctx context.Context, databaseUrl string) (*Storage, error) {
	const op = "storage.postgres.new"
	conn, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{Conn: conn}, nil
}

func (s *Storage) Close() {
	s.Conn.Close()
}


func (s* Storage) GetStudentsSafe(ctx context.Context, id int) ([]models.Student, error) {
	const op = "storage.get_students_safe"
	query := "SELECT id, age, sex, card_id, name FROM students WHERE id = $1"
	
	rows, err := s.Conn.Query(ctx, query, id)
	if err != nil {
		// var pgErr pgconn.PgError
		// if errors.As(err, &pgErr) && pgErr.Code == pgconn
		return nil, fmt.Errorf("%s: query: %w", op, err)
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var st models.Student
		if err := rows.Scan(&st.Id, &st.Age, &st.Sex, &st.CardId, &st.Name); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, storage.ErrParsingQuery)
		}
		students = append(students, st)
	}
	
	return students, nil
}

func (s* Storage) GetStudentInjection(ctx context.Context, id string) ([]models.Student, error) {
	// TODO: add ability to explode someone else injection types
	const op = "storage.get_students_injection"
	query := fmt.Sprintf("SELECT id, age, sex, card_id, name FROM students WHERE id = %s", id)
	rows, err := s.Conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: query: %w", op, err)
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var st models.Student
		if err := rows.Scan(&st.Id, &st.Age, &st.Sex, &st.CardId, &st.Name); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		students = append(students, st)
	}
	
	return students, nil
}