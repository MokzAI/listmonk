package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// L represents a translation system
type L struct{}

// Store translations in a package-level variable
var translations = map[string]string{
	"email.optin.confirmSubWelcome": "Hi",
	"email.optin.confirmSub":        "Confirm subscription",
}

// Update the Ts method to use the translations map
func (l L) Ts(key string) string {
	if val, ok := translations[key]; ok {
		return val
	}
	return key
}

func (l L) T(key string) string {
	return l.Ts(key)
}

func main() {
	// Define configuration constants at the top of main
	const (
		baseURL = "http://localhost:8000"
		logoURL = baseURL + "/public/static/logo.svg"
	)

	// Template handler for email templates
	emailHandler := func(w http.ResponseWriter, r *http.Request) {
		// Extract template name from URL path
		rawTemplateName := strings.TrimPrefix(r.URL.Path, "/email-templates/")
		rawTemplateName = strings.TrimSuffix(rawTemplateName, ".html")
		
		if rawTemplateName == "" {
			// Show list of available templates
			files, err := filepath.Glob("static/email-templates/*.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Create simple HTML list of templates with better styling
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `
				<html>
				<head>
					<style>
						body { font-family: Arial, sans-serif; margin: 20px; }
						h1 { color: #333; }
						ul { list-style-type: none; padding: 0; }
						li { margin: 10px 0; }
						a { color: #0055d4; text-decoration: none; }
						a:hover { text-decoration: underline; }
					</style>
				</head>
				<body>
					<h1>/email-templates</h1>
					<ul>
			`)
			
			for _, file := range files {
				name := filepath.Base(file)
				if name != "base.html" {  // Don't show base.html
					fmt.Fprintf(w, `<li><a href="/email-templates/%s">%s</a></li>`, name, name)
				}
			}
			
			fmt.Fprintf(w, `
					</ul>
				</body>
				</html>
			`)
			return
		}

		// Map of filename (without .html) to template definition names
		templateNameMap := map[string]string{
			"subscriber-optin-campaign": "optin-campaign",
			// Add other mappings as needed
		}

		// Get the correct template name to execute
		templateName := templateNameMap[rawTemplateName]
		if templateName == "" {
			templateName = rawTemplateName // fallback to the original name
		}

		if templateName == "" {
			// Show list of available templates
			files, err := filepath.Glob("static/email-templates/*.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Create simple HTML list of templates
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "<h1>Available Templates:</h1><ul>")
			for _, file := range files {
				name := filepath.Base(file)
				if name != "base.html" {  // Don't show base.html as it's not meant to be rendered directly
					fmt.Fprintf(w, `<li><a href="/email-templates/%s">%s</a></li>`, name, name)
				}
			}
			fmt.Fprintf(w, "</ul>")
			return
		}

		// Get all template files
		files, err := filepath.Glob("static/email-templates/*.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create template with functions
		tmpl := template.New("")
		
		// Add all necessary template functions
		tmpl = tmpl.Funcs(template.FuncMap{
			"L": func() L { return L{} },
			"RootURL": func() string { 
				return "http://localhost:8000" 
			},
			"LogoURL": func() string {
				return logoURL
			},
			"Safe": func(s string) template.HTML {
				return template.HTML(s)
			},
			"Date": func(t interface{}) string {
				return fmt.Sprintf("%v", t)
			},
			"UnixTime": func(t interface{}) int64 {
				return 0 // placeholder
			},
			"ToUpper": strings.ToUpper,
			"ToLower": strings.ToLower,
			"ne": func(a, b interface{}) bool {
				return a != b
			},
			"eq": func(a, b interface{}) bool {
				return a == b
			},
		})

		// Parse all templates
		tmpl = template.Must(tmpl.ParseFiles(files...))

		// Template data with all possible fields
		data := struct {
			L          L
			LogoURL    string
			Lists      []struct{ Name, Type string }
			Subscriber struct{ 
				FirstName string
				Email     string
				UUID     string
			}
			OptinURL    string
			UnsubURL    string
			Campaign    struct{ Subject string }
			Links       struct{ Manage string }
			Campaigns   []struct{ Subject string }
			Token      string
			SiteURL    string
			MessageURL string
		}{
			L:       L{},
			LogoURL: logoURL,
			Lists: []struct{ Name, Type string }{
				{Name: "Test List", Type: "public"},
			},
			Subscriber: struct{ 
				FirstName string
				Email     string
				UUID     string
			}{
				FirstName: "John",
				Email:     "john@example.com",
				UUID:     "test-uuid",
			},
			OptinURL:    "#",
			UnsubURL:    "#",
			Campaign:    struct{ Subject string }{Subject: "Test Campaign"},
			Links:       struct{ Manage string }{Manage: "#"},
			Campaigns:   []struct{ Subject string }{{Subject: "Test Campaign"}},
			Token:      "test-token",
			SiteURL:    "http://localhost:8000",
			MessageURL: "#",
		}

		// Debug information in case of error
		err = tmpl.ExecuteTemplate(w, templateName, data)
		if err != nil {
			// List all available templates
			var availableTemplates []string
			for _, t := range tmpl.Templates() {
				availableTemplates = append(availableTemplates, t.Name())
			}
			
			errMsg := fmt.Sprintf("Template error: %v\nRequested template: %s\nAvailable templates: %v", 
				err, templateName, availableTemplates)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}
	}

	// Create a custom file server handler that enables directory listing
	publicTemplatesHandler := func(w http.ResponseWriter, r *http.Request) {
		// Strip /public/templates from the path
		path := strings.TrimPrefix(r.URL.Path, "/public/templates")
		if path == "" || path == "/" {
			// Show directory listing with debug info
			files, err := filepath.Glob("static/public/templates/*.html")
			if err != nil {
				// Show the error and current working directory
				currentDir, _ := os.Getwd()
				errMsg := fmt.Sprintf("Error: %v\nCurrent working directory: %s\n", err, currentDir)
				http.Error(w, errMsg, http.StatusInternalServerError)
				return
			}

			// Try to read the directory directly
			entries, err := os.ReadDir("static/public/templates")
			if err != nil {
				currentDir, _ := os.Getwd()
				errMsg := fmt.Sprintf("Error reading directory: %v\nCurrent working directory: %s\n", err, currentDir)
				http.Error(w, errMsg, http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `
				<html>
				<head>
					<style>
						body { font-family: Arial, sans-serif; margin: 20px; }
						h1, h2 { color: #333; }
						ul { list-style-type: none; padding: 0; }
						li { margin: 10px 0; }
						a { color: #0055d4; text-decoration: none; }
						a:hover { text-decoration: underline; }
						.debug { background: #f5f5f5; padding: 10px; margin: 20px 0; }
					</style>
				</head>
				<body>
					<h1>Public Templates:</h1>
					<div class="debug">
						<h2>Debug Information:</h2>
						<p>Looking for files in: static/public/templates/*.html</p>
						<p>Directory contents:</p>
						<ul>
			`)
			
			// List all files in the directory
			for _, entry := range entries {
				fmt.Fprintf(w, `<li>%s (IsDir: %v)</li>`, entry.Name(), entry.IsDir())
			}
			
			fmt.Fprintf(w, `
						</ul>
					</div>
					<h2>Template Files:</h2>
					<ul>
			`)
			
			for _, file := range files {
				name := filepath.Base(file)
				fmt.Fprintf(w, `<li><a href="/public/templates/%s">%s</a></li>`, name, name)
			}
			
			fmt.Fprintf(w, `
					</ul>
				</body>
				</html>
			`)
			return
		}

		// Serve the actual file
		filePath := filepath.Join("static/public/templates", path)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, fmt.Sprintf("File not found: %s", filePath), http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, filePath)
	}

	// Create a handler for /public/ that shows both directories and files
	publicHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/public/" {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `
				<html>
				<head>
					<meta charset="UTF-8">
					<style>
						body { font-family: Arial, sans-serif; margin: 20px; }
						h1, h2 { color: #333; }
						ul { list-style-type: none; padding: 0; }
						li { margin: 10px 0; }
						a { color: #0055d4; text-decoration: none; }
						a:hover { text-decoration: underline; }
						.section { margin-bottom: 30px; }
					</style>
				</head>
				<body>
					<h1>Public Directory</h1>
			`)

			// Show static directory contents
			staticFiles, err := filepath.Glob("static/public/static/*")
			if err == nil && len(staticFiles) > 0 {
				fmt.Fprintf(w, `
					<div class="section">
						<h2>/static</h2>
						<ul>
				`)
				for _, file := range staticFiles {
					name := filepath.Base(file)
					fmt.Fprintf(w, `<li><a href="/public/static/%s">%s</a></li>`, name, name)
				}
				fmt.Fprintf(w, `</ul></div>`)
			}

			// Show template files
			templateFiles, err := filepath.Glob("static/public/templates/*.html")
			if err == nil && len(templateFiles) > 0 {
				fmt.Fprintf(w, `
					<div class="section">
						<h2>/templates</h2>
						<ul>
				`)
				for _, file := range templateFiles {
					name := filepath.Base(file)
					fmt.Fprintf(w, `<li><a href="/public/templates/%s">%s</a></li>`, name, name)
				}
				fmt.Fprintf(w, `</ul></div>`)
			}

			fmt.Fprintf(w, `</body></html>`)
			return
		}
		
		// If not the root /public/ path, serve the file
		http.ServeFile(w, r, filepath.Join("static", r.URL.Path))
	}

	// Route handlers
	http.HandleFunc("/email-templates/", emailHandler)
	http.HandleFunc("/public/templates/", publicTemplatesHandler)
	http.HandleFunc("/public/", publicHandler)
	http.Handle("/", http.FileServer(http.Dir("static")))  // Keep this last

	log.Print("[Server] 🏃 http://localhost:8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}