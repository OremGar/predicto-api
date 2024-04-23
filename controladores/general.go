package controladores

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/OremGar/predicto-api/respuestas"
)

func ResetearBD(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("/bin/sh", "-c", "echo Aut201104 | sudo -S systemctl restart postgresql")
	err := cmd.Run()

	log.Print(err)

	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error when running script: %v", err))
		return
	}

	respuestas.JsonResponse(w, http.StatusOK, nil, 0, nil)
}
