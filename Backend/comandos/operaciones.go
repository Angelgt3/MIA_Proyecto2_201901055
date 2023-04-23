package comandos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
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
func newBloque(disco Disco) int {
	bm := get_bitmap(disco, false)
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
func get_usuarios(usuarios string) [][]string {
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

// retorna el bloque carpeta deseado por medio de un index
func get_bloque_carpeta(index int, disco Disco) BLOQUE_CARPETA {
	var sb SUPER_BLOQUE
	puntero := 64 * index
	archivo, err := os.OpenFile(disco.Path, os.O_RDWR, 0660) // Apertura del archivo
	if err != nil {
		fmt.Println("ERROR: NO SE LOGRO ABRIR EL ARCHIVO (get_bloque_carpeta)")
	}
	Pstart := strings.Split(string(disco.Part.Part_start[:]), "\x00")
	Ps, _ := strconv.Atoi(Pstart[0])
	archivo.Seek(int64(Ps), 0) // Posicion inicial

	rre := binary.Read(archivo, binary.BigEndian, &sb) // se convierte en arreglo de byte
	if rre != nil {
		print(rre)
	}
	var Bcarpeta BLOQUE_CARPETA
	// SE OBTIENE EL BLOQUE DEL INDICE
	inosta := strings.Split(string(sb.S_block_start[:]), "\x00")
	is, _ := strconv.Atoi(inosta[0])
	archivo.Seek(int64(is), 0)
	archivo.Seek(int64(puntero), 1)
	binary.Read(archivo, binary.BigEndian, &Bcarpeta)
	archivo.Close()
	return Bcarpeta
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
			pos_ba := newBloque(disco)
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
