package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	//Analizar el archivo
	nombreArchivo := "/home/angel/Escritorio/MIA/Semestre3/[MIA]Proyecto2_201901055/MIA_Proyecto2_201901055/Backend/entrada.eea"
	leer_archivo(nombreArchivo)
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

// ------------------------------ANALIZADOR---------------------------------------------
func leer_archivo(nombreArchivo string) {
	//abrimos el archivo
	bytesLeidos, err := ioutil.ReadFile(nombreArchivo)
	if err != nil {
		fmt.Printf("Error leyendo archivo: %v", err)
	}
	//convierto todo el texto en minusculas}
	contenido := strings.ToLower(string(bytesLeidos))
	//separa el contenido por lineas
	lineasComoArreglo := strings.Split(string(contenido), "\n")

	var parametros [15]string
	for _, nombre := range lineasComoArreglo {
		if nombre != "" {
			//fmt.Println(nombre)
			//separa comando y parametros
			cpara := strings.Split(string(nombre), ">")
			cont := 0
			for _, paaa := range cpara {
				r := strings.Replace(paaa, " ", "", -1) //quito los espacios en blanco
				//fmt.Println(r)
				parametros[cont] = r //el primero siempre sera el comando -seguido de los parametros
				cont++
			}
			ejecutar_script(parametros)
		}

	}
}

func ejecutar_script(comando [15]string) {
	if comando[0] == "" { //no aceptar espacios vacios
		return
	}
	if string(comando[0][0]) == "#" { //quito los comentarios al inicio
		return
	}
	fmt.Println("comando:" + comando[0])
	if comando[0] == "mkdisk" { //MKDISK-ANALIZADOR
		var size int = 0
		for i := 1; i < 15; i++ {
			if comando[i] == "" {
				break
			}
			part := strings.Split(string(comando[i]), "=")
			//fmt.Println(part[0])
			//fmt.Println(part[1])
			if part[0] == "size" {
				size, _ = strconv.Atoi(part[1])
			}

		}
		fmt.Println(size)
	}
}
