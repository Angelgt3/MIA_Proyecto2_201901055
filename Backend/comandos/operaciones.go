package comandos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// ------------------------------VARIABLES---------------------------------------------
var Tmontadas List
var usuario_activo USUARIO

// ------------------------------FUNCIONES---------------------------------------------

func newAjuste() Ajuste {
	var l [6]int
	e := Ajuste{l, l, l, l}
	return e
}

func newEbr(st, f, s, tm, tms, n string) EBR {
	e := EBR{}
	copy(e.Part_status[:], st)
	copy(e.Part_fit[:], f)
	copy(e.Part_start[:], s)
	copy(e.Part_size[:], tm)
	copy(e.Part_next[:], tms)
	copy(e.Part_name[:], n)
	return e
}

func newMBR(t, f, d, ft string) MBR {
	e := MBR{}
	copy(e.Mbr_tamano[:], t)
	copy(e.Mbr_fecha_creacion[:], f)
	copy(e.Mbr_dsk_signature[:], d)
	copy(e.Dsk_fit[:], ft)
	return e
}

func newTlist() Disco {
	e := Disco{}
	return e
}

func newInodo() TINODOS {
	var inodo TINODOS
	for i := 0; i < 64; i = i + 4 {
		copy(inodo.I_block[i:i+3], "-")
	}
	return inodo
}

func newCarpeta() BLOQUE_CARPETA {
	var carpeta BLOQUE_CARPETA
	for i := 0; i < 4; i++ {
		copy(carpeta.B_content[i].B_inodo[0:3], "-")
	}
	return carpeta
}

// retorna la posicion para un nuevo bloque
func new_index_bloque(disco Disco) int {
	bm := get_bitmap(disco, false)
	for i := 0; i < len(bm); i++ {
		if bm[i] == '0' {
			return i
		}
	}
	return -1
}

func new_index_inodo(disco Disco) int {
	bm := get_bitmap(disco, true)
	for i := 0; i < len(bm); i++ {
		if bm[i] == '0' {
			return i
		}
	}
	return -1
}

func separar_ruta(path string) []string {
	rutas := strings.Split(path, "/")
	rutas = append(rutas[1:])
	return rutas
}

func unir_ruta(ruta []string) string {
	completo := ""
	for _, nombre := range ruta {
		completo += "/" + nombre
	}
	return completo
}

// retorna si existe la particion montada y su estruct
func existe_montada(id string, part *Disco) bool {
	q := Tmontadas
	for i := 0; i < q.Len(); i++ {
		t := q.GetValue(i).(Disco)
		if t.Id == id {
			*part = t
			return true
		}
	}
	return false
}

// retorna todos los usuarios del users.txt
func get_registros(usuarios string) [][]string {
	linea := strings.Split(string(usuarios), "\n")
	todos := make([][]string, len(linea))
	for i := 0; i < len(linea); i++ {
		usu := strings.Split(linea[i], ",")
		todos[i] = make([]string, len(usu))
		for j := 0; j < len(usu); j++ {
			usu[j] = strings.TrimSpace(usu[j])
			todos[i][j] = usu[j]
		}
	}
	return todos
}

// retorna el indice de las carpetas
func Indices_BC(inodo TINODOS, disco Disco) []int {
	apuntadores := make([]int, 0)
	if inodo.I_type[0] == '0' {
		for i := 0; i < 64; i = i + 4 {
			inosta := strings.Split(string(inodo.I_block[i:i+3]), "\x00")
			if inosta[0] != "-" {
				apunt, _ := strconv.Atoi(inosta[0])
				apuntadores = append(apuntadores, apunt)
			}
		}
	}
	return apuntadores
}

