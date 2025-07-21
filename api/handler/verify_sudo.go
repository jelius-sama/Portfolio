package handler

import (
	"KazuFolio/util"
	"net/http"
)

func VerifySudo(w http.ResponseWriter, r *http.Request) {
	util.VerifySudo(w, r)

	w.WriteHeader(http.StatusNoContent)
}
