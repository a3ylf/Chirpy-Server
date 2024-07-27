package handlers

import (
	"fmt"
	"net/http"
)

func (cfg *Apiconfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverhits++
        next.ServeHTTP(w,r)
    })
}

func (cfg *Apiconfig)HandlerMetrics(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte(fmt.Sprintf(`
<html>

<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>

</html>
	`,cfg.fileserverhits)))
}
