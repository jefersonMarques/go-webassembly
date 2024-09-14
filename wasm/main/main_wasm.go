package main

import (
	"syscall/js"
)

// Função para simular navegação sem alterar a URL visível
// Esta função é responsável por mudar a página da aplicação sem modificar a URL no navegador.
// O estado da página é armazenado no histórico do navegador, mas a URL visível não é alterada.
func navigate(page string) {
	// Adiciona o estado da página ao histórico sem alterar a URL visível
	js.Global().Get("history").Call("pushState", js.ValueOf(page), "", "")

	// Altera o conteúdo da página com base na página solicitada
	changeContent(page)
}

// Função para alterar o conteúdo da página com base no parâmetro "page".
// Esta função carrega o conteúdo apropriado (HTML) de acordo com a página solicitada.
// Se a página for "weather", o módulo WebAssembly do clima é carregado dinamicamente.
func changeContent(page string) {
	switch page {
	case "home":
		// Conteúdo da página "Home"
		loadHome()
	case "weather":
		// Carrega o módulo WASM do Weather sob demanda
		loadWeatherModule()
	case "about":
		// Conteúdo da página "About"
		content := `<h1>SObre</h1><p>Este é o conteúdo da página sobre nós.</p>`
		setContent(content)
	default:
		// Se a página não for encontrada, exibe a página de erro 404
		content := `<h1>404 - Página não encontrada</h1>`
		setContent(content)
	}
}

func loadHome() {
	homeWasmPath := "/static/home.wasm"       // Caminho do arquivo WASM home
	goInstance := js.Global().Get("Go").New() // Cria uma nova instância do Go para rodar o WASM

	// Faz a requisição para obter o arquivo "home.wasm"
	fetchWasm := js.Global().Call("fetch", homeWasmPath)

	fetchWasm.Call("then", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		resp := p[0]

		if js.Global().Get("WebAssembly").Get("instantiateStreaming").Truthy() {
			// Usa instantiateStreaming se disponível para carregar e instanciar o WASM a partir do stream
			js.Global().Get("WebAssembly").Call("instantiateStreaming", resp, goInstance.Get("importObject")).Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				result := args[0]
				goInstance.Call("run", result.Get("instance")) // Executa o módulo WASM da home
				return nil
			}))
		} else {
			// Fallback: carrega o arquivo WASM como ArrayBuffer e usa WebAssembly.instantiate
			resp.Call("arrayBuffer").Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				buffer := args[0]
				js.Global().Get("WebAssembly").Call("instantiate", buffer, goInstance.Get("importObject")).Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					result := args[0]
					goInstance.Call("run", result.Get("instance")) // Executa o módulo WASM da home
					return nil
				}))
				return nil
			}))
		}
		return nil
	})).Call("catch", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		// Caso ocorra um erro ao carregar o arquivo WASM, exibe uma mensagem no console
		js.Global().Get("console").Call("error", "Falha ao carregar weather.wasm", p[0])
		return nil
	}))
}

// Função para atualizar o conteúdo da div #content
// Esta função recebe o HTML como string e o insere no elemento "content" do DOM.
// Usada para atualizar o conteúdo exibido na página com base na navegação.
func setContent(content string) {
	document := js.Global().Get("document")
	contentDiv := document.Call("getElementById", "content")
	contentDiv.Set("innerHTML", content)
}

// Função para carregar o módulo WASM do Weather sob demanda
// Esta função carrega o arquivo WebAssembly "weather.wasm" apenas quando a página de clima é solicitada.
// Usa WebAssembly.instantiateStreaming se disponível, com fallback para WebAssembly.instantiate.
func loadWeatherModule() {
	weatherWasmPath := "/static/weather.wasm" // Caminho do arquivo WASM do clima
	goInstance := js.Global().Get("Go").New() // Cria uma nova instância do Go para rodar o WASM

	// Faz a requisição para obter o arquivo "weather.wasm"
	fetchWasm := js.Global().Call("fetch", weatherWasmPath)

	// Verifica se o navegador suporta WebAssembly.instantiateStreaming
	fetchWasm.Call("then", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		resp := p[0]

		if js.Global().Get("WebAssembly").Get("instantiateStreaming").Truthy() {
			// Usa instantiateStreaming se disponível para carregar e instanciar o WASM a partir do stream
			js.Global().Get("WebAssembly").Call("instantiateStreaming", resp, goInstance.Get("importObject")).Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				result := args[0]
				goInstance.Call("run", result.Get("instance")) // Executa o módulo WASM do clima
				return nil
			}))
		} else {
			// Fallback: carrega o arquivo WASM como ArrayBuffer e usa WebAssembly.instantiate
			resp.Call("arrayBuffer").Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				buffer := args[0]
				js.Global().Get("WebAssembly").Call("instantiate", buffer, goInstance.Get("importObject")).Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					result := args[0]
					goInstance.Call("run", result.Get("instance")) // Executa o módulo WASM do clima
					return nil
				}))
				return nil
			}))
		}
		return nil
	})).Call("catch", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		// Caso ocorra um erro ao carregar o arquivo WASM, exibe uma mensagem no console
		js.Global().Get("console").Call("error", "Falha ao carregar weather.wasm", p[0])
		return nil
	}))
}

// Função que escuta mudanças no histórico (evento popstate)
// Esta função é chamada quando o usuário usa o botão "voltar" ou "avançar" do navegador.
// Com base no estado armazenado no histórico, ela recarrega a página correta.
func listenToPopState(this js.Value, p []js.Value) interface{} {
	state := js.Global().Get("history").Get("state").String() // Obtém o estado da página

	// Caso o estado seja vazio (página inicial), define como 'home'
	if state == "" {
		state = "home"
	}

	// Carrega o conteúdo correto com base no estado do histórico
	changeContent(state)
	return nil
}

// Função para delegação de eventos no nav (barra de navegação)
// Esta função adiciona um único EventListener ao elemento <nav> e usa delegação de eventos para gerenciar os cliques nos links.
// Isso economiza memória ao registrar apenas um ouvinte de evento para todos os links de navegação.
func addEventDelegation() {
	document := js.Global().Get("document")

	// Obtém o elemento <nav> e adiciona um EventListener para capturar cliques
	nav := document.Call("getElementById", "main-nav")
	nav.Call("addEventListener", "click", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		event := p[0]
		event.Call("preventDefault") // Evita o comportamento padrão de redirecionamento de link

		// Verifica qual link foi clicado (Home, Weather ou About) e navega para a página correta
		targetID := event.Get("target").Get("id").String()
		switch targetID {
		case "nav-home":
			navigate("home")
		case "nav-weather":
			navigate("weather")
		case "nav-about":
			navigate("about")
		}
		return nil
	}))
}

func main() {
	// Inicializa a delegação de eventos na barra de navegação
	addEventDelegation()

	// Escuta mudanças no histórico do navegador
	js.Global().Get("window").Call("addEventListener", "popstate", js.FuncOf(listenToPopState))

	// Mantém o WebAssembly rodando indefinidamente
	select {}
}
