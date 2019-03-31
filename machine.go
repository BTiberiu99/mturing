package main

import (
	"bufio"
	"fmt"
	"mturing/machine"
	"os"
	"os/exec"
	"strings"
)

func main() {
	test()
}

func startMachine() {
	machine := &machine.Machine{}
	machine.ConstructMachine("../machine.txt")
	//reader := bufio.NewReader(os.Stdin)
	for {

	}
}

type Command struct {
	Name        string
	Description string
}

func test() {
	reader := bufio.NewReader(os.Stdin)
	machine := &machine.Machine{}
	fmt.Println("Machine started")
	commandsMachine := []Command{Command{Name: "construct", Description: "Construieste o machine din fisierul dat ca al 2 lea argument"}, Command{Name: "parse", Description: "Parseaza un cuvant cu masina construita primit ca al 2 lea argument"}, Command{Name: "print", Description: "Printeaza masina"}, Command{Name: "close", Description: "Inchide programul"}, Command{Name: "clear", Description: "Goleste consola"}, Command{Name: "help", Description: "Afiseaza informatii despre comenzzi"}, Command{Name: "debug", Description: "Afiseaza informatii din iteratiile de parsare"}}
	// , "parse": "Parseaza un cuvant cu masina construita primit ca al 2 lea argument", "print": "Printeaza masina", "close": "Inchide programul", "clear": "Goleste consola", "help": "Afiseaza informatii despre comenzzi"
	infinite := true
	for infinite {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		commands := strings.Split(text, " ")
		if len(commands) < 1 {
			fmt.Println("Introduceti o comanda")
			continue
		}
		switch commands[0] {
		case commandsMachine[0].Name:
			if len(commands) < 2 {
				fmt.Println("Introduceti numele fisierului")
				continue
			}
			err := machine.ConstructMachine(commands[1])
			if err != nil {
				fmt.Println(err.Error())
			}
			break
		case commandsMachine[1].Name:
			if len(commands) < 2 {
				fmt.Println("Introduceti cuvantul")
				continue
			}
			accepted := machine.ParseWord(commands[1])
			if accepted {
				fmt.Println(commands[1], "este accepetat")
			} else {
				fmt.Println(commands[1], "nu este accepetat")
			}
			break
		case commandsMachine[2].Name:
			machine.PrintMachine()
			break
		case commandsMachine[3].Name:
			infinite = false
			fmt.Println("Program inchis")
			break
		case commandsMachine[4].Name:
			c := exec.Command("clear")
			c.Stdout = os.Stdout
			c.Run()
			break
		case commandsMachine[5].Name:
			fmt.Println("Comenzi:")
			for _, item := range commandsMachine {
				fmt.Println(item.Name + " ---> " + item.Description)
			}
			break
		case commandsMachine[6].Name:
			if len(commands) > 1 {
				machine.Print = commands[1] == "true"
			} else {
				fmt.Println("Comanda print are nevoie de al 2-lea argument true/false")
			}
		default:
			fmt.Println("Commanda nerecunoscuta")
		}

	}

}
