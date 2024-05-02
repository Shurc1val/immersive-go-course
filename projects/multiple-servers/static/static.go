package static

import (
	"net/http"
	"strconv"
)

func Run(path string, port int) {
	http.Handle("/*", http.FileServer(http.Dir("./"+path)))
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
