package machine

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	//BLANK ... Blank item of a band
	BLANK = "B"
)

type Machine struct {
	begin       string
	end         []string
	transitions []Transition
	isPrinting  bool
}

type Config struct {
	NameFile string
	Print    bool
	Buffer   io.Reader
}

func NewMachine(conf *Config) (*Machine, error) {
	machine := &Machine{
		isPrinting: conf.Print,
	}

	file := conf.Buffer

	if file == nil {
		file, err := os.Open(conf.NameFile)

		if err != nil {
			return nil, err
		}

		defer file.Close()
	}

	err := machine.constructMachine(file)

	if err != nil {
		return nil, err
	}

	return machine, nil
}

func (m *Machine) SetPrint(print bool) {
	m.isPrinting = print
}

func (m *Machine) constructMachine(file io.Reader) error {

	m.transitions = make([]Transition, 0)

	scanner := bufio.NewScanner(file)
	nrRow := 0

	//Take Begin Transition
	nrRow++
	if scanner.Scan() {
		m.begin = scanner.Text()
	} else {
		fmt.Println("Nu exista niciun rand de parsat")
	}

	if scanner.Scan() {
		m.end = strings.Split(scanner.Text(), ",")
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

			m.transitions = append(m.transitions, Transition{
				begin: items[0],
				end:   items[2],
				do: TuringMovement{
					read:      items2[0],
					write:     items2[1],
					direction: dir,
				},
			})

		} else {
			fmt.Println(fmt.Sprintf("Eroare pe randul %d tranzitia trebuie sa fie de forma (q0,BBR,q1) - (trans, ReadWriteDirection ,trans)", nrRow))

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

func (m *Machine) Print() {

	fmt.Printf("Inceput:%s\n", m.begin)
	fmt.Println("F=", m.end)
	for _, trans := range m.transitions {
		fmt.Printf("(%s,%s,%s)\n", trans.begin, trans.do.ToString(), trans.end)
	}
}

func (m *Machine) Speak(str ...interface{}) {
	if m.isPrinting {
		fmt.Println(str...)
	}
}

func (m *Machine) ParseWord(word string) bool {

	m.Speak("WORD:", word)

	var (
		accepted  = false
		stop      = 0
		maximStop = 100000
	)

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

	for _, trans := range m.transitions {
		if trans.begin == m.begin {
			currentTransition = trans
			break
		}
	}

	accepted = currentTransition.Final(m.end)

	for !accepted {
		stop++
		if stop == maximStop {
			fmt.Println("Este posibil sa se fi intrat intr-un loop infinit")
			break
		}

		loop := true
		for !accepted && loop {
			loop = currentTransition.end == currentTransition.begin
			if currentTransition.do.read == banda[j] {
				banda[j] = currentTransition.do.write

				if currentTransition.do.direction == 0 {
					j--
				} else if currentTransition.do.direction == 2 {
					j++
				}

				m.Speak(banda, "INDEX", j)
			} else {
				break
			}
		}
		next := false

		for i := range m.transitions {
			if banda[j] == m.transitions[i].do.read && currentTransition.end == m.transitions[i].begin {

				m.Speak("NEW TRANSITION", m.transitions[i])

				currentTransition = m.transitions[i]
				next = true
				break
			}
		}

		if !next {
			break
		}

		accepted = currentTransition.Final(m.end)

	}

	return accepted
}

type Transition struct {
	begin string
	end   string
	do    TuringMovement
}

func (t *Transition) Final(finalTrans []string) bool {
	for i := range finalTrans {
		if t.begin == finalTrans[i] || t.end == finalTrans[i] {
			return true
		}
	}
	return false
}

type TuringMovement struct {
	direction int8 // true = "RIGHT" , false ="LEFT"
	read      string
	write     string
}

func (mov *TuringMovement) ToString() string {

	move := mov.read + mov.write

	if mov.direction == 0 {
		move += "R"
	} else if mov.direction == 1 {
		move += "S"
	} else if mov.direction == 2 {
		move += "L"
	}

	return move
}
