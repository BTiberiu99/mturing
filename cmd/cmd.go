package cmd

import (
	"bufio"
	"fmt"
	"machine/machine"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Command struct {
	Name        string
	Description string
}

func (cmd *Command) ToString() string {
	return fmt.Sprintf("%s ---> %s", cmd.Name, cmd.Description)
}

func Run() {

	fmt.Println("Machine started")
	var (
		mch *machine.Machine
	)

	reader := bufio.NewReader(os.Stdin)

	commandsMachine := []Command{
		Command{Name: "construct", Description: "Construieste o machine din fisierul dat ca al 2 lea argument"},
		Command{Name: "parse", Description: "Parseaza un cuvant cu masina construita primit ca al 2 lea argument"},
		Command{Name: "print", Description: "Printeaza masina"},
		Command{Name: "close", Description: "Inchide programul"},
		Command{Name: "clear", Description: "Goleste consola"},
		Command{Name: "help", Description: "Afiseaza informatii despre comenzzi"},
		Command{Name: "debug", Description: "Afiseaza informatii din iteratiile de parsare"},
	}
	run := true
	for run {
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
			m, err := machine.NewMachine(&machine.Config{
				NameFile: commands[1],
			})

			if err != nil {
				fmt.Println(err)
			} else {
				mch = m
			}
			break
		case commandsMachine[1].Name:
			if len(commands) < 2 {
				fmt.Println("Introduceti cuvantul")
				continue
			}

			if mch == nil {
				fmt.Println("Construiti o masina intai!")
				continue
			}
			accepted := mch.ParseWord(commands[1])
			if accepted {
				fmt.Println(commands[1], "este acceptat")
			} else {
				fmt.Println(commands[1], "nu este acceptat")
			}
			break
		case commandsMachine[2].Name:
			mch.Print()
			break
		case commandsMachine[3].Name:

			run = false
			fmt.Println("Program inchis")

			break
		case commandsMachine[4].Name:
			c := exec.Command("clear")
			c.Stdout = os.Stdout
			c.Run()
			break
		case commandsMachine[5].Name:
			fmt.Println("Comenzi:")
			for _, cmd := range commandsMachine {
				fmt.Println(cmd.ToString())
			}
			break
		case commandsMachine[6].Name:
			if mch == nil {
				fmt.Println("Construiti o masina intai!")
				continue
			}
			if len(commands) > 1 {
				mch.SetPrint(commands[1] == "true")
			} else {
				fmt.Println("Comanda print are nevoie de al 2-lea argument true/false")
			}
		default:
			fmt.Println("Comanda nerecunoscuta")
		}

		time.Sleep(100 * time.Millisecond)

	}
}
