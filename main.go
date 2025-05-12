package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// folder for the notes
const NotesDir = ".tnotes"

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	// Ensure hidden notes directory exists
	if err := os.MkdirAll(NotesDir, 0755); err != nil {
		fmt.Printf("Error creating notes dir: %v\n", err)
		return
	}

	// Hide the folder on Windows
	if runtime.GOOS == "windows" {
		hideWindowsDir(NotesDir)
	}

	// Handle commands
	switch {
	case os.Args[1] == "-v" || os.Args[1] == "-view":
		listNotes()
	case os.Args[1] == "-d" || os.Args[1] == "-delete":
		if len(os.Args) < 3 {
			fmt.Println("Please specify note ID to delete")
			return
		}
		deleteNote(os.Args[2])
	case os.Args[1] == "-help" || os.Args[1] == "-h":
		showHelp()
	case strings.HasPrefix(os.Args[1], "-n"):
		handleTerminalCapture()
	default:
		// If first arg starts with '-', but is not a recognized command, treat it as invalid
		if strings.HasPrefix(os.Args[1], "-") {
			fmt.Println("Invalid command. Use 'tnote -help' for usage")
			return
		}
		// treat the rest of the args as a message
		// if the first arg is not a recognized command, treat it as a message
		message := strings.Join(os.Args[1:], " ")
		saveNote(message, "")
	}
}

// Here for future implementation
func handleTerminalCapture() {
	if len(os.Args) < 3 {
		fmt.Println("Please specify line count and text")
		return
	}

	lineCountStr := strings.TrimPrefix(os.Args[1], "-n")
	lineCount, err := strconv.Atoi(lineCountStr)
	if err != nil || lineCount < 1 {
		fmt.Println("Invalid line count. Use -n<number> where number >= 1")
		return
	}

	message := strings.Join(os.Args[2:], " ")
	terminalHistory := captureTerminalHistory(lineCount)
	fullNote := fmt.Sprintf("%s\n\nTerminal context:\n%s", message, terminalHistory)
	saveNote(fullNote, "term_")
}

func captureTerminalHistory(lines int) string {
	//if runtime.GOOS == "windows" {
	//	return captureWindowsHistory(lines)
	//}
	//return captureUnixTerminalHistory(lines)

	// return a placeholder until I figure out if and how its feasable to do it
	fmt.Println("terminal history capture not implemented yet, couldn't get", lines, "lines")
	return fmt.Sprintf("[Terminal history capture not implemented for %s]", runtime.GOOS)
}

// // this DOES NOT WORK
// func captureWindowsHistory(lines int) string {
//	if history := getCurrentSessionHistory(lines); history != "" {
//		return history
//	}

//
// 	if history := readWindowsTerminalHistory(); history != "" {
// 		return getLastLines(history, lines)
// 	}

// 	return "[Could not capture history - try running in PowerShell]"
// }

// func getCurrentSessionHistory(lines int) string {
// 	// Create a temporary PowerShell script
// 	script := fmt.Sprintf(`
//     $history = Get-History -Count %d | ForEach-Object { $_.CommandLine }
//     $history | Out-File -FilePath "$env:TEMP\\tnote_history.txt" -Encoding utf8
//     `, lines)

// 	// Execute the script
// 	cmd := exec.Command("powershell", "-Command", script)
// 	if err := cmd.Run(); err != nil {
// 		return ""
// 	}

// 	// Read the output file
// 	content, err := os.ReadFile(os.Getenv("TEMP") + "\\tnote_history.txt")
// 	if err != nil {
// 		return ""
// 	}

// 	return string(content)
// }

// func readWindowsTerminalHistory() string {
// 	// Try multiple possible history file locations
// 	paths := []string{
// 		filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Windows Terminal", "State", "history.json"),
// 		filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Microsoft", "Windows", "PowerShell", "history.json"),
// 	}

// 	for _, path := range paths {
// 		if content, err := os.ReadFile(path); err == nil {
// 			return string(content)
// 		}
// 	}
// 	return ""
// }

