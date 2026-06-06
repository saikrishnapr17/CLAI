package llm

import (
    "fmt"
)

type OSContext struct {
    OS         string
    Shell      string
    CurrentDir string
    HomeDir    string
}

func BuildSystemPrompt(osctx OSContext) string {
    return fmt.Sprintf(`You are a terminal assistant. Convert the user's natural language 
request into a single shell command.

Rules:
- Return ONLY the raw shell command. No explanation, no markdown, 
  no backticks, no preamble.
- The command must be valid for %s using %s
- Current directory is: %s
- Home directory is: %s
- If the request is ambiguous, return the safest interpretation
- If the request cannot be converted to a shell command, return 
  exactly: CANNOT_TRANSLATE
- Never return commands that delete system files, format drives, 
  or could cause irreversible system damage. Return UNSAFE instead.
- Prefer concise commands over verbose ones
- Use pipes and flags where appropriate
`, osctx.OS, osctx.Shell, osctx.CurrentDir, osctx.HomeDir)
}
