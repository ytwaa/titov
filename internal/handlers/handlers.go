package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"myproject/internal/models"
)

var (
	expenses []models.Expense
	nextID   = 1
	mu       sync.Mutex
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/index.html",
	)
	if err != nil {
		log.Printf("Ошибка парсинга шаблонов: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	mu.Lock()
	list := make([]models.Expense, len(expenses))
	copy(list, expenses)
	mu.Unlock()

	sort.Slice(list, func(i, j int) bool {
		return list[i].Date.After(list[j].Date)
	})

	tmpl.ExecuteTemplate(w, "layout", list)
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		addExpense(w, r)
		return
	}

	tmpl, err := template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/add.html",
	)
	if err != nil {
		log.Printf("Ошибка парсинга шаблонов: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", nil)
}

func addExpense(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	sum, err := strconv.ParseFloat(r.FormValue("sum"), 64)
	if err != nil || sum <= 0 {
		http.Error(w, "Неверный формат или значение суммы", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		http.Error(w, "Неверный формат даты", http.StatusBadRequest)
		return
	}

	mu.Lock()
	expenses = append(expenses, models.Expense{
		ID:          nextID,
		Sum:         sum,
		Description: r.FormValue("description"),
		Date:        date,
	})
	nextID++
	mu.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func About(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Учебный проект «Учёт трат» на Go: веб-интерфейс для записи расходов.")
}

func Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprint(w, "pong")
}