// escribe inodo en el disco por medio de su index
func escribir_inodo(index int, inodo TINODOS, particion Disco) {
	archivo, err := os.OpenFile(particion.Path, os.O_RDWR, 0660) // Apertura del archivo
	if err != nil {
		fmt.Println("ERROR: NO SE LOGRO ABRIR EL ARCHIVO (escribir_inood)")
	}
	Pstart := strings.Split(string(particion.Part.Part_start[:]), "\x00")
	Ps, _ := strconv.Atoi(Pstart[0])
	archivo.Seek(int64(Ps), 0)

	//modifico el super bloque
	var sb SUPER_BLOQUE
	rre := binary.Read(archivo, binary.BigEndian, &sb)
	if rre != nil {
		print(rre)
	}
	ind := index + 1
	ino := strconv.Itoa(ind)
	copy(sb.S_first_ino[:], ino)
	//lo escribo en el archivo
	archivo.Seek(int64(Ps), 0)
	binary.Write(archivo, binary.BigEndian, sb)

	//mofico el bitmap de inodos
	inodoCon := strings.Split(string(sb.S_inodes_count[:]), "\x00")
	Ic, _ := strconv.Atoi(inodoCon[0])
	bitmap_inodos := make([]byte, int(Ic))
	bmStart := strings.Split(string(sb.S_bm_inode_start[:]), "\x00")
	bm, _ := strconv.Atoi(bmStart[0])
	archivo.Seek(int64(bm), 0)
	rre = binary.Read(archivo, binary.BigEndian, bitmap_inodos)
	if rre != nil {
		print(rre)
	}
	copy(bitmap_inodos[index:index+1], "1")
	archivo.Seek(int64(bm), 0)
	binary.Write(archivo, binary.BigEndian, bitmap_inodos)
	//escribo el inodo en el archivo
	inosta := strings.Split(string(sb.S_inode_start[:]), "\x00")
	is, _ := strconv.Atoi(inosta[0])
	archivo.Seek(int64(is), 0)
	var TI TINODOS
	lenTi := new(bytes.Buffer)
	binary.Write(lenTi, binary.BigEndian, TI)
	archivo.Seek(int64(index*len(lenTi.Bytes())), 1)
	binary.Write(archivo, binary.BigEndian, inodo)
	archivo.Close()
}

// escribo el bitmap en el archivo disco
func escribir_bitmap(bitmap []byte, disco Disco, tipo bool) { //TRUE = INODO | FALSE = BLOQUE
	var sb SUPER_BLOQUE
	size := -1
	archivo, err := os.OpenFile(disco.Path, os.O_RDWR, 0660) // Apertura del archivo
	if err != nil {
		fmt.Println("ERROR: NO SE LOGRO ABRIR EL ARCHIVO (escribir_bitmap)")
	}
	Pstart := strings.Split(string(disco.Part.Part_start[:]), "\x00")
	Ps, _ := strconv.Atoi(Pstart[0])
	archivo.Seek(int64(Ps), 0) // Posicion inicial

	rre := binary.Read(archivo, binary.BigEndian, &sb) // se convierte en arreglo de byte
	if rre != nil {
		print(rre)
	}
	if tipo { // Inodo
		inosta := strings.Split(string(sb.S_bm_inode_start[:]), "\x00")
		is, _ := strconv.Atoi(inosta[0])
		archivo.Seek(int64(is), 0)
		bl := strings.Split(string(sb.S_bm_block_start[:]), "\x00")
		bls, _ := strconv.Atoi(bl[0])
		ino := strings.Split(string(sb.S_bm_inode_start[:]), "\x00")
		nos, _ := strconv.Atoi(ino[0])
		size = bls - nos
	} else { // Bloque
		inosta := strings.Split(string(sb.S_bm_block_start[:]), "\x00")
		is, _ := strconv.Atoi(inosta[0])
		archivo.Seek(int64(is), 0)
		bl := strings.Split(string(sb.S_inode_start[:]), "\x00")
		bls, _ := strconv.Atoi(bl[0])
		ino := strings.Split(string(sb.S_bm_block_start[:]), "\x00")
		nos, _ := strconv.Atoi(ino[0])
		size = bls - nos
	}

	bm := make([]byte, size)
	copy(bm[:], bitmap)
	binary.Write(archivo, binary.BigEndian, bm) // fuardo el sb
	archivo.Close()
}

