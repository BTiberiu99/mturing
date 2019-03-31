package machine

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	//BLANK ... Blank item of a band
	BLANK = "B"
)

type Machine struct {
	Begin       string
	End         []string
	Transitions []Transition
	Print       bool
}

func (m *Machine) ConstructMachine(nameFile string) error {
	file, err := os.Open(nameFile)
	if err != nil {
		return err
	}
	defer file.Close()

	m.Transitions = make([]Transition, 0)

	scanner := bufio.NewScanner(file)
	nrRow := 0

	//Take Begin Transition
	nrRow++
	if scanner.Scan() {
		m.Begin = scanner.Text()
	} else {
		fmt.Println("Nu exista niciun rand de parsat in " + nameFile)
	}

	if scanner.Scan() {
		m.End = strings.Split(scanner.Text(), ",")
	}
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		nrRow++
		row := match(scanner.Text())

		items := strings.Split(row, ",")
		items2 := strings.Split(items[1], "")
		if len(items) == 3 && len(items2) == 3 {
			dir := int8(-1)
			if items2[2] == "R" {
				dir = 2
			} else if items2[2] == "S" {
				dir = 1
			} else if items2[2] == "L" {
				dir = 0
			}
			m.Transitions = append(m.Transitions, Transition{
				Begin: items[0],
				End:   items[2],
				Do: TuringMovement{
					Read:      items2[0],
					Write:     items2[1],
					Direction: dir,
				},
			})
		} else {
			fmt.Println("Eroare pe randul " + fmt.Sprintf("%d", nrRow) + " tranzitia trebuie sa fie de forma (q0,BBR,q1) - (trans, ReadWriteDirection ,trans)")

		}
	}

	return scanner.Err()
}

func match(s string) string {
	i := strings.Index(s, "(")
	if i >= 0 {
		j := strings.Index(s[i:], ")")
		if j >= 0 {
			return s[i+1 : j+i]
		}
	}
	return ""
}

func (m *Machine) PrintMachine() {

	fmt.Printf("Inceput:%s\n", m.Begin)
	fmt.Println("F=", m.End)
	for _, trans := range m.Transitions {
		move := trans.Do.ToString()
		fmt.Println("(" + trans.Begin + "," + move + "," + trans.End + ")")
	}
}

func (m *Machine) ParseWord(word string) bool {
	if m.Print {
		fmt.Println("WORD:", word)
	}
	accepted := false
	stop := 0
	maximStop := 100000
	letters := strings.Split(word, "")
	lenWord := len(letters)
	banda := make([]string, lenWord+2)
	for i := range letters {
		banda[i+1] = letters[i]
	}
	banda[0] = BLANK
	banda[lenWord+1] = BLANK
	j := 0
	var currentTransition Transition
	for _, trans := range m.Transitions {
		if trans.Begin == m.Begin {
			currentTransition = trans
			break
		}
	}
	accepted = currentTransition.Final(m.End)
	for !accepted {
		stop++
		if stop == maximStop {
			fmt.Println("Este posibil sa se fi intrat intr-un loop infinit")
			break
		}

		loop := true
		for !accepted && loop {
			loop = currentTransition.End == currentTransition.Begin
			if currentTransition.Do.Read == banda[j] {
				banda[j] = currentTransition.Do.Write

				if currentTransition.Do.Direction == 0 {
					j--
				} else if currentTransition.Do.Direction == 2 {
					j++
				}
				if m.Print {
					fmt.Println(banda, "INDEX", j)
				}
			} else {
				break
			}
		}
		next := false
		for i := range m.Transitions {
			if banda[j] == m.Transitions[i].Do.Read && currentTransition.End == m.Transitions[i].Begin {
				if m.Print {
					fmt.Println("NEW TRANSITION", m.Transitions[i])
				}
				currentTransition = m.Transitions[i]
				next = true
				break
			}
		}
		if !next {
			break
		}
		accepted = currentTransition.Final(m.End)

	}

	return accepted
}

type Transition struct {
	Begin string
	End   string
	Do    TuringMovement
}

func (t *Transition) Final(finalTrans []string) bool {
	for i := range finalTrans {
		if t.Begin == finalTrans[i] || t.End == finalTrans[i] {
			return true
		}
	}
	return false
}

type TuringMovement struct {
	Direction int8 // true = "RIGHT" , false ="LEFT"
	Read      string
	Write     string
}

func (mov *TuringMovement) ToString() string {
	move := mov.Read + mov.Write
	if mov.Direction == 0 {
		move += "R"
	} else if mov.Direction == 1 {
		move += "S"
	} else if mov.Direction == 2 {
		move += "L"
	}
	return move
}
