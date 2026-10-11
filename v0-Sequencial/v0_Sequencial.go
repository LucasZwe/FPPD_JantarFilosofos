package main

import (
	"fmt"
	"time"
)

func main() {
	// configurações
	numFilosofos := 5
	numCiclos := 20
	tempoPensar := 100 * time.Millisecond
	tempoComer := 150 * time.Millisecond
	mostrarExecucao := true

	// guardar valores
	var ciclosConcluidos [5]int
	var somaEspera time.Duration
	var maiorEspera time.Duration
	totalEsperas := 0
	comendoAgora := 0
	maxComendo := 0

	// codigo
	inicio := time.Now()

	if mostrarExecucao {
		fmt.Println("[F1] [F2] [F3] [F4] [F5]")
	}

	for ciclo := 1; ciclo <= numCiclos; ciclo++ {
		for id := 0; id < numFilosofos; id++ {

			// pensa
			time.Sleep(tempoPensar)

			if mostrarExecucao {
				fmt.Printf("%*s T%d\n", id*5, "", ciclo)
			}

			// pega garfos
			antes := time.Now()

			espera := time.Since(antes)

			somaEspera = somaEspera + espera
			totalEsperas = totalEsperas + 1
			if espera > maiorEspera {
				maiorEspera = espera
			}

			// come
			comendoAgora = comendoAgora + 1
			if comendoAgora > maxComendo {
				maxComendo = comendoAgora
			}
			time.Sleep(tempoComer)

			if mostrarExecucao {
				fmt.Printf("%*s E%d\n", id*5, "", ciclo)
			}

			// solta garfos
			comendoAgora = comendoAgora - 1
			ciclosConcluidos[id] = ciclosConcluidos[id] + 1
		}
	}

	tempoTotal := time.Since(inicio)

	// resultados
	fmt.Println("=== V0 - Sequencial ===")
	fmt.Println("Tempo total:", tempoTotal)
	fmt.Println("Tempo médio de espera:", somaEspera/time.Duration(totalEsperas))
	fmt.Println("Maior tempo de espera:", maiorEspera)
	fmt.Println("Máximo de filósofos comendo ao mesmo tempo:", maxComendo)
	for id := 0; id < numFilosofos; id++ {
		fmt.Println("Filósofo", id+1, "concluiu", ciclosConcluidos[id], "ciclos")
	}
}