// escribe el bloque archivo en el archivo disco
func escribir_BA(bloque BLOQUE_ARCHIVO, index int, disco Disco) {
	archivo, err := os.OpenFile(disco.Path, os.O_RDWR, 0660)
	if err != nil {
		fmt.Println("ERROR: NO SE LOGRO ABRIR EL ARCHIVO")
	}
	Pstart := strings.Split(string(disco.Part.Part_start[:]), "\x00")
	Ps, _ := strconv.Atoi(Pstart[0])
	archivo.Seek(int64(Ps), 0)
	var sb SUPER_BLOQUE
	rre := binary.Read(archivo, binary.BigEndian, &sb)
	if rre != nil {
		print(rre)
	}
	copy(sb.S_first_blo[:], strconv.Itoa(index+1))
	archivo.Seek(int64(Ps), 0)
	binary.Write(archivo, binary.BigEndian, sb)

	inost := strings.Split(string(sb.S_inode_start[:]), "\x00")
	Inos, _ := strconv.Atoi(inost[0])
	blmst := strings.Split(string(sb.S_bm_block_start[:]), "\x00")
	blo, _ := strconv.Atoi(blmst[0])

	bitmap := make([]byte, Inos-blo)
	archivo.Seek(int64(blo), 0)
	rre = binary.Read(archivo, binary.BigEndian, &bitmap)
	if rre != nil {
		print(rre)
	}
	bitmap[index] = '1'
	archivo.Seek(int64(blo), 0)
	binary.Write(archivo, binary.BigEndian, bitmap)

	blst := strings.Split(string(sb.S_block_start[:]), "\x00")
	bloq, _ := strconv.Atoi(blst[0])
	archivo.Seek(int64(bloq), 0)
	var bar BLOQUE_ARCHIVO
	lenBar := new(bytes.Buffer)
	binary.Write(lenBar, binary.BigEndian, bar)
	archivo.Seek(int64(index*len(lenBar.Bytes())), 1)
	binary.Write(archivo, binary.BigEndian, bloque)
	archivo.Close()
}

// escribe el bloque carpeta en el archivo disco
func escribir_BC(bloque BLOQUE_CARPETA, index int, disco Disco) {
	archivo, err := os.OpenFile(disco.Path, os.O_RDWR, 0660)
	if err != nil {
		fmt.Println("ERROR: NO SE LOGRO ABRIR EL ARCHIVO")
	}
	Pstart := strings.Split(string(disco.Part.Part_start[:]), "\x00")
	Ps, _ := strconv.Atoi(Pstart[0])
	archivo.Seek(int64(Ps), 0)
	var sb SUPER_BLOQUE
	rre := binary.Read(archivo, binary.BigEndian, &sb)
	if rre != nil {
		print(rre)
	}
	copy(sb.S_first_blo[:], strconv.Itoa(index+1))
	archivo.Seek(int64(Ps), 0)
	binary.Write(archivo, binary.BigEndian, sb)

	inost := strings.Split(string(sb.S_inode_start[:]), "\x00")
	Inos, _ := strconv.Atoi(inost[0])

	blmst := strings.Split(string(sb.S_bm_block_start[:]), "\x00")
	blo, _ := strconv.Atoi(blmst[0])

	bitmap := make([]byte, Inos-blo)
	archivo.Seek(int64(blo), 0)
	rre = binary.Read(archivo, binary.BigEndian, &bitmap) // se convierte en arreglo de byte
	if rre != nil {
		print(rre)
	}
	bitmap[index] = '1'
	archivo.Seek(int64(blo), 0)
	binary.Write(archivo, binary.BigEndian, bitmap)

	blst := strings.Split(string(sb.S_block_start[:]), "\x00")
	bloq, _ := strconv.Atoi(blst[0])
	archivo.Seek(int64(bloq), 0)

	apunt := index * 64
	archivo.Seek(int64(apunt), 1)
	binary.Write(archivo, binary.BigEndian, bloque)
	archivo.Close()
}

// retorna el inodo deseado por medio de un index
func get_inodo(index int, disco Disco) TINODOS {
	var sb SUPER_BLOQUE
	var inodo TINODOS

	archivo, err := os.OpenFile(disco.Path, os.O_RDWR, 0660)
	if err != nil {
		fmt.Println("ERROR: NO SE LOGRO ABRIR EL ARCHIVO (get_inodo)")
	}
	Pstart := strings.Split(string(disco.Part.Part_start[:]), "\x00")
	Ps, _ := strconv.Atoi(Pstart[0])
	archivo.Seek(int64(Ps), 0)

	rre := binary.Read(archivo, binary.BigEndian, &sb)
	if rre != nil {
		print(rre)
	}
	//se obtiene el inodo deseado
	inosta := strings.Split(string(sb.S_inode_start[:]), "\x00")
	is, _ := strconv.Atoi(inosta[0])
	archivo.Seek(int64(is), 0)
	lenTi := new(bytes.Buffer)
	binary.Write(lenTi, binary.BigEndian, inodo)
	archivo.Seek(int64(index*len(lenTi.Bytes())), 1)
	binary.Read(archivo, binary.BigEndian, &inodo)
	archivo.Close()
	return inodo
}

