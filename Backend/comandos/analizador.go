package comandos

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// ------------------------------ANALIZADOR---------------------------------------------
func Leer_archivo(nombreArchivo string) {
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
	//fmt.Println("comando:" + comando[0])
	if comando[0] == "mkdisk" { //MKDISK-ANALIZADOR
		var size int = 0
		var path string = ""
		var fit string = "ff"
		var unit string = "m"
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
			if part[0] == "fit" {
				fit = part[1]
			}
			if part[0] == "path" {
				if path == "" {
					path = part[1]
				}
			}
			if part[0] == "unit" {
				unit = part[1]
			}

		}
		mkdisk(path, size, fit, unit)
	}
	if comando[0] == "rmdisk" { //RMDISK-ANALIZADOR
		var path string = ""
		for i := 1; i < 15; i++ {
			if comando[i] == "" {
				break
			}
			part := strings.Split(string(comando[i]), "=")
			if part[0] == "path" {
				if path == "" {
					path = part[1]
				}
			}
		}
		fmt.Println(path)
		rmdisk(path)
	}
}
