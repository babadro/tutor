package middlewares

import "net/http"

// Cors adds CORS headers to the response.
func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust the port as per your Flutter app
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		// Check if it's a preflight request
		if r.Method == http.MethodOptions {
			// Preflight request; respond with 200 OK without further processing
			w.WriteHeader(http.StatusOK)
			return
		}

		// Handle the actual request
		next.ServeHTTP(w, r)
	})
}
