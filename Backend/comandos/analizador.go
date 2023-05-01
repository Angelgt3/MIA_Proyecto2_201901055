package comandos

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// ------------------------------ANALIZADOR---------------------------------------------
func Abir_archivo(ruta string) string {
	bytesLeidos, err := ioutil.ReadFile(ruta)
	if err != nil {
		//fmt.Printf("Error leyendo archivo: %v", err)
		respuesta += "ERROR: NO SE PUDO ABRIR EL ARCHIVO (analizador)"
	}
	return string(bytesLeidos)
}

func Leer_archivo(bytesLeidos string) string {
	//abrimos el archivo
	fmt.Println("EJECUTANDO SCRIPT ... ")
	//convierto todo el texto en minusculas}
	contenido := strings.ToLower(string(bytesLeidos))
	//separa el contenido por lineas
	lineasComoArreglo := strings.Split(string(contenido), "\n")
	var parametros [15]string
	for _, nombre := range lineasComoArreglo {
		if nombre != "" {
			//separa comando y parametros
			cpara := strings.Split(string(nombre), ">")
			cont := 0
			for _, paaa := range cpara {
				r := strings.Replace(paaa, " ", "", -1) //quito los espacios en blanco
				parametros[cont] = r                    //el primero siempre sera el comando -seguido de los parametros
				cont++
			}
			ejecutar_script(parametros)
		}
	}
	tem := respuesta
	respuesta = ""
	fmt.Println("SE TERMINO DE EJECUTAR EL SCRIPT :)")
	return tem
}

func ejecutar_script(comando [15]string) {
	if comando[0] == "" { //no aceptar espacios vacios
		return
	}
	if string(comando[0][0]) == "#" { //quito los comentarios al inicio
		return
	}
	if comando[0] == "mkdisk" { //MKDISK-ANALIZADOR
		var size int = 0
		var path, fit, unit string = "", "ff", "m"
		for i := 1; i < 15; i++ {
			if comando[i] == "" {
				break
			}
			part := strings.Split(string(comando[i]), "=")
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
	} else if comando[0] == "login" { //LOGIN-ANALIZADOR
		var user, pass, id string = "", "", ""
		for i := 1; i < 15; i++ {
			if comando[i] == "" {
				break
			}
			part := strings.Split(string(comando[i]), "=")
			if part[0] == "user" {
				if user == "" {
					user = part[1]
				}
			} else if part[0] == "pwd" {
				if pass == "" {
					pass = part[1]
				}
			} else if part[0] == "id" {
				if id == "" {
					id = part[1]
				}
			}
		}
		Login(user, pass, id)
	} else if comando[0] == "logout" { //LOGOUT-ANALIZADOR
		Logout()
	} else if comando[0] == "mkgrp" { //MKGRP-ANALIZADOR
		var name string = ""
		for i := 1; i < 15; i++ {
			if comando[i] == "" {
				break
			}
			part := strings.Split(string(comando[i]), "=")
			if part[0] == "name" {
				if name == "" {
					name = part[1]
				}
			}
		}
		Mkgrp(name)
	} else if comando[0] == "rmgrp" { //RMGRP-ANALIZADOR
		var name string = ""
		for i := 1; i < 15; i++ {
			if comando[i] == "" {
				break
			}
			part := strings.Split(string(comando[i]), "=")
			if part[0] == "name" {
				if name == "" {
					name = part[1]
				}
			}
		}
		Rmgrp(name)
	} else if comando[0] == "mkuser" { //MKUSER-ANALIZADOR
		var user, pwd, grp string = "", "", ""
		for i := 1; i < 15; i++ {
			if comando[i] == "" {
				break
			}
			part := strings.Split(string(comando[i]), "=")
			if part[0] == "user" {
				if user == "" {
					user = part[1]
				}
			} else if part[0] == "pwd" {
				if pwd == "" {
					pwd = part[1]
				}
			} else if part[0] == "grp" {
				if grp == "" {
					grp = part[1]
				}
			}
		}
		Mkuser(user, pwd, grp)
	} else if comando[0] == "rmusr" { //RMUSR-ANALIZADOR
		var user string = ""
		for i := 1; i < 15; i++ {
			if comando[i] == "" {
				break
			}
			part := strings.Split(string(comando[i]), "=")
			if part[0] == "user" {
				if user == "" {
					user = part[1]
				}
			}
		}
		Rmusr(user)
	} else if comando[0] == "mkdir" { //MKDIR-ANALIZADOR
		var path string = ""
		var r bool = false
		for i := 1; i < 15; i++ {
			if comando[i] == "" {
				break
			}
			part := strings.Split(string(comando[i]), "=")
			if part[0] == "path" {
				if path == "" {
					path = part[1]
				}
			} else if part[0] == "r" {
				r = true
			}
		}
		Mkdir(path, r)
	} else if comando[0] == "mkfile" { //MKFILE-ANALIZADOR
		var path, size, cont string = "", "0", ""
		var r bool = false
		for i := 1; i < 15; i++ {
			if comando[i] == "" {
				break
			}
			part := strings.Split(string(comando[i]), "=")
			if part[0] == "path" {
				if path == "" {
					path = part[1]
				}
			} else if part[0] == "r" {
				r = true
			} else if part[0] == "size" {
				if size == "0" {
					size = part[1]
				}
			} else if part[0] == "cont" {
				if cont == "" {
					cont = part[1]
				}
			}
		}
		Mkfile(path, r, size, cont)
	} else if comando[0] == "pause" { //PAUSE-ANALIZADOR

		fmt.Println(" PAUSA: Presione cualguien tecla para continuar")
		pause := ""
		fmt.Scanln(&pause)

	}
}
