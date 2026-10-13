package main

import (
	"fmt"
	"sync"
	"time"
)

var garfos [5]sync.Mutex
var wg sync.WaitGroup

func filosofo(id int) {
	esquerda := id
	direita := id + 1
	if direita == 5 {
		direita = 0
	}

	primeiro := esquerda
	segundo := direita

	if id == 4 {
		primeiro = direita
		segundo = esquerda
	}

	for ciclo := 1; ciclo <= 20; ciclo++ {
		time.Sleep(100 * time.Millisecond)

		garfos[primeiro].Lock()
		fmt.Println("F", id+1, "pegou o primeiro garfo", primeiro+1, "e quer o segundo", segundo+1)

		garfos[segundo].Lock()

		fmt.Println("F", id+1, "COMENDO, ciclo", ciclo)
		time.Sleep(150 * time.Millisecond)

		garfos[segundo].Unlock()
		garfos[primeiro].Unlock()
	}

	wg.Done()
}

func main() {
	for id := 0; id < 5; id++ {
		wg.Add(1)
		go filosofo(id)
	}
	wg.Wait()
	fmt.Println("Terminou sem deadlock")
}