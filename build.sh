#!/bin/bash

# Compilar main_wasm
echo "Compilando main_wasm.go..."
GOOS=js GOARCH=wasm go build -o frontend/main.wasm wasm/main/main_wasm.go

# Compilar weather_wasm
echo "Compilando weather_wasm.go..."
GOOS=js GOARCH=wasm go build -o frontend/weather.wasm wasm/weather/weather_wasm.go

# Compilar home_wasm
echo "Compilando home..."
GOOS=js GOARCH=wasm go build -o frontend/home.wasm wasm/home/home_wasm.go

echo "Compilação completa!"
