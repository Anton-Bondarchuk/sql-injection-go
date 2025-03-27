package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Student struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Sex    bool   `json:"sex"`
	CardId int    `json:"card_id"`
}

type CardCredit struct {
	Id         int `json:"id"`
	StudentId  int `json:"student_id"`
	CardNumber int `json:"card_number"`
	Expiration int `json:"expiration"`
	Cvv        int `json:"cvv"`
}

type StudentData struct {
	Student    Student    `json:"student"`
	CardCredit CardCredit `json:"card_credit,omitempty"`
}

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
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	studentHandler := New(conn)

	router := gin.Default()
	// Указываем директорию для шаблонов в папке "templates"
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/search")
	})

	// Уязвимый endpoint – поиск по имени, позволяющий злоумышленнику выполнить SQL-инъекцию.
	// Если параметр "query" содержит подстроку "join" (без учета регистра), то:
	// 1. Если строка начинается с "*", то будет использован предоставленный злоумышленником список полей.
	// 2. Иначе будет выполнен запрос с предопределенным списком полей из таблиц students и card_credits.
	router.GET("/students", studentHandler.GetStudentWithSqlInjection)
	// Безопасный endpoint, использующий параметризованный запрос и НЕ отдающий данные card_credits.
	router.GET("/students_safe", studentHandler.GetStudentsSafe)
	// Рендеринг страницы поиска
	router.GET("/search", studentHandler.RenderSearch)

	router.Run("0.0.0.0:8080")
}

type StudentHandler struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *StudentHandler {
	return &StudentHandler{conn: conn}
}

// GetStudentWithSqlInjection осуществляет поиск студентов по имени с использованием небезопасного конкатенирования SQL-запроса.
// Если в параметре "query" присутствует подстрока "join", то злоумышленник может самостоятельно добавить JOIN на таблицу card_credits
// и указать поля для выборки (либо выбрать все, используя "*").
// Примеры использования:
//   1. Злоумышленник может передать в качестве query:
//         "JOIN card_credits cc ON cc.student_id = s.id WHERE TRUE"
//      в этом случае запрос будет выполнен как:
//         SELECT s.id, s.age, s.sex, s.card_id, s.name, cc.student_id, cc.id, cc.card_number, cc.expiration, cc.cvv FROM students s JOIN card_credits cc ON cc.student_id = s.id WHERE s.name LIKE '%John%'
//   2. Или злоумышленник может передать:
//         "* FROM students s JOIN card_credits cc ON cc.student_id = s.id WHERE s.name LIKE '%John%'"
//      тогда итоговый запрос будет:
//         SELECT * FROM students s JOIN card_credits cc ON cc.student_id = s.id WHERE s.name LIKE '%John%'"
func (h *StudentHandler) GetStudentWithSqlInjection(c *gin.Context) {
	queryParam := c.Query("query")

	var sqlQuery string
	// Если параметр содержит "join" (без учета регистра), позволяем злоумышленнику добавить JOIN и поля для выборки.
	if strings.Contains(strings.ToLower(queryParam), "join") {
		trimmed := strings.TrimSpace(queryParam)
		if strings.HasPrefix(trimmed, "*") {
			// Пользователь выбрал все поля
			sqlQuery = "SELECT " + trimmed
		} else {
			// Используем предопределенный список полей из таблицы students и card_credits
			sqlQuery = "SELECT s.id, s.age, s.sex, s.card_id, s.name, cc.student_id, cc.id, cc.card_number, cc.expiration, cc.cvv FROM students s " + queryParam
		}
		rows, err := h.conn.Query(context.Background(), sqlQuery)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error(), "query": sqlQuery})
			return
		}
		defer rows.Close()

		var studentData []StudentData
		for rows.Next() {
			// Подразумевается, что злоумышленник добавил данные для card_credits,
			// поэтому пытаемся сканировать их
			var s Student
			var cc CardCredit
			err := rows.Scan(&s.Id, &s.Age, &s.Sex, &s.CardId, &s.Name, &cc.StudentId, &cc.Id, &cc.CardNumber, &cc.Expiration, &cc.Cvv)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error(), "query": sqlQuery})
				return
			}
			studentData = append(studentData, StudentData{
				Student:    s,
				CardCredit: cc,
			})
		}
		c.JSON(200, studentData)
		return
	}

	sqlQuery = fmt.Sprintf("SELECT * FROM students WHERE id = %s", queryParam)
	rows, err := h.conn.Query(context.Background(), sqlQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "query": sqlQuery})
		return
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.Id, &s.Name, &s.Age, &s.Sex, &s.CardId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "query": sqlQuery})
			return
		}
		students = append(students, s)
	}
	c.JSON(http.StatusOK, students)
}


// GetStudentsSafe осуществляет поиск студентов по имени, используя параметризованный запрос.
// В отличие от небезопасного варианта, этот endpoint не позволяет выполнить SQL-инъекцию и не возвращает данные card_credits.
func (h *StudentHandler) GetStudentsSafe(c *gin.Context) {
	
	queryParam := strings.ToLower(c.Query("query"))
	id, err := strconv.Atoi(queryParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student id", "query": queryParam})
		return
	}

	studentQuery := "SELECT id, age, sex, card_id, name FROM students WHERE id = $1"

	_, err = h.conn.Prepare(context.Background(), "studentSafeStmt", studentQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "query": studentQuery})
		return
	}

	rows, err := h.conn.Query(context.Background(), "studentSafeStmt", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "query": studentQuery})
		return
	}
	defer rows.Close()


	var students []Student
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.Id, &s.Age, &s.Sex, &s.CardId, &s.Name); err != nil {
			c.JSON(500, gin.H{"error": err.Error(), "query": studentQuery})
			return
		}
		students = append(students, s)
	}
	c.JSON(200, students)
}
// RenderSearch рендерит HTML-страницу с двумя вариантами поиска.
// Первый вариант использует уязвимый endpoint (/students), позволяющий через SQL-инъекцию добавить JOIN и вернуть данные card_credits.
// Второй вариант – безопасный поиск (/students_safe), который не возвращает данные card_credits.
func (h *StudentHandler) RenderSearch(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/search.html")
	if err != nil {
		log.Printf("Ошибка парсинга шаблона: %v", err)
		c.String(500, "Ошибка шаблона")
		return
	}

	err = tmpl.Execute(c.Writer, nil)
	if err != nil {
		log.Printf("Ошибка выполнения шаблона: %v", err)
		c.String(500, "Ошибка шаблона")
		return
	}
}