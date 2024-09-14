declare var Go: any;

async function initWasm() {
    debugger
    const go = new Go();  // Definido no wasm_exec.js
    try {
        const wasmModule = await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject);
        console.log("WebAssembly module loaded");

        go.run(wasmModule.instance);

        if ((window as any).getWeather === undefined) {
            console.error("getWeather is not defined in WebAssembly");
        } else {
            console.log("getWeather function is defined.");
        }
    } catch (err) {
        console.error("Error loading WebAssembly:", err);
    }
}

// Inicia o WebAssembly
initWasm();

const greetBtn = document.getElementById("greet-btn");
if (greetBtn) {
    console.log("Button found");
    greetBtn.addEventListener("click", () => {
        const cityInput = (document.getElementById("city-input") as HTMLInputElement).value;
        console.log("City input:", cityInput);
        
        if (cityInput) {
            console.log("Calling getWeather with city:", cityInput);
            (window as any).getWeather(cityInput);  // Chama a função getWeather do WebAssembly com o valor da cidade
        } else {
            console.error("City input is empty.");
        }
    });
} else {
    console.error("Button not found");
}



