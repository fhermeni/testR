package rest

import "net/http"

func index(c Context) error {
	http.Redirect(c.w, c.r, "/assets/root.html", http.StatusMovedPermanently)
	return nil
}
