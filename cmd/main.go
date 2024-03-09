package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 100, "Número total de requests")
	concurrency := flag.Int("concurrency", 10, "Número de chamadas simultâneas")
	flag.Parse()

	// Convertendo ponteiros para valores
	urlValue := *url
	requestsValue := *requests
	concurrencyValue := *concurrency

	fmt.Printf("URL: %s\n", urlValue)
	fmt.Printf("Número total de requests: %d\n", requestsValue)
	fmt.Printf("Número de chamadas simultâneas: %d\n", concurrencyValue)

	// Agora você pode usar urlValue, requestsValue e concurrencyValue diretamente
	if urlValue == "" {
		fmt.Println("URL é obrigatório")
		return
	}

	fmt.Println("Iniciando testes...")

	start := time.Now()

	jobs := make(chan struct{}, requestsValue)
	results := make(chan int, requestsValue)

	var wg sync.WaitGroup
	wg.Add(concurrencyValue)

	for i := 0; i < concurrencyValue; i++ {
		go worker(urlValue, jobs, results, &wg)
	}

	for i := 0; i < requestsValue; i++ {
		jobs <- struct{}{}
	}
	close(jobs)

	wg.Wait()
	close(results)

	totalTime := time.Since(start)
	statusCodes := make(map[int]int)
	for result := range results {
		statusCodes[result]++
	}

	fmt.Printf("Teste concluído em %v\n", totalTime)
	fmt.Printf("Requests realizados: %d\n", requestsValue)
	for status, count := range statusCodes {
		fmt.Printf("Status %d: %d requests\n", status, count)
	}
}

func worker(url string, jobs <-chan struct{}, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for range jobs {
		statusCode := makeRequest(url)
		results <- statusCode
	}
}

func makeRequest(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Erro ao realizar request: %v\n", err)
		return 0
	}
	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		fmt.Printf("Erro ao ler o corpo da resposta: %v\n", readErr)
		return 0 // Considerando 0 como um código de erro interno
	}

	fmt.Println(fmt.Sprintf("HTTP %d > %s | %s", resp.StatusCode, url, fmt.Sprintf(string(body))))

	// Retorna o status code da resposta HTTP

	return resp.StatusCode
}
