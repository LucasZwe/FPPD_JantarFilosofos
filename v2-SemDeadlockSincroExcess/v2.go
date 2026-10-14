package main

import (
	"fmt"
	"sync"
	"time"
)

const numFilosofos = 5
const numCiclos = 20
const tempoPensar = 100 * time.Millisecond
const tempoComer = 150 * time.Millisecond

// false só medições, true todo processo
const mostrarExecucao = true

var garfos [numFilosofos]sync.Mutex // um mutex para cada garfo
var mesa sync.Mutex                 // "garçom" - SÓ 1 filósofo por vez pode pegar garfos e comer
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

		mesa.Lock()
		if mostrarExecucao {
			fmt.Println("F", id+1, "entrou na mesa e quer os garfos", esquerda+1, "e", direita+1)
		}

		garfos[esquerda].Lock()
		garfos[direita].Lock()
		espera := time.Since(antes)

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

		if mostrarExecucao {
			fmt.Println("F", id+1, "COMENDO, ciclo", ciclo)
		}
		time.Sleep(tempoComer)

		trava.Lock()
		comendoAgora = comendoAgora - 1
		ciclosConcluidos[id] = ciclosConcluidos[id] + 1
		trava.Unlock()

		garfos[direita].Unlock()
		garfos[esquerda].Unlock()

		// libera a mesa para o próximo
		mesa.Unlock()
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

	fmt.Println("=== V2 - Sincronização excessiva (mutex global) ===")
	fmt.Println("Tempo total:", tempoTotal)
	fmt.Println("Tempo médio de espera:", somaEspera/time.Duration(totalEsperas))
	fmt.Println("Maior tempo de espera:", maiorEspera)
	fmt.Println("Máximo de filósofos comendo ao mesmo tempo:", maxComendo)
	for id := 0; id < numFilosofos; id++ {
		fmt.Println("Filósofo", id+1, "concluiu", ciclosConcluidos[id], "ciclos")
	}
}