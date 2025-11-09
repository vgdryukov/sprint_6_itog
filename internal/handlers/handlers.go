package handlers

import (
	"io"
	"net/http"
	"os"
	"sprint_6_itog/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RootHandler обрабатывает корневой эндпоинт /
func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Читаем и отдаем index.html
	html, err := os.ReadFile("index.html")
	if err != nil {
		http.Error(w, "Could not read index.html", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(html)
}

// UploadHandler обрабатывает эндпоинт /upload
func (h *Handler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсим форму
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusInternalServerError)
		return
	}

	// Получаем файл из формы
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file from form", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Проверяем, что файл не пустой
	if header.Size == 0 {
		http.Error(w, "File is empty", http.StatusBadRequest)
		return
	}

	// Читаем данные из файла
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file data", http.StatusInternalServerError)
		return
	}

	// Передаем данные в функцию автоопределения
	convertedString, err := h.service.AutoDetectAndConvert(string(fileData))
	if err != nil {
		http.Error(w, "Conversion error", http.StatusInternalServerError)
		return
	}

	// Создаем локальный файл для результата
	/*timestamp := time.Now().UTC().Format("2006-01-02_15-04-05")
	originalExt := filepath.Ext(header.Filename)
	outputFilename := fmt.Sprintf("result_%s%s", timestamp, originalExt)

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		http.Error(w, "Unable to create output file", http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	// Записываем результат конвертации в файл
	_, err = outputFile.WriteString(convertedString)
	if err != nil {
		http.Error(w, "Unable to write to output file", http.StatusInternalServerError)
		return
	}*/

	// Возвращаем результат конвертации
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(convertedString))
}