// func getLastLines(text string, n int) string {
// 	lines := strings.Split(strings.TrimSpace(text), "\n")
// 	start := len(lines) - n
// 	if start < 0 {
// 		start = 0
// 	}
// 	return strings.Join(lines[start:], "\n")
// }

// func captureUnixTerminalHistory(lines int) string {
// 	// Unix implementation (Linux/Mac)
// 	home, err := os.UserHomeDir()
// 	if err != nil {
// 		return fmt.Sprintf("[Failed to get home directory: %v]", err)
// 	}

// 	historyFile := filepath.Join(home, ".bash_history")
// 	if runtime.GOOS == "darwin" {
// 		// On macOS, check .zsh_history if using zsh
// 		shell := os.Getenv("SHELL")
// 		if strings.Contains(shell, "zsh") {
// 			historyFile = filepath.Join(home, ".zsh_history")
// 		}
// 	}

// 	file, err := os.Open(historyFile)
// 	if err != nil {
// 		return fmt.Sprintf("[Failed to open history file: %v]", err)
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	var history []string

// 	// Read all lines into memory (history files are typically small)
// 	for scanner.Scan() {
// 		history = append(history, scanner.Text())
// 	}

// 	if err := scanner.Err(); err != nil {
// 		return fmt.Sprintf("[Error reading history: %v]", err)
// 	}

// 	// Get the last 'lines' entries
// 	start := len(history) - lines
// 	if start < 0 {
// 		start = 0
// 	}

// 	var result strings.Builder
// 	for i := start; i < len(history); i++ {
// 		result.WriteString(history[i])
// 		result.WriteString("\n")
// 	}

// 	return result.String()
// }

func hideWindowsDir(dir string) error {
	// Since apparently ".path" does not hide folders on windows by default
	// we need to set the hidden attribute manually
	dirW, err := syscall.UTF16PtrFromString(dir)
	if err != nil {
		return err
	}
	return syscall.SetFileAttributes(dirW, syscall.FILE_ATTRIBUTE_HIDDEN)
}

func saveNote(message string, prefix string) {
	// save the note with a timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := filepath.Join(NotesDir, prefix+timestamp)

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating note: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(message)
	if err != nil {
		fmt.Printf("Error writing note: %v\n", err)
		return
	}

	fmt.Printf("Note saved: %s\n", filename)
}

func listNotes() {
	files, err := os.ReadDir(NotesDir)
	if err != nil {
		fmt.Printf("No notes found (dir: %s)\n", NotesDir)
		return
	}

	if len(files) == 0 {
		fmt.Println("No notes found")
		return
	}

	fmt.Println("Your tnotes:")
	fmt.Println(strings.Repeat("=", 50))

	for i, file := range files {
		content, _ := os.ReadFile(filepath.Join(NotesDir, file.Name()))
		fmt.Printf("[%d] %s\n", i+1, file.Name())
		fmt.Println(string(content))
		fmt.Println(strings.Repeat("-", 50))
	}
}

func deleteNote(id string) {
	// this is a bit hacky as it relies on the fact that the files are sorted by name
	// but it works for now
	// in the future I might want to add a more robust way to handle this
	files, err := os.ReadDir(NotesDir)
	if err != nil {
		fmt.Printf("No notes found (dir: %s)\n", NotesDir)
		return
	}

	index := 0
	_, err = fmt.Sscan(id, &index)
	if err != nil || index < 1 || index > len(files) {
		fmt.Println("Invalid note ID")
		return
	}

	noteToDelete := files[index-1]
	err = os.Remove(filepath.Join(NotesDir, noteToDelete.Name()))
	if err != nil {
		fmt.Printf("Error deleting note: %v\n", err)
		return
	}

	fmt.Printf("Deleted note: %s\n", noteToDelete.Name())
}

func showHelp() {
	fmt.Println(`tnote - Simple Terminal Note Taking App

Usage:
  tnote <text>               Create new note with given text
  tnote -n<lines> <text>     Create note with terminal context (not implemented)
  tnote -v                   View all notes
  tnote -d <id>              Delete note with specified ID
  tnote -help                Show this help message

Examples:
  tnote Remember to fix the bug
  tnote -n5 Capture last 5 commands (not implemented)
  tnote -v
  tnote -d 2`)
}
