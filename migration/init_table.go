package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Получаем строку подключения к БД
	connString := os.Getenv("DATABASE_URL")

	if connString == "" {
		log.Fatal("DATABASE_URL не найден в .env")
	}
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Невозможно подключиться к базе данных: %v", err)
	}
	defer conn.Close(context.Background())

	// todo: convert to must load 
	// todo: write personnaly file for each migration
	sqlBytesInitTable, err := os.ReadFile("./migration/init_table.sql")
	if err != nil {
		log.Fatalf("Ошибка чтения файла миграции: %v", err)
	}
	sqlBytesSeeds, err := os.ReadFile("./migration/seeds.sql")
	if err != nil {
		log.Fatalf("Ошибка чтения файла миграции: %v", err)
	}
	sqlScriptInitTable := string(sqlBytesInitTable)
	sqlScriptSeeds := string(sqlBytesSeeds)

	_, err = conn.Exec(context.Background(), sqlScriptInitTable)
	op := "sql init table"
	if err != nil {
		log.Fatalf("Ошибка выполнения миграции: " + op + "%v", err)
	}

	fmt.Println("Миграция " + op + " успешно выполнены!")

	_, err = conn.Exec(context.Background(), sqlScriptSeeds)
	op = "seeds"
	if err != nil {
		log.Fatalf("Ошибка выполнения миграции: " + op + "%v", err)
	}

	fmt.Println("Миграция " + op + " успешно выполнены!")
}