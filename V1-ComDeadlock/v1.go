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

	for ciclo := 1; ciclo <= 20; ciclo++ {
		// pensa
		time.Sleep(100 * time.Millisecond)

		// pega os garfos
		garfos[esquerda].Lock()
		fmt.Println("F", id+1, "pegou o garfo", esquerda+1, "e quer o garfo", direita+1)

		garfos[direita].Lock()

		// come
		fmt.Println("F", id+1, "COMENDO, ciclo", ciclo)
		time.Sleep(150 * time.Millisecond)

		// solta os garfos
		garfos[direita].Unlock()
		garfos[esquerda].Unlock()
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
