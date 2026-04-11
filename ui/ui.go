package ui

import (
	"fmt"
	"strings"
)

type CallbackMenu func()

type MenuCmd struct {
	Name string
	Callback CallbackMenu
}

func startCallback() {
	fmt.Println("Starting...")
}

func exitCallback() {
	fmt.Println("Exiting...")
}

func Clear() {
	fmt.Print("\033[H\033[2J")
}

/*

func saveCursor() {
	fmt.Print("\033[s")
}

func restoreCursor() {
	fmt.Print("\033[u")
}

func moveCursorAbs(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

func getCursorPos() (row, col int) {
	fmt.Print("\033[6n")

	reader := bufio.NewReader(os.Stdin)
	var ch byte
	for {
		b, _ := reader.ReadByte()
		if b == 0x1b {
			ch, _ = reader.ReadByte()
			if ch == '[' {
				break
			}
		}
	}

	fmt.Fscanf(reader, "%d;%dR", &row, &col)

	return
}

func getConsoleDim() (row, col int) {
	saveCursor()

	moveCursorAbs(1000, 1000)	
	row, col = getCursorPos()

	restoreCursor()

	return
}
*/

func HandleIntInput(lower, upper int) int {
	input := 0
	for {
		fmt.Print("> ")
		_, err := fmt.Scanf("%d", &input)
		if err == nil{
			if  input < lower || input > upper {
				fmt.Println("Input out of bounds")
				continue
			}
			break
		}
		fmt.Println("Invalid input")
	}

	return input
}

func Menu() {
	Clear()
	options := map[int]MenuCmd{
		1: {
			Name: "Start Game",
			Callback: startCallback,
		},
		2: {
			Name: "Exit",
			Callback: exitCallback,
		},
	}

	var cmdLines strings.Builder
	optionsNums := []int { 1, 2 }
	for _, v := range optionsNums {
		fmt.Fprintf(&cmdLines, "%d. %s\n", v, options[v].Name)
	}

	fmt.Println("+=== Main Menu ===+")	
	fmt.Print(cmdLines.String())
	fmt.Println("+=================+")

	input := HandleIntInput(1, 2)
	
	options[input].Callback()
}

func WaitingScreen(minPlayers, maxPlayers, currPlayers int) {
	Clear()

	remaining := maxPlayers - currPlayers
	ready := currPlayers > minPlayers
	
	fmt.Printf("%d/%d\n", currPlayers, maxPlayers)
	
	if remaining >= 0 {
		fmt.Printf("Remainig %d\n", remaining)
	}

	if ready{
		fmt.Println("Can Start")
	}
}
