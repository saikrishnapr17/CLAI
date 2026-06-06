package llm

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "time"
)

type GroqClient struct {
    APIKey     string
    Model      string
    HTTPClient *http.Client
}

func NewGroqClient(apiKey string) *GroqClient {
    return &GroqClient{
        APIKey: apiKey,
        Model:  "llama-4-scout-17b-16e-instruct",
        HTTPClient: &http.Client{
            Timeout: 15 * time.Second,
        },
    }
}

type chatMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type chatRequest struct {
    Model       string        `json:"model"`
    Messages    []chatMessage `json:"messages"`
    Temperature float32       `json:"temperature"`
    MaxTokens   int           `json:"max_tokens"`
}

type chatChoice struct {
    Message chatMessage `json:"message"`
}

type chatResponse struct {
    ID      string       `json:"id"`
    Object  string       `json:"object"`
    Created int64        `json:"created"`
    Choices []chatChoice `json:"choices"`
    Error   interface{}  `json:"error"`
}

func (c *GroqClient) Translate(userInput string, osContext string) (string, error) {
    reqBody := chatRequest{
        Model: c.Model,
        Messages: []chatMessage{
            {Role: "system", Content: osContext},
            {Role: "user", Content: fmt.Sprintf("Translate this to a shell command: %s", userInput)},
        },
        Temperature: 0.1,
        MaxTokens:   200,
    }

    data, err := json.Marshal(reqBody)
    if err != nil {
        return "", err
    }

    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewReader(data))
    if err != nil {
        return "", err
    }
    req.Header.Set("Authorization", "Bearer "+c.APIKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        // detect timeout
        if errors.Is(err, context.DeadlineExceeded) {
            return "", fmt.Errorf("request timed out")
        }
        return "", err
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return "", fmt.Errorf("groq API error: %s", string(body))
    }

    var cr chatResponse
    if err := json.Unmarshal(body, &cr); err != nil {
        return "", fmt.Errorf("failed to parse groq response: %w", err)
    }

    if len(cr.Choices) == 0 {
        return "", fmt.Errorf("no choices returned from Groq")
    }

    return cr.Choices[0].Message.Content, nil
}
