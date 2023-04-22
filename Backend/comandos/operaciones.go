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
	for i := 0; i < 16; i++ {
		copy(inodo.I_block[i:i+1], "-")
	}
	return inodo
}

func newCarpeta() BLOQUE_CARPETA {
	var carpeta BLOQUE_CARPETA
	for i := 0; i < 4; i++ {
		copy(carpeta.B_content[i].B_inodo[0:1], "-")
	}
	return carpeta
}

// escribe inodo en el disco por medio de su index
func escribir_inodo(index int, inodo TINODOS, particion Disco) {
	disco, err := os.OpenFile(particion.Path, os.O_RDWR, 0660) // Apertura del archivo
	if err != nil {
		fmt.Println("ERROR: NO SE LOGRO ABRIR EL ARCHIVO (escribir_inood)")
	}
	Pstart := strings.Split(string(particion.Part.Part_start[:]), "\x00")
	Ps, _ := strconv.Atoi(Pstart[0])
	disco.Seek(int64(Ps), 0)

	//modifico el super bloque
	var sb SUPER_BLOQUE
	rre := binary.Read(disco, binary.BigEndian, &sb)
	if rre != nil {
		print(rre)
	}
	ind := index + 1
	ino := strconv.Itoa(ind)
	copy(sb.S_first_ino[:], ino)
	//lo escribo en el archivo
	disco.Seek(int64(Ps), 0)
	binary.Write(disco, binary.BigEndian, sb)

	//mofico el bitmap de inodos
	inodoCon := strings.Split(string(sb.S_inodes_count[:]), "\x00")
	Ic, _ := strconv.Atoi(inodoCon[0])
	bitmap_inodos := make([]byte, int(Ic))
	bmStart := strings.Split(string(sb.S_bm_inode_start[:]), "\x00")
	bm, _ := strconv.Atoi(bmStart[0])
	disco.Seek(int64(bm), 0)
	rre = binary.Read(disco, binary.BigEndian, bitmap_inodos)
	if rre != nil {
		print(rre)
	}
	copy(bitmap_inodos[index:index+1], "1")
	disco.Seek(int64(bm), 0)
	binary.Write(disco, binary.BigEndian, bitmap_inodos)
	//escribo el inodo en el archivo
	inosta := strings.Split(string(sb.S_inode_start[:]), "\x00")
	is, _ := strconv.Atoi(inosta[0])
	disco.Seek(int64(is), 0)
	var TI TINODOS
	lenTi := new(bytes.Buffer)
	binary.Write(lenTi, binary.BigEndian, TI)
	disco.Seek(int64(index*len(lenTi.Bytes())), 1)
	binary.Write(disco, binary.BigEndian, inodo)
	disco.Close()
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
