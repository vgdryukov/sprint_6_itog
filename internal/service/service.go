package service

import (
	"errors"
	"log"
	"sprint_6_itog/pkg/morse"
	"strings"
	"unicode"
)

type Service struct {
	logger *log.Logger
	config Config
}

type Config struct {
	StrictValidation bool
	MinInputLength   int
	MaxInputLength   int
}

func NewService(logger *log.Logger) *Service {
	return &Service{
		logger: logger,
		config: Config{
			StrictValidation: true,
			MinInputLength:   1,
			MaxInputLength:   10000,
		},
	}
}

// DefiningContentType определяет тип входящей строки - кириллица или азбука Морзе
func (s *Service) DefiningContentType(input string) (string, error) {
	if input == "" {
		return "", errors.New("input string is empty")
	}

	// Убираем лишние пробелы по краям
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", errors.New("input contains only whitespace characters")
	}

	// Анализируем строку для определения типа
	isMorse := s.isMorseCode(trimmed)
	isCyrillic := s.isCyrillicText(trimmed)

	switch {
	case isMorse:
		return "Morse", nil
	case isCyrillic:
		return "Cyrillic", nil
	default:
		return "", errors.New("the text does not correspond to either Cyrillic or Morse code")
	}
}

// isMorseCode проверяет, является ли строка кодом Морзе
func (s *Service) isMorseCode(input string) bool {
	morseCharCount := 0
	totalChars := 0

	for _, char := range input {
		totalChars++
		switch char {
		case '.', '-', ' ', '\t', '\n':
			morseCharCount++
		default:
			return false
		}
	}

	// Если слишком мало морзе-символов, вероятно это не код Морзе
	if totalChars > 0 && float64(morseCharCount)/float64(totalChars) < 0.8 {
		return false
	}

	// Разбиваем на слова и проверяем каждое слово отдельно
	words := strings.Fields(input)
	if len(words) == 0 {
		return false
	}

	for _, word := range words {
		if !s.isValidMorseWord(word) {
			return false
		}
	}

	return true
}

// isValidMorseWord проверяет валидность отдельного слова в морзе-коде
func (s *Service) isValidMorseWord(word string) bool {
	if len(word) == 0 || len(word) > 6 {
		return false
	}

	// Проверяем, что слово состоит только из точек и тире
	for _, char := range word {
		if char != '.' && char != '-' {
			return false
		}
	}

	// Проверяем на наличие хотя бы одного тире или точки
	hasDot := strings.Contains(word, ".")
	hasDash := strings.Contains(word, "-")

	return hasDot || hasDash
}

// isCyrillicText проверяет, является ли строка текстом на кириллице
func (s *Service) isCyrillicText(input string) bool {
	cyrillicCount := 0
	letterCount := 0

	for _, char := range input {
		// Игнорируем пробелы и пунктуацию при подсчете
		if unicode.IsSpace(char) || unicode.IsPunct(char) {
			continue
		}

		if unicode.IsLetter(char) {
			letterCount++
			if unicode.Is(unicode.Cyrillic, char) {
				cyrillicCount++
			} else if unicode.IsLetter(char) {
				// Если есть латинские буквы - это не чистая кириллица
				return false
			}
		}
	}

	// Если нет букв вообще
	if letterCount == 0 {
		return false
	}

	// Если более 80% букв - кириллица, считаем что это кириллический текст
	return float64(cyrillicCount)/float64(letterCount) >= 0.8
}

func (s *Service) AutoDetectAndConvert(input string) (string, error) {
	// Логирование
	if s.logger != nil {
		s.logger.Printf("Processing input, length: %d", len(input))
	}

	// Валидация длины
	if len(input) < s.config.MinInputLength {
		return "", errors.New("input too short")
	}
	if len(input) > s.config.MaxInputLength {
		return "", errors.New("input too long")
	}

	// Определение типа и конвертация
	contentType, err := s.DefiningContentType(input)
	if err != nil {
		return "", err
	}

	switch contentType {
	case "Cyrillic":
		return morse.ToMorse(input), nil
	case "Morse":
		return morse.ToText(input), nil
	default:
		return "", errors.New("unsupported content type")
	}
}
