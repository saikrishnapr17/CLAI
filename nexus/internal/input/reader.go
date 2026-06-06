package input

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func ReadInput(prompt string) (string, error) {
    fmt.Print(prompt)
    scanner := bufio.NewScanner(os.Stdin)
    if !scanner.Scan() {
        if err := scanner.Err(); err != nil {
            return "", err
        }
        return "", nil
    }
    return strings.TrimSpace(scanner.Text()), nil
}

func Confirm(question string) bool {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Printf("%s [y/n]: ", question)
        line, err := reader.ReadString('\n')
        if err != nil {
            return false
        }
        v := strings.TrimSpace(strings.ToLower(line))
        if v == "y" || v == "yes" {
            return true
        }
        if v == "n" || v == "no" {
            return false
        }
        fmt.Println("Please type y or n.")
    }
}
