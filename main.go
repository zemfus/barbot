package main

import (
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, HTTPS world!"))
}

func main() {
	// Функция для запуска HTTPS сервера с autocert
	startHTTPS := func(domain string) {
		m := &autocert.Manager{
			Cache:      autocert.DirCache("cert-cache"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domain), // Замените "yourdomain.com" вашим доменом
		}

		server := &http.Server{
			Addr:      ":https",
			TLSConfig: m.TLSConfig(),
			Handler:   http.DefaultServeMux,
		}

		// Ручка для hello
		http.HandleFunc("/hello", helloHandler)

		// Запуск HTTP сервера для редиректа на HTTPS и ACME-вызова
		go func() {
			log.Fatal(http.ListenAndServe(":http", m.HTTPHandler(nil)))
		}()

		log.Println("Starting HTTPS server on :443...")
		err := server.ListenAndServeTLS("", "") // Пути к сертификатам пусты, так как autocert их автоматически управляет
		if err != nil {
			log.Fatalf("Failed to start HTTPS server: %v", err)
		}
	}

	// Запуск HTTP сервера для локальной разработки или тестирования без SSL
	//startHTTP := func() {
	//	// Ручка для hello
	//	http.HandleFunc("/hello", helloHandler)
	//
	//	// Запуск HTTP сервера
	//	log.Println("Starting HTTP server on :8080...")
	//	err := http.ListenAndServe(":8080", nil)
	//	if err != nil {
	//		log.Fatalf("Failed to start HTTP server: %v", err)
	//	}
	//}

	// Здесь вы можете определить, какой сервер запустить, в зависимости от вашего окружения
	// Например, использовать флаг командной строки, переменную окружения или любую другую логику

	// Предположим, что для примера мы запускаем HTTPS сервер с autocert для домена "yourdomain.com"
	// В продакшене замените "yourdomain.com" на ваш реальный домен и расскомментируйте следующую строку
	startHTTPS("zemfus.com")

	// Для локальной разработки запускаем HTTP сервер
	//startHTTP()
}
