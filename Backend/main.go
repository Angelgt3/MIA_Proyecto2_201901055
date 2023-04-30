package main

import (
	"comandos"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	//routes
	http.HandleFunc("/ejecutar", ejecutar)

	//crea el servidor
	fmt.Println("Run server: http://localhost:3000")
	http.ListenAndServe("localhost:3000", nil)

}

type PETICION struct {
	Nombre string `json:"nombre"`
	Texto  string `json:"texto"`
}

func ejecutar(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	//leo el json que recibo
	contenido, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	//LEER UN JSON EN DONDE RECIBE EL LINK DONDE SE ENCUENTRA LA ENTRADA
	/*
		peti := PETICION{}
		err = json.Unmarshal(contenido, &peti)
		if err != nil {
			log.Fatal(err)
		}
			//Analizar el archivo
			var res string = ""
			if string(peti.Nombre) == "script" { //lee un script
				res = comandos.Abir_archivo(string(peti.Texto))
			}
	*/

	//Analizar el archivo comandos texto
	res := comandos.Leer_archivo(string(contenido))

	//Escribo el json para enviar
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["nombre"] = "201901055 - Angel"
	fmt.Println(res)
	resp["result"] = res
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
