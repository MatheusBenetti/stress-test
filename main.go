package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	// Definir flags de linha de comando
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 0, "Número total de requests")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")
	flag.Parse()

	if *url == "" || *requests == 0 || *concurrency == 0 {
		fmt.Println("Por favor, forneça todos os parâmetros necessários.")
		return
	}

	// Iniciar o teste de carga
	startTime := time.Now()
	var wg sync.WaitGroup
	counter := make(chan int, *requests)

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			makeRequests(*url, *requests / *concurrency, counter)
		}()
	}

	go func() {
		wg.Wait()
		close(counter)
	}()

	// Agregar resultados
	totalRequests := 0
	status200 := 0
	otherStatus := make(map[int]int)

	for count := range counter {
		totalRequests++
		if count == 200 {
			status200++
		} else {
			otherStatus[count]++
		}
	}

	// Gerar relatório
	duration := time.Since(startTime)
	fmt.Printf("Tempo total gasto na execução: %v\n", duration)
	fmt.Printf("Quantidade total de requests realizados: %d\n", totalRequests)
	fmt.Printf("Quantidade de requests com status HTTP 200: %d\n", status200)

	for code, count := range otherStatus {
		fmt.Printf("Quantidade de requests com status HTTP %d: %d\n", code, count)
	}
}

func makeRequests(url string, numRequests int, counter chan<- int) {
	for i := 0; i < numRequests; i++ {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Erro ao realizar request:", err)
			counter <- 500
		} else {
			counter <- resp.StatusCode
			resp.Body.Close()
		}
	}
}
