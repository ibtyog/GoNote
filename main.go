package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ibtyog/GoNote/note"
)

func main() {
	fmt.Println("Welcome to GoNote!")
	for {
		option := getUserInput("1. Create new note\n2. View notes\n3. Exit")
		switch option {
		case "1":
			newNote, err := createNote()
			if err != nil {
				fmt.Println("Failed to create note object.")
				return
			}
			newNote.Save()
		case "2":
			readAllNotes()
		case "3":
			fmt.Println("Thank you for using GoNote!")
			return
		default:
			fmt.Println("Please provide a valid option.")
		}
	}
}

func getUserInput(promptText string) string {
	fmt.Println(promptText)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')

	if err != nil {
		return "I/O Reader failed to read your input."
	}

	text = strings.TrimSuffix(text, "\n")
	text = strings.TrimSuffix(text, "\r")
	return text
}

func createNote() (note.Note, error) {
	title := getUserInput("Note title: ")
	for {
		if _, err := os.Stat("notes/" + title + ".json"); err == nil {
			fmt.Println("File already exists.")
			title = getUserInput("Note title: ")
		} else {
			break
		}
	}

	content := getUserInput("Note content: ")
	newNote, err := note.New(title, content)
	if err != nil {
		fmt.Println("Failed to create note.")
		return newNote, err
	}
	return newNote, nil
}

func readAllNotes() {
	file, err := os.Open("notes")
	if err != nil {
		return
	}
	defer file.Close()
	list, _ := file.Readdirnames(0)
	counter := 1
	for _, name := range list {
		fmt.Printf("%v. %v\n", counter, name)
		counter++
	}

	operation := getUserInput("Choose operation (view/quit): ")
	for {
		switch operation {
		case "view":
			i := getUserInput("Choose number of desired file: ")
			num, err := strconv.Atoi(i)

			if err != nil {
				return
			}

			if num < 1 || num > len(list) {
				fmt.Println("Provided invalid number.")
				return
			}

			viewNote(list[num-1])
			return
		case "quit":
			return
		default:
			fmt.Print("Provide valid operation ")
			operation = getUserInput("(view/quit): ")
		}

	}

}

func viewNote(name string) {
	jsonFile, err := os.Open("notes/" + name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	fmt.Println(result["content"])
}
