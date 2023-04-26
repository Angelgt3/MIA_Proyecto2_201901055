package comandos

import (
	"bufio"
	"os"
	"strconv"
)

func Mkfile(path string, r bool, siz string, cont string) {
	if path == "" {
		//fmt.Println("ERROR MKFILE: NO SE INGRESO LOS PARAMETROS OBLIGATORIOS")
		respuesta += "\nERROR MKFILE: NO SE INGRESO LOS PARAMETROS OBLIGATORIOS"
		return
	}
	var contenido_archivo string = ""
	size, _ := strconv.Atoi(siz)
	if cont == "" {
		if size != 0 && size > 0 {
			// llenar el contenido de 0-9 por size
			for len(contenido_archivo) < size {
				for i := 0; i <= 9; i++ {
					if len(contenido_archivo) >= size {
						break
					}
					contenido_archivo += strconv.Itoa(i)
				}
			}
		}

	} else if cont != "" {
		if _, err := os.Stat(cont); os.IsNotExist(err) {
			//fmt.Println("ERROR MKFILE: NO EXISTE RUTA CONT")
			respuesta += "\nERROR MKFILE: NO EXISTE RUTA CONT"
			return
		}
		archivo, err := os.OpenFile(cont, os.O_RDWR, 0660)
		if err != nil {
			//fmt.Println("ERROR MKFILE: NO SE LOGRO ABRIR EL ARCHIVO")
			respuesta += "\nERROR MKFILE: NO SE LOGRO ABRIR EL ARCHIVO"
			return
		}
		fileScanner := bufio.NewScanner(archivo)
		fileScanner.Split(bufio.ScanLines)
		var lines []string
		for fileScanner.Scan() {
			lines = append(lines, fileScanner.Text())
		}
		archivo.Close()
		for _, line := range lines {
			contenido_archivo += line
			contenido_archivo += "\n"
		}
		contenido_archivo = contenido_archivo[0 : len(contenido_archivo)-1]
	}
	folderPath := directorio(path)
	fileName := nombre_archivo(path) + "." + extension(path)
	// si existe el archivo
	if existe_ruta(path, usuario_activo.Montada, 0) != -1 {
		//fmt.Println("MKFILE: SE SOBREESCRIBIRA EL ARCHIVO")
		respuesta += "\nMKFILE: SE SOBREESCRIBIRA EL ARCHIVO"
	}
	//si existe la ruta
	IndiceInodoFolder := existe_ruta(folderPath, usuario_activo.Montada, 0)
	if IndiceInodoFolder == -1 {
		if r {
			IndiceInodoFolder = crear_ruta(folderPath, usuario_activo.Montada)
		} else {
			//fmt.Println("ERROR MKFILE: NO EXISTE LA RUTA INDICADA PARA CREAR EL ARCHIVO")
			respuesta += "\nERROR MKFILE: NO EXISTE LA RUTA INDICADA PARA CREAR EL ARCHIVO"
			return
		}
	}
	IndiceInodoArchivo := crear_archivo(size, fileName, IndiceInodoFolder, usuario_activo.Montada)
	if IndiceInodoFolder == -1 {
		//fmt.Println("ERROR MKFILE: NO SE PUDO CREAR EL ARCHIVO EN LA CARPETA DESEADA")
		respuesta += "\nERROR MKFILE: NO SE PUDO CREAR EL ARCHIVO EN LA CARPETA DESEADA"
		return
	}
	modificar_archivo(IndiceInodoArchivo, contenido_archivo, usuario_activo.Montada)
	//fmt.Println("SE CREO EL ARCHIVO CORRECTAMENTE")
	respuesta += "\nSE CREO EL ARCHIVO CORRECTAMENTE"
}
