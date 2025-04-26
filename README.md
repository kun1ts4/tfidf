# Анализатор TF-IDF (тестовое Леста)

## Функционал 

* Загрузка текстового файла
* Извлечение слов из текста
* Расчет TF-IDF
* Отображение топ-50 слов

## Установка
```bash
git clone https://github.com/kun1ts4/tfidf.git
cd tfidf
go mod tidy
```

## Запуск
```bash
go run cmd/api/main.go
```
Откройте http://localhost:8080

> [!IMPORTANT]  
> Каждый абзац в загружаемом файле считается отдельным документом, а файл — набором документов

## Структура проекта

* cmd/api/main.go — точка входа в приложение
* internal/handler — обработчики HTTP-запросов
* internal/parser — функции для обработки текста
* internal/service — логика расчета TF-IDF
* internal/model — структуры данных
* web/templates — HTML-шаблоны для интерфейса
* web/templates/styles.css — стили для веб-интерфейса