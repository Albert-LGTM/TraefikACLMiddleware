package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
)

func main() {
    middleware := func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            authAPIEndpoint := "http://100.113.188.36:4444/api/authenticated"

            // Make a request to the authentication API endpoint
            resp, err := http.Get(authAPIEndpoint)
            if err != nil {
                http.Error(w, "Failed to check authentication status", http.StatusInternalServerError)
                return
            }
            defer resp.Body.Close()

            // Ensure successful HTTP response
            if resp.StatusCode != http.StatusOK {
                http.Error(w, fmt.Sprintf("Unexpected status code from authentication API: %d", resp.StatusCode), http.StatusInternalServerError)
                return
            }

            // Read and unmarshal JSON response
            body, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                http.Error(w, "Failed to read authentication status", http.StatusInternalServerError)
                return
            }

            var authResponse struct {
                Authenticated bool `json:"authenticated"`
            }
            if err := json.Unmarshal(body, &authResponse); err != nil {
                http.Error(w, "Invalid JSON response from authentication API", http.StatusBadRequest)
                return
            }

            // Handle authentication result
            if authResponse.Authenticated {
                next(w, r) // Allow the request to proceed
            } else {
                http.Error(w, "Unauthorized access", http.StatusUnauthorized)
            }
        }
    }

    http.Handle("/", middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Protected content")
    })))

    port := os.Getenv("PORT")
    if port == "" {
        port = "80"
    }

    err := http.ListenAndServe(":"+port, nil)
    if err != nil {
        fmt.Println("Error starting server:", err)
    }
}

