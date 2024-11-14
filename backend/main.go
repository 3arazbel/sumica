package main

import (
    "encoding/json"
    "log"
    "net/http"
    "html/template"
    "bytes"
    "os"

    "github.com/go-resty/resty/v2"
)

type Song struct {
    Title        string `json:"title"`
    Artist       string `json:"artist"`
    Field        string `json:"field"`
    Id           string `json:"id"`
    CollectionId string `json:"collectionId"`
}

type PocketBaseResponse struct {
    Items []Song `json:"items"`
}

const songTemplate = `
{{range .Items}}
<div class="song-item">
    <h3>{{.Title}}</h3>
    <p>by {{.Artist}}</p>
    <audio controls>
        <source src="http://localhost:8090/api/files/{{.CollectionId}}/{{.Id}}/{{.Field}}" type="audio/mp3">
        Your browser does not support the audio element.
    </audio>
</div>
{{end}}
`

func getSongs(w http.ResponseWriter, r *http.Request) {
    client := resty.New()
    pocketbaseURL := os.Getenv("POCKETBASE_URL")
    if pocketbaseURL == "" {
        pocketbaseURL = "http://pocketbase:8090"
    }
    resp, err := client.R().Get(pocketbaseURL + "/api/collections/songs/records")
    if err != nil {
        http.Error(w, "Failed to fetch songs", http.StatusInternalServerError)
        log.Println("Error fetching songs:", err)
        return
    }

    var pocketbaseResp PocketBaseResponse
    if err := json.Unmarshal(resp.Body(), &pocketbaseResp); err != nil {
        http.Error(w, "Error parsing JSON", http.StatusInternalServerError)
        log.Println("JSON Unmarshal error:", err)
        return
    }

    tmpl, err := template.New("songs").Parse(songTemplate)
    if err != nil {
        http.Error(w, "Template error", http.StatusInternalServerError)
        log.Println("Template error:", err)
        return
    }

    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, pocketbaseResp); err != nil {
        http.Error(w, "Template error", http.StatusInternalServerError)
        log.Println("Template error:", err)
        return
    }

    w.Header().Set("Content-Type", "text/html")
    w.Write(buf.Bytes())
}

func main() {
    http.Handle("/", http.FileServer(http.Dir("/app/frontend")))
    http.HandleFunc("/getSongs", getSongs)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server running on http://localhost:%s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