// retorna el bloque archivo deseado por medio de un index
func get_bloque_archivo(index int, disco Disco) BLOQUE_ARCHIVO {
	var barchivo BLOQUE_ARCHIVO
	var sb SUPER_BLOQUE
	puntero := 64 * index
	archivo, err := os.OpenFile(disco.Path, os.O_RDWR, 0660) // Apertura del archivo
	if err != nil {
		fmt.Println("ERROR: NO SE LOGRO ABRIR EL ARCHIVO (get_bloque_archivo)")
	}
	Pstart := strings.Split(string(disco.Part.Part_start[:]), "\x00")
	Ps, _ := strconv.Atoi(Pstart[0])
	archivo.Seek(int64(Ps), 0) // Posicion inicial

	rre := binary.Read(archivo, binary.BigEndian, &sb) // se convierte en arreglo de byte
	if rre != nil {
		print(rre)
	}
	// SE OBTIENE EL BLOQUE DEL INDICE
	inosta := strings.Split(string(sb.S_block_start[:]), "\x00")
	is, _ := strconv.Atoi(inosta[0])
	archivo.Seek(int64(is), 0)
	archivo.Seek(int64(puntero), 1)
	binary.Read(archivo, binary.BigEndian, &barchivo)
	archivo.Close()
	return barchivo
}

// retorna el bloque carpeta deseado por medio de un index
func get_bloque_carpeta(index int, disco Disco) BLOQUE_CARPETA {
	var sb SUPER_BLOQUE
	puntero := 64 * index
	archivo, err := os.OpenFile(disco.Path, os.O_RDWR, 0660)
	if err != nil {
		fmt.Println("ERROR: NO SE LOGRO ABRIR EL ARCHIVO")
	}
	Pstart := strings.Split(string(disco.Part.Part_start[:]), "\x00")
	Ps, _ := strconv.Atoi(Pstart[0])
	archivo.Seek(int64(Ps), 0)
	rre := binary.Read(archivo, binary.BigEndian, &sb)
	if rre != nil {
		print(rre)
	}
	var bc BLOQUE_CARPETA
	inosta := strings.Split(string(sb.S_block_start[:]), "\x00")
	is, _ := strconv.Atoi(inosta[0])
	archivo.Seek(int64(is), 0)
	archivo.Seek(int64(puntero), 1)
	binary.Read(archivo, binary.BigEndian, &bc)
	archivo.Close()
	return bc
}

// retorna el BITMAP
func get_bitmap(disco Disco, tipo bool) []byte { //TRUE = INODO | FALSE = BLOQUE
	var sb SUPER_BLOQUE
	size := -1
	archivo, err := os.OpenFile(disco.Path, os.O_RDWR, 0660)
	if err != nil {
		fmt.Println("ERROR: NO SE LOGRO ABRIR EL ARCHIVO (get_bitmap)")
	}
	Pstart := strings.Split(string(disco.Part.Part_start[:]), "\x00")
	Ps, _ := strconv.Atoi(Pstart[0])
	archivo.Seek(int64(Ps), 0)
	rre := binary.Read(archivo, binary.BigEndian, &sb)
	if rre != nil {
		print(rre)
	}
	if tipo { // Inodo
		inosta := strings.Split(string(sb.S_bm_inode_start[:]), "\x00")
		is, _ := strconv.Atoi(inosta[0])
		archivo.Seek(int64(is), 0)
		bl := strings.Split(string(sb.S_bm_block_start[:]), "\x00")
		bls, _ := strconv.Atoi(bl[0])
		ino := strings.Split(string(sb.S_bm_inode_start[:]), "\x00")
		nos, _ := strconv.Atoi(ino[0])
		size = bls - nos
	} else { // Bloque
		inosta := strings.Split(string(sb.S_bm_block_start[:]), "\x00")
		is, _ := strconv.Atoi(inosta[0])
		archivo.Seek(int64(is), 0)
		bl := strings.Split(string(sb.S_inode_start[:]), "\x00")
		bls, _ := strconv.Atoi(bl[0])
		ino := strings.Split(string(sb.S_bm_block_start[:]), "\x00")
		nos, _ := strconv.Atoi(ino[0])
		size = bls - nos
	}
	bm := make([]byte, size+1)
	rre = binary.Read(archivo, binary.BigEndian, &bm)
	if rre != nil {
		print(rre)
	}
	archivo.Close()
	return bm
}

