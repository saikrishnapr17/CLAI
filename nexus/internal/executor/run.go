package executor

import (
    "errors"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

func SplitCommand(command string) []string {
    var parts []string
    var cur strings.Builder
    inSingle := false
    inDouble := false
    escape := false

    for _, r := range command {
        if escape {
            cur.WriteRune(r)
            escape = false
            continue
        }
        if r == '\\' {
            escape = true
            continue
        }
        if r == '\'' && !inDouble {
            inSingle = !inSingle
            continue
        }
        if r == '"' && !inSingle {
            inDouble = !inDouble
            continue
        }
        if r == ' ' && !inSingle && !inDouble {
            if cur.Len() > 0 {
                parts = append(parts, cur.String())
                cur.Reset()
            }
            continue
        }
        cur.WriteRune(r)
    }
    if cur.Len() > 0 {
        parts = append(parts, cur.String())
    }
    return parts
}

func Execute(command string) error {
    args := SplitCommand(command)
    if len(args) == 0 {
        return errors.New("empty command")
    }

    cmd := exec.Command(args[0], args[1:]...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    cmd.Dir, _ = os.Getwd()

    if err := cmd.Start(); err != nil {
        if errors.Is(err, exec.ErrNotFound) {
            return errors.New("command not found: " + args[0])
        }
        return err
    }

    if err := cmd.Wait(); err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok {
            if exitErr.ProcessState != nil {
                code := exitErr.ProcessState.ExitCode()
                return errors.New("exit status " + strconv.Itoa(code))
            }
        }
        return err
    }

    return nil
}
