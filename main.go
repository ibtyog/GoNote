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
	"github.com/inancgumus/screen"
)

func main() {
	fmt.Println("Welcome to GoNote!")
	for {
		option := getUserInput("1. Create new note\n2. View notes\n3. Exit\nYour choice:")
		screen.Clear()
		switch option {
		case "1":
			newNote, err := createNote()
			if err != nil {
				fmt.Println("Failed to create note object.")
				return
			}
			newNote.Save()
			screen.Clear()
		case "2":
			readAllNotes()
			screen.Clear()
		case "3":
			fmt.Println("Thank you for using GoNote!")
			return
		default:
			fmt.Println("Please provide a valid option.")
		}
	}
}

func getUserInput(promptText string) string {
	fmt.Print(promptText)
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
			screen.Clear()
			viewNote(list[num-1])
			operation = getUserInput("Choose operation (edit/remove/quit): ")
			if operation == "edit" {
				content := getUserInput("Note content: ")
				noteTitle := strings.ReplaceAll(list[num-1], ".json", "")
				editNote, err := note.New(noteTitle, content)
				if err != nil {
					return
				}
				editNote.Save()
			} else if operation == "remove" {
				os.Remove("notes/" + list[num-1])
			} else {
				return
			}

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
