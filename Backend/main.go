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
	//Analizar el archivo
	nombreArchivo := "/home/angel/Escritorio/MIA/Semestre3/[MIA]Proyecto2_201901055/MIA_Proyecto2_201901055/Backend/entrada.eea"
	comandos.Leer_archivo(nombreArchivo)
	/*
		//routes
		http.HandleFunc("/ejecutar", ejecutar)

		//crea el servidor
		fmt.Println("El servidor esta correindo en el puerto 3000")
		fmt.Println("Run server: http://localhost:3000")
		http.ListenAndServe("localhost:3000", nil)
	*/
}

type CONTENIDO struct {
	texto     string
	contenido string
}

var cont CONTENIDO

func ejecutar(w http.ResponseWriter, r *http.Request) {

	//leo el json que recibo
	contenido, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	cont := contenido
	fmt.Println(string(cont))

	//Analizar el archivo
	//analizar()

	//Escribo el json para enviar
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["carnet"] = "201901055"
	resp["nombre"] = "Angel Geovany Aragón Pérez"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}
