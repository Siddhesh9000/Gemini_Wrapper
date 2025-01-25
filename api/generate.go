package generate

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// HomeHandler serves the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/index.html"))
	tmpl.Execute(w, nil)
}

// GenerateHandler handles the content generation request
func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	prompt := r.FormValue("prompt")
	if prompt == "" {
		http.Error(w, "Prompt is required", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		http.Error(w, "API key not set", http.StatusInternalServerError)
		log.Println("GEMINI_API_KEY environment variable is not set")
		return
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		http.Error(w, "Failed to create client", http.StatusInternalServerError)
		log.Printf("Failed to create client: %v\n", err)
		return
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-pro")

	var promptParts []genai.Part
	if file, _, err := r.FormFile("image"); err == nil {
		defer file.Close()
		imgData, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read image", http.StatusInternalServerError)
			log.Printf("Failed to read image: %v\n", err)
			return
		}
		promptParts = append(promptParts, genai.ImageData("jpeg", imgData))
	}

	promptParts = append(promptParts, genai.Text(prompt))
	resp, err := model.GenerateContent(ctx, promptParts...)

	if err != nil {
		http.Error(w, "Failed to generate content", http.StatusInternalServerError)
		log.Printf("Failed to generate content: %v\n", err)
		return
	}

	if resp != nil && resp.Candidates != nil && len(resp.Candidates) > 0 {
		content := resp.Candidates[0].Content
		if content != nil && len(content.Parts) > 0 {
			result := content.Parts[0]
			fmt.Fprintf(w, "Generated Content: %s", result)
			return
		}
	}

	fmt.Fprint(w, "No content generated")
}

// Handler function for API endpoint, called by Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Route requests based on the path
	switch r.URL.Path {
	case "/":
		HomeHandler(w, r)
	case "/generate":
		GenerateHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}
