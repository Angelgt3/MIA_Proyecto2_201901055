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
		var path, fit, unit string = "", "ff", "m"
		for i := 1; i < 15; i++ {
			if comando[i] == "" {
				break
			}
			part := strings.Split(string(comando[i]), "=")
			//fmt.Println(part[0])
			//fmt.Println(part[1])
			if part[0] == "size" {
				size, _ = strconv.Atoi(part[1])
			} else if part[0] == "fit" {
				fit = part[1]
			} else if part[0] == "path" {
				if path == "" {
					path = part[1]
				}
			} else if part[0] == "unit" {
				unit = part[1]
			}

		}
		mkdisk(path, size, fit, unit)
	} else if comando[0] == "rmdisk" { //RMDISK-ANALIZADOR
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
		rmdisk(path)
	} else if comando[0] == "fdisk" { //FDISK-ANALIZADOR
		var path, unit, typee, fit, name string = "", "k", "p", "wf", ""
		var size int = 0
		for i := 1; i < 15; i++ {
			if comando[i] == "" {
				break
			}
			part := strings.Split(string(comando[i]), "=")
			if part[0] == "path" {
				if path == "" {
					path = part[1]
				}
			} else if part[0] == "size" {
				size, _ = strconv.Atoi(part[1])
			} else if part[0] == "unit" {
				unit = part[1]
			} else if part[0] == "type" {
				typee = part[1]
			} else if part[0] == "fit" {
				fit = part[1]
			} else if part[0] == "name" {
				if name == "" {
					name = part[1]
				}
			}
		}
		fdisk(size, unit, path, typee, fit, name)
	} else if comando[0] == "mount" { //MOUNT-ANALIZADOR
		var path, name, mostrar string = "", "", ""
		for i := 1; i < 15; i++ {
			if comando[i] == "" {
				break
			}
			part := strings.Split(string(comando[i]), "=")
			if part[0] == "path" {
				if path == "" {
					path = part[1]
				}
			} else if part[0] == "name" {
				if name == "" {
					name = part[1]
				}
			} else if part[0] == "mostrar" {
				if name == "" {
					mostrar = "!"
				}
			}

		}
		if mostrar == "!" {
			Mostrar_mount()
		} else {
			Mount(path, name)
		}

	} else if comando[0] == "rep" { //REP-ANALIZADOR
		var path, name, id, ruta string = "", "", "", ""
		for i := 1; i < 15; i++ {
			if comando[i] == "" {
				break
			}
			part := strings.Split(string(comando[i]), "=")
			if part[0] == "path" {
				if path == "" {
					path = part[1]
				}
			} else if part[0] == "name" {
				if name == "" {
					name = part[1]
				}
			} else if part[0] == "id" {
				id = part[1]
			} else if part[0] == "ruta" {
				ruta = part[1]
			}
		}
		Crear_reporte(name, path, id, ruta)
	} else if comando[0] == "mkfs" { //MKFS-ANALIZADOR
		var id, typee string = "", ""
		for i := 1; i < 15; i++ {
			if comando[i] == "" {
				break
			}
			part := strings.Split(string(comando[i]), "=")
			if part[0] == "id" {
				if id == "" {
					id = part[1]
				}
			} else if part[0] == "type" {
				if typee == "" {
					typee = part[1]
				}
			}
		}
		mkfs(id, typee)
	}
}
