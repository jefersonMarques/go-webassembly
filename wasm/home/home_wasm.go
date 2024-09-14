package main

import (
	"syscall/js"
)

// Função que será chamada para renderizar a página Weather
// Essa função cria o conteúdo HTML da página de clima, incluindo um formulário
// onde o usuário pode inserir o nome de uma cidade e submeter para obter a previsão do tempo.
// Após gerar o HTML, ele adiciona o evento de submissão para o formulário de busca de clima.
func renderWeatherPage() {
	// Define o HTML da página "Weather", incluindo um formulário para inserir o nome da cidade.
	content := `<h1>Página inicial</h1>`

	// Atualiza o conteúdo da div com id "content", substituindo o conteúdo atual pelo novo HTML da página "Weather".
	document := js.Global().Get("document")
	contentDiv := document.Call("getElementById", "content")
	contentDiv.Set("innerHTML", content)
}

func main() {
	// Chama a função para renderizar a página Weather quando o módulo WASM for carregado.
	renderWeatherPage()

	// Mantém o WebAssembly rodando indefinidamente. Sem isso, o módulo seria finalizado imediatamente após carregar a página.
	select {}
}
