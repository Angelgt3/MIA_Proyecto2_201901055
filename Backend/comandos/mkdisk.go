package comandos

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// ------------------------------MKDISK---------------------------------------------
func mkdisk(path string, size int, fit string, unit string) {
	//fmt.Println("path:" + path)
	//fmt.Println("size:" + strconv.Itoa(size))
	//fmt.Println("unit:" + unit)
	//fmt.Println("fit:" + fit)
	//se crea el disco
	if size <= 0 {
		fmt.Println("ERROR: size:" + strconv.Itoa(size) + " invalido")
		return
	}
	var tam int
	//convierto el tam a bytes
	if unit == "k" {
		tam = size
	} else {
		tam = size * 1024
	}
	limite := 0
	bloque := make([]byte, 1024)
	//Preparacion del bloque a escribir en archivo
	for j := 0; j < 1024; j++ {
		bloque[j] = 0
	}
	//verificar si existe la ruta del archivo
	ruta := ""

	carpetas := strings.Split(path, "/")
	for _, carp := range carpetas {
		if carp == "" {
			continue
		}
		if !strings.HasSuffix(carp, ".dsk") {
			//os.Chmod(ruta, 0755)
			ruta = ruta + "/" + carp
			if _, err := os.Stat(ruta); os.IsNotExist(err) {
				err = os.Mkdir(ruta, 0755)
				if err != nil {
					panic(err)
				}
			}
		}
	}
	//Creacion del archivo
	disco, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	//Escritura de 0
	for limite < tam {
		_, err := disco.Write(bloque)
		if err != nil {
			panic(err)
		}
		limite++
	}
	//Cierre de archivo
	disco.Close()

	//Creacion del mbr
	tams := strconv.Itoa(tam * 1024)                // size
	tim := time.Now().Format("2006-01-02 15:04:05") // fecha de creacion
	ale := strconv.Itoa(rand.Intn(101))             // numero random de 0 a 100
	ft := string(fit)                               // tipo de ajuste
	mbr := newMBR(tams, tim, ale, ft)

	//Guarda el mbr
	disco, err = os.OpenFile(string(path), os.O_RDWR, 0660) // Apertura del archivo
	if err != nil {
		fmt.Println("ERROR: NO SE LOGRO ABRIR EL ARCHIVO")
	}

	disco.Seek(0, 0)
	binary.Write(disco, binary.BigEndian, mbr) // se convierte en arreglo de byte
	disco.Close()
	fmt.Println("SE CREO EL DISCO:  " + path)

}
