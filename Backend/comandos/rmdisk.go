package comandos

import (
	"fmt"
	"os"
)

func rmdisk(path string) {
	//verificar si el archivo existe
	//fmt.Println(path)
	if !archivoExiste(path) {
		fmt.Println("ERROR: NO EXISTE LA RUTA")
		return
	}
	// confirmacion de eliminar
	var input string
	fmt.Println("DESEA ELIMINAR EL DISCO: " + path + "? (S/N)")
	fmt.Scanln(&input)
	if input == "S" || input == "s" {
		err := os.Remove(path) // elimina el archivo
		if err != nil {
			fmt.Printf("ERROR AL ELIMINAR EL ARCHIVO %v\n", err)
			return
		}
	}
}

func archivoExiste(ruta string) bool {
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}
