package comandos

import "fmt"

func Mkdir(path string, r bool) {
	if path == "" {
		fmt.Println("ERROR MKDIR: FALTA PARAMETRO DE PATH")
		return
	}

	pathSep := separar_ruta(path)
	NameCarpetaNew := pathSep[len(pathSep)-1]
	pathSep = append(pathSep[:len(pathSep)-1])
	pathCarpetaPadre := unir_ruta(pathSep)
	if !usuario_activo.Logeado {
		fmt.Println("ERRROR: NO HAY USUARIO LOGEADO")
		return
	}
	indiceInodoCarpetaPadre := index_inodo_ruta(pathCarpetaPadre, usuario_activo.Montada, 0)
	if indiceInodoCarpetaPadre == -1 {
		if r {
			indiceInodoCarpetaPadre = crear_ruta(pathCarpetaPadre, usuario_activo.Montada)
		} else {
			fmt.Println("ERROR: NO EXISTE LA RUTA INDICADA")
			return
		}
	}
	crear_carpeta(indiceInodoCarpetaPadre, NameCarpetaNew, usuario_activo.Montada)

	fmt.Println("SE CREO LA CARPETA CORRECTAMENTE")

}
