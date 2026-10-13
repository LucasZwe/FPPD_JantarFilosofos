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

// false só medicoes, true todo processo
const mostrarExecucao = true

var garfos [numFilosofos]sync.Mutex
var wg sync.WaitGroup
var trava sync.Mutex

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

	// todos pegam primeiro MENOS numero
	// o ultimo inverte a ordem
	primeiro := esquerda
	segundo := direita

	if id == numFilosofos-1 {
		primeiro = direita
		segundo = esquerda
	}

	for ciclo := 1; ciclo <= numCiclos; ciclo++ {
		//pensa
		time.Sleep(tempoPensar)

		antes := time.Now()

		garfos[primeiro].Lock()
		if mostrarExecucao {
			fmt.Println("F", id+1, "pegou o primeiro garfo", primeiro+1, "e quer o segundo", segundo+1)
		}

		garfos[segundo].Lock()
		espera := time.Since(antes)

		//atualiza metrica
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

		//come
		if mostrarExecucao {
			fmt.Println("F", id+1, "COMENDO, ciclo", ciclo)
		}
		time.Sleep(tempoComer)

		trava.Lock()
		comendoAgora = comendoAgora - 1
		ciclosConcluidos[id] = ciclosConcluidos[id] + 1
		trava.Unlock()

		// solta garfos
		garfos[segundo].Unlock()
		garfos[primeiro].Unlock()
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

	fmt.Println("=== V3 - Ordenação dos garfos ===")
	fmt.Println("Tempo total:", tempoTotal)
	fmt.Println("Tempo médio de espera:", somaEspera/time.Duration(totalEsperas))
	fmt.Println("Maior tempo de espera:", maiorEspera)
	fmt.Println("Máximo de filósofos comendo ao mesmo tempo:", maxComendo)
	for id := 0; id < numFilosofos; id++ {
		fmt.Println("Filósofo", id+1, "concluiu", ciclosConcluidos[id], "ciclos")
	}
}