// retorna el contenido del bloque de contenido de un archivo por su index
func leer_archivo(index int, disco Disco) string {
	inodo := get_inodo(index, disco)
	contenido := ""
	for i := 0; i < 64; i = i + 4 {
		inosta := strings.Split(string(inodo.I_block[i]), "\x00")
		if inosta[0] != "-" {
			is, _ := strconv.Atoi(inosta[0])
			bloqueU := get_bloque_archivo(is, disco)
			acont := strings.Split(string(bloqueU.B_content[:]), "\x00")
			contenido += acont[0]
		}

	}
	return contenido
}

// elimina todos los bloques del archivo
func eliminar_bloques_archivo(index int, disco Disco) {
	inodo := get_inodo(index, disco)
	for i := 0; i < 64; i = i + 4 {
		inosta := strings.Split(string(inodo.I_block[i:i+3]), "\x00")
		if inosta[0] == "-" {
			break
		}

		bloques := get_bitmap(disco, false)
		is, _ := strconv.Atoi(inosta[0])
		copy(bloques[is:is+3], "0")
		escribir_bitmap(bloques, disco, false)
		copy(inodo.I_block[i:i+3], "-")
	}
	escribir_inodo(index, inodo, disco)
}

// crea los bloques de un archivo
func escribir_bloques_archivo(index int, texto string, disco Disco) {
	inodo := get_inodo(index, disco)
	btes := 0
	lenn := len(texto)
	for i := 0; i < 64; i = i + 4 {
		if btes >= lenn {
			break
		}
		inosta := strings.Split(string(inodo.I_block[i:i+3]), "\x00")
		if inosta[0] == "-" {
			var ba BLOQUE_ARCHIVO
			pos_ba := new_index_bloque(disco)
			if len(texto) > 63 {
				copy(ba.B_content[:], texto[0:63])
				texto = texto[63:]
			} else {
				copy(ba.B_content[:], texto[0:])
			}

			btes = btes + 63
			sf := strconv.Itoa(pos_ba)
			copy(inodo.I_block[i:i+3], sf)
			escribir_BA(ba, pos_ba, disco)
		}
	}
	escribir_inodo(index, inodo, disco)
}

// retorna la posicion del inodo de un path
func index_inodo_ruta(path string, disco Disco, index int) int {
	if path == "" {
		return 0
	}
	var sb SUPER_BLOQUE
	ruta := separar_ruta(path)

	archivo, err := os.OpenFile(disco.Path, os.O_RDWR, 0660) // Apertura del archivo
	if err != nil {
		fmt.Println("ERROR: NO SE LOGRO ABRIR EL ARCHIVO")
	}
	Pstart := strings.Split(string(disco.Part.Part_start[:]), "\x00")
	Ps, _ := strconv.Atoi(Pstart[0])
	archivo.Seek(int64(Ps), 0) // Posicion inicial

	rre := binary.Read(archivo, binary.BigEndian, &sb) // se convierte en arreglo de byte
	if rre != nil {
		print(rre)
	}
	// EMPIEZA A OBTNEER LOS APUNTADORES A CARPETAS
	inodo := get_inodo(index, disco)
	iBcarpeta := Indices_BC(inodo, disco)
	// RECORRER TODOS LOS BLOQUES
	for _, bloque := range iBcarpeta {
		var bc BLOQUE_CARPETA
		bls := strings.Split(string(sb.S_block_start[:]), "\x00")
		blsi, _ := strconv.Atoi(bls[0])
		archivo.Seek(int64(blsi), 0)                      // Posicion inicial
		archivo.Seek(int64(bloque*64), 1)                 // Posicion inicial
		rre = binary.Read(archivo, binary.BigEndian, &bc) // se convierte en arreglo de byte
		if rre != nil {
			print(rre)
		}
		// VERIFICO EN CADA APUNTADOR DEL BLOQUE
		for i := 0; i < 4; i++ {
			name := strings.Split(string(bc.B_content[i].B_name[:]), "\x00")
			if name[0] == ruta[0] {
				rutahijo := ruta
				rutahijo = append(rutahijo[1:])
				ret := strings.Split(string(bc.B_content[i].B_inodo[:]), "\x00")
				retn, _ := strconv.Atoi(ret[0])
				// CUANDO SE QUEDA VACIO SE LLEGO AL FINAL
				if len(rutahijo) == 0 {
					archivo.Close()
					return retn
				}
				// SINO SE USA EL PATH HIJO
				pp := unir_ruta(rutahijo)
				// Y SE SIGUE BUSCANDO DE FORMA RECURSIVA CON EL PATH HIIJO
				resBusqueda := index_inodo_ruta(pp, disco, retn)
				if resBusqueda != -1 {
					return resBusqueda
				}
			}
		}

	}

	return -1
}

