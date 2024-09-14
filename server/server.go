package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	pb "grpc-client/web" // Ajuste o caminho para o pacote gerado

	"google.golang.org/grpc"
)

// Estrutura para armazenar a resposta do clima que será enviada ao cliente
type WeatherResponse struct {
	City        string  `json:"city"`
	Description string  `json:"description"`
	Temperature float32 `json:"temperature"`
}

// Função para buscar os dados de clima via gRPC
func getWeatherData(city string) (*WeatherResponse, error) {
	// Conecta ao servidor gRPC
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao servidor gRPC: %v", err)
	}
	defer conn.Close()

	client := pb.NewWeatherServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Faz a requisição gRPC para obter os dados de clima
	res, err := client.GetWeather(ctx, &pb.WeatherRequest{City: city})
	if err != nil {
		return nil, fmt.Errorf("erro ao obter dados de clima: %v", err)
	}

	// Prepara a resposta com os dados de clima
	return &WeatherResponse{
		City:        res.City,
		Description: res.Description,
		Temperature: res.Temperature,
	}, nil
}

// Função para lidar com a rota /weather e buscar o clima via gRPC
func handleWeather(w http.ResponseWriter, r *http.Request) {
	// Obtém a cidade da query string (ex: ?city=SaoPaulo)
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "Cidade não especificada", http.StatusBadRequest)
		return
	}

	// Faz a chamada ao gRPC para buscar os dados do clima
	weatherData, err := getWeatherData(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Define o cabeçalho da resposta como JSON
	w.Header().Set("Content-Type", "application/json")
	// Envia a resposta JSON para o frontend
	json.NewEncoder(w).Encode(weatherData)
}

// Função para servir o arquivo index.html
func serveIndex(w http.ResponseWriter, r *http.Request) {
	// Serve o arquivo index.html da pasta frontend
	http.ServeFile(w, r, "frontend/index.html")
}

func main() {
	// Servir arquivos estáticos
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend"))))

	// Rota para servir o index.html
	http.HandleFunc("/", serveIndex)

	// Rota para buscar o clima via HTTP e gRPC
	http.HandleFunc("/weather", handleWeather)

	log.Println("Servidor rodando em http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
