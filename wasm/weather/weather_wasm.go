package main

import (
	"fmt"
	"syscall/js"
)

// Função que será chamada para renderizar a página Weather
// Essa função cria o conteúdo HTML da página de clima, incluindo um formulário
// onde o usuário pode inserir o nome de uma cidade e submeter para obter a previsão do tempo.
// Após gerar o HTML, ele adiciona o evento de submissão para o formulário de busca de clima.
func renderWeatherPage() {
	// Define o HTML da página "Weather", incluindo um formulário para inserir o nome da cidade.
	content := `<h1>Weather Page</h1><p>Insira uma cidade para buscar o clima:</p>
	<form id="weatherForm" onsubmit="event.preventDefault(); if (typeof getWeather === 'function') getWeather(event);">
		<input type="text" id="cityInput" placeholder="Nome da cidade"/>
		<button type="submit">Buscar Clima</button>
	</form>
	<div id="output"></div>`

	// Atualiza o conteúdo da div com id "content", substituindo o conteúdo atual pelo novo HTML da página "Weather".
	document := js.Global().Get("document")
	contentDiv := document.Call("getElementById", "content")
	contentDiv.Set("innerHTML", content)

	// Adiciona um listener de evento para o formulário de busca de clima.
	// O evento "submit" dispara a função getWeather quando o usuário submeter o formulário.
	form := document.Call("getElementById", "weatherForm")
	form.Call("addEventListener", "submit", js.FuncOf(getWeather))
}

// Função para buscar o clima no backend
// Essa função é disparada quando o formulário é submetido.
// Ela evita o comportamento padrão de reload da página e obtém o valor inserido pelo usuário (nome da cidade).
// Em seguida, chama a função fetchWeather para fazer a requisição ao backend com o nome da cidade.
func getWeather(this js.Value, p []js.Value) interface{} {
	event := p[0]
	event.Call("preventDefault") // Evita o comportamento padrão de reload da página ao submeter o formulário.

	// Obtém o valor digitado no campo de input (nome da cidade).
	document := js.Global().Get("document")
	city := document.Call("getElementById", "cityInput").Get("value").String()

	// Inicia uma requisição ao backend de forma assíncrona para buscar o clima da cidade inserida.
	go fetchWeather(city)

	return nil
}

// Função para realizar a requisição HTTP ao backend
// Essa função constrói a URL para o backend usando o nome da cidade inserida pelo usuário.
// Faz uma requisição HTTP utilizando "fetch" e processa a resposta com promises para obter os dados do clima.
// A função atualiza a interface com os dados da cidade, temperatura e descrição do clima ou exibe uma mensagem de erro caso a requisição falhe.
func fetchWeather(city string) {
	// Constroi a URL da API do backend para buscar o clima da cidade inserida.
	url := "/weather?city=" + city

	// Realiza a requisição HTTP ao backend usando fetch.
	fetch := js.Global().Call("fetch", url)

	// Processa a resposta do servidor com promises. Primeiro, obtemos a resposta no formato JSON.
	fetch.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resp := args[0]
		return resp.Call("json")
	})).Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Decodifica o JSON e extrai os dados do clima (nome da cidade, descrição e temperatura).
		json := args[0]
		cityName := json.Get("city").String()           // Nome da cidade
		description := json.Get("description").String() // Descrição do clima (ex.: "nublado")
		temperature := json.Get("temperature").Float()  // Temperatura em graus Celsius

		// Formata as informações do clima para exibição na interface.
		weatherInfo := "Cidade: " + cityName + "\nTemperatura: " + fmt.Sprintf("%.2f", temperature) + "°C\nDescrição: " + description

		// Atualiza a interface exibindo as informações de clima.
		updateOutput(weatherInfo)
		return nil
	})).Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Caso haja um erro na requisição, exibe uma mensagem de erro.
		updateOutput("Erro ao obter dados de clima")
		return nil
	}))
}

// Função para atualizar a saída do clima
// Essa função recebe os dados de clima como string e exibe-os no elemento com id "output".
// Se ocorrer um erro, a função exibe a mensagem apropriada no mesmo elemento.
func updateOutput(data string) {
	document := js.Global().Get("document")
	output := document.Call("getElementById", "output")
	output.Set("innerText", data) // Define o texto exibido no elemento de saída (div com id "output").
}

func main() {
	// Chama a função para renderizar a página Weather quando o módulo WASM for carregado.
	renderWeatherPage()

	// Mantém o WebAssembly rodando indefinidamente. Sem isso, o módulo seria finalizado imediatamente após carregar a página.
	select {}
}
