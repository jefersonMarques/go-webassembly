package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"

	pb "grpc-client/web" // Ajuste para o caminho correto dos arquivos gerados

	"google.golang.org/grpc"
)

// Defina sua chave API do OpenWeather
const apiKey = "858d50452536881dc0b2ce882156a3f8" // Insira sua chave de API aqui
const apiURL = "http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric"

// Estrutura para resposta da API OpenWeather
type WeatherAPIResponse struct {
	Main struct {
		Temp float32 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

// Implementação do servidor gRPC
type server struct {
	pb.UnimplementedWeatherServiceServer
}

// Função que faz a chamada para a API do OpenWeather
func getWeatherFromAPI(city string) (string, float32, error) {
	// Codifica o nome da cidade corretamente para a URL
	encodedCity := url.QueryEscape(city)

	// Monta a URL da API com a cidade e a chave de API
	url := fmt.Sprintf(apiURL, encodedCity, apiKey)

	// Faz a requisição HTTP para a API do OpenWeather
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, fmt.Errorf("falha na requisição HTTP: %v", err)
	}
	defer resp.Body.Close()

	// Lê a resposta da API
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("falha ao ler a resposta: %v", err)
	}

	// Adicione este log para inspecionar a resposta
	log.Printf("Resposta da API OpenWeather: %s", string(body))

	// Decodifica o JSON da resposta
	var weatherData WeatherAPIResponse
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return "", 0, fmt.Errorf("falha ao decodificar JSON: %v", err)
	}

	// Retorna a descrição e a temperatura
	description := weatherData.Weather[0].Description
	temperature := weatherData.Main.Temp

	return description, temperature, nil
}

// Implementação do método GetWeather do servidor gRPC
func (s *server) GetWeather(ctx context.Context, req *pb.WeatherRequest) (*pb.WeatherResponse, error) {
	log.Printf("Recebendo requisição para cidade: %s", req.City)

	// Obtém os dados reais da API
	description, temperature, err := getWeatherFromAPI(req.City)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter dados do clima: %v", err)
	}

	// Retorna a resposta gRPC com os dados reais
	return &pb.WeatherResponse{
		City:        req.City,
		Description: description,
		Temperature: temperature,
	}, nil
}

func main() {
	// Cria um listener na porta 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Falha ao escutar: %v", err)
	}

	// Cria uma instância do servidor gRPC
	s := grpc.NewServer()
	pb.RegisterWeatherServiceServer(s, &server{})

	log.Printf("Servidor gRPC rodando na porta 50051")

	// Inicia o servidor gRPC
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Falha ao servir: %v", err)
	}
}