// crea la ruta en el sistema de archivo
func crear_ruta(ruta string, disco Disco) int {
	carpetas := separar_ruta(ruta)
	completo := 0
	for i := 0; i < len(carpetas); i++ {
		raiz := make([]string, 0)
		hijo := make([]string, 0)
		raiz = append(carpetas[0:i])
		hijo = append(carpetas[0 : i+1])

		IndiceR := index_inodo_ruta(unir_ruta(raiz), disco, 0)
		if IndiceR == -1 {
			return -1
		}
		IndiceH := index_inodo_ruta(unir_ruta(hijo), disco, 0)
		if IndiceH != -1 {
			completo = IndiceH

		} else {
			completo = crear_carpeta(IndiceR, carpetas[i], disco)
		}
	}
	return completo
}

// crea una carpeta
func crear_carpeta(index int, carpeta string, l Disco) int {
	raiz := get_inodo(index, l)
	carpetaNew := new_index_inodo(l)
	Bcarpeta := new_index_bloque(l)
	IcarpetaNew := newInodo()
	IcarpetaNew.I_type[0] = '0'
	copy(IcarpetaNew.I_block[0:3], strconv.Itoa(Bcarpeta))
	copy(IcarpetaNew.I_uid[:], "1")
	copy(IcarpetaNew.I_gid[:], "1")
	copy(IcarpetaNew.I_size[:], "0")
	copy(IcarpetaNew.I_atime[:], time.Now().Format("2006-01-02 15:04:05"))
	copy(IcarpetaNew.I_ctime[:], time.Now().Format("2006-01-02 15:04:05"))
	copy(IcarpetaNew.I_mtime[:], time.Now().Format("2006-01-02 15:04:05"))
	copy(IcarpetaNew.I_perm[:], "664")
	bloqueCarpetaNueva := newCarpeta()
	copy(bloqueCarpetaNueva.B_content[0].B_name[:], ".")
	copy(bloqueCarpetaNueva.B_content[0].B_inodo[:], strconv.Itoa(carpetaNew))
	copy(bloqueCarpetaNueva.B_content[1].B_name[:], "..")
	copy(bloqueCarpetaNueva.B_content[1].B_inodo[:], strconv.Itoa(index))
	escribir_inodo(carpetaNew, IcarpetaNew, l)
	escribir_BC(bloqueCarpetaNueva, Bcarpeta, l)
	var bcontent CONTENIDO
	copy(bcontent.B_name[:], carpeta)
	copy(bcontent.B_inodo[:], strconv.Itoa(carpetaNew))
	listo := false
	for i := 0; i < 64; i = i + 4 {
		inosta := strings.Split(string(raiz.I_block[i:i+3]), "\x00")
		if inosta[0] != "-" {
			apunt, _ := strconv.Atoi(inosta[0])
			bloqueDeCarpetas := get_bloque_carpeta(apunt, l)
			for j := 0; j < 4; j++ {
				blosta := strings.Split(string(bloqueDeCarpetas.B_content[j].B_inodo[:]), "\x00")
				if blosta[0] == "-" {
					bloqueDeCarpetas.B_content[j] = bcontent
					escribir_BC(bloqueDeCarpetas, apunt, l)
					listo = true
					break
				}
			}
		} else {
			Bcarpeta = new_index_bloque(l)
			copy(raiz.I_block[i:i+3], strconv.Itoa(Bcarpeta))
			bloqueDeCarpeta := newCarpeta()
			bloqueDeCarpeta.B_content[0] = bcontent
			escribir_inodo(index, raiz, l)
			escribir_BC(bloqueDeCarpeta, Bcarpeta, l)
			listo = true
		}
		if listo {
			break
		}
	}
	return carpetaNew
}
