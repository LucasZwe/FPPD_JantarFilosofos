package main

import (
	"fmt"
	"sync"
	"time"
)

// Configurações
const numFilosofos = 5
const numCiclos = 20
const tempoPensar = 100 * time.Millisecond
const tempoComer = 150 * time.Millisecond

// 0 = V1 normal. 10 * time.Millisecond força deadlock.
const pausaEntreGarfos = 0 * time.Millisecond

// false só medições, true todo processo
const mostrarExecucao = true

var garfos [numFilosofos]sync.Mutex // um mutex para cada garfo
var wg sync.WaitGroup               // espera os filósofos terminarem
var trava sync.Mutex                // protege as métricas abaixo

var ciclosConcluidos [numFilosofos]int
var somaEspera time.Duration
var maiorEspera time.Duration
var totalEsperas int
var comendoAgora int
var maxComendo int

func filosofo(id int) {
	esquerda := id
	direita := id + 1
	if direita == numFilosofos {
		direita = 0
	}

	for ciclo := 1; ciclo <= numCiclos; ciclo++ {

		// pensa
		time.Sleep(tempoPensar)

		// pega garfos
		antes := time.Now()

		garfos[esquerda].Lock()
		if mostrarExecucao {
			fmt.Println("F", id+1, "pegou o garfo", esquerda+1, "e quer o garfo", direita+1)
		}

		if pausaEntreGarfos > 0 {
			time.Sleep(pausaEntreGarfos)
		}

		garfos[direita].Lock()
		espera := time.Since(antes)

		// atualiza métricas
		trava.Lock()
		somaEspera = somaEspera + espera
		totalEsperas = totalEsperas + 1
		if espera > maiorEspera {
			maiorEspera = espera
		}
		comendoAgora = comendoAgora + 1
		if comendoAgora > maxComendo {
			maxComendo = comendoAgora
		}
		trava.Unlock()

		// come
		if mostrarExecucao {
			fmt.Println("F", id+1, "COMENDO, ciclo", ciclo)
		}
		time.Sleep(tempoComer)

		// solta garfos
		trava.Lock()
		comendoAgora = comendoAgora - 1
		ciclosConcluidos[id] = ciclosConcluidos[id] + 1
		trava.Unlock()

		garfos[direita].Unlock()
		garfos[esquerda].Unlock()
	}

	wg.Done()
}

func main() {
	inicio := time.Now()

	for id := 0; id < numFilosofos; id++ {
		wg.Add(1)
		go filosofo(id)
	}
	wg.Wait()

	tempoTotal := time.Since(inicio)

	fmt.Println("=== V1 - Terminou sem deadlock ===")
	fmt.Println("Tempo total:", tempoTotal)
	fmt.Println("Tempo médio de espera:", somaEspera/time.Duration(totalEsperas))
	fmt.Println("Maior tempo de espera:", maiorEspera)
	fmt.Println("Máximo de filósofos comendo ao mesmo tempo:", maxComendo)
	for id := 0; id < numFilosofos; id++ {
		fmt.Println("Filósofo", id+1, "concluiu", ciclosConcluidos[id], "ciclos")
	}
}
