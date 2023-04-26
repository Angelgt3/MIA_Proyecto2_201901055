package comandos

import (
	"bytes"
	"encoding/binary"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func mkfs(id string, typee string) {
	//validar los parametros
	if id == "" {
		//fmt.Println("ERROR MKFS: FALTA PARAMETRO DE ID")
		respuesta += "\nERROR MKFS: FALTA PARAMETRO DE ID"
		return
	}
	if typee != "full" {
		//fmt.Println("ERROR MKFS: EL VALOR DE TYPE NO ES VALIDO")
		respuesta += "\nERROR MKFS: EL VALOR DE TYPE NO ES VALIDO"
		return
	}

	//se busca la particion montada
	var t Disco
	for i := 0; i < Tmontadas.Len(); i++ {
		t = Tmontadas.GetValue(i).(Disco)
		if t.Id == id { //Se encontro
			ext2(t)
			return
		}
	}
	//fmt.Println("ERROR MKFS: NO SE ENCONTRO LA PARTICION")
	respuesta += "\nERROR MKFS: NO SE ENCONTRO LA PARTICION"
}

func ext2(particion Disco) {
	//se crea el super bloque
	var SB SUPER_BLOQUE
	lenSb := new(bytes.Buffer)
	binary.Write(lenSb, binary.BigEndian, SB)

	var TI TINODOS
	lenTi := new(bytes.Buffer)
	binary.Write(lenTi, binary.BigEndian, TI)

	var BC BLOQUE_CARPETA
	lenBc := new(bytes.Buffer)
	binary.Write(lenBc, binary.BigEndian, BC)

	//numeros_estructuras
	siz := particion.Size
	num := float64((siz - len(lenSb.Bytes())) / (1 + 3 + len(lenTi.Bytes()) + 3*len(lenBc.Bytes())))
	num = math.Floor(num)
	tim := time.Now().Format("2006-01-02 15:04:05")

	copy(SB.S_filesystem_type[:], "2")
	copy(SB.S_inodes_count[:], strconv.Itoa(int(num)))
	copy(SB.S_blocks_count[:], strconv.Itoa(3*int(num)))
	copy(SB.S_free_blocks_count[:], strconv.Itoa((3*int(num))-2))
	copy(SB.S_free_inodes_count[:], strconv.Itoa((int(num))-2))
	copy(SB.S_mtime[:], tim)
	copy(SB.S_mnt_count[:], "1")
	copy(SB.S_magic[:], strconv.Itoa(0xEF53))
	copy(SB.S_inode_size[:], strconv.Itoa(len(lenTi.Bytes())))
	copy(SB.S_block_size[:], "64")
	copy(SB.S_first_ino[:], "2")
	copy(SB.S_first_blo[:], "2")

	//inicio del bitmaps
	Pstart := strings.Split(string(particion.Part.Part_start[:]), "\x00")
	Ps, _ := strconv.Atoi(Pstart[0])
	copy(SB.S_bm_inode_start[:], strconv.Itoa(Ps+len(lenSb.Bytes())))
	Bstar := Ps + len(lenSb.Bytes())
	copy(SB.S_bm_block_start[:], strconv.Itoa(Bstar+int(num)))
	copy(SB.S_inode_start[:], strconv.Itoa(Bstar+int(num)+(3*int(num))))
	copy(SB.S_block_start[:], strconv.Itoa(Bstar+int(num)+(3*int(num))+int(num)*len(lenTi.Bytes())))

	//creacion del bitmap de inodos
	BitInodo := make([]byte, int(num))
	copy(BitInodo[0:1], "1")
	copy(BitInodo[1:2], "1")
	for i := 2; i < int(num); i++ {
		copy(BitInodo[i:i+1], "0")
	}

	//creacion del bitmap de bloques
	BitBloque := make([]byte, 3*int(num))
	copy(BitBloque[0:1], "1")
	copy(BitBloque[1:2], "1")
	for i := 2; i < 3*int(num); i++ {
		copy(BitBloque[i:i+1], "0")
	}

	//escribir todo en el archivo dsk
	archivo, err := os.OpenFile(particion.Path, os.O_RDWR, 0660) // Apertura del archivo
	if err != nil {
		//fmt.Println("ERROR MKFS: NO SE LOGRO ABRIR EL ARCHIVO")
		respuesta += "\nERROR MKFS: NO SE LOGRO ABRIR EL ARCHIVO"
	}
	// se escribe le super bloque
	archivo.Seek(int64(Ps), 0) // Posicion inicial de la particion
	rre := binary.Write(archivo, binary.BigEndian, SB)
	if rre != nil {
		print(rre)
	}

	// se escribe el bitmap de inodos
	rre = binary.Write(archivo, binary.BigEndian, BitInodo)
	if rre != nil {
		print(rre)
	}

	//se escribe el bitmap de bloques
	rre = binary.Write(archivo, binary.BigEndian, BitBloque)
	if rre != nil {
		print(rre)
	}

	//se escribe el inodo
	var inodo = newInodo() //se crea un nuevo inodo
	for i := 0; i < int(num); i++ {
		rre = binary.Write(archivo, binary.BigEndian, inodo)
		if rre != nil {
			print(rre)
		}
	}
	// se escribe el espacio de los bloques con 0
	Bloqus := make([]byte, 64)
	for i := 0; i < 64; i++ {
		copy(Bloqus[i:i+1], "0")
	}

	for i := 0; i < 3*int(num); i++ {
		rre = binary.Write(archivo, binary.BigEndian, Bloqus)
		if rre != nil {
			print(rre)
		}
	}
	archivo.Close()
	// se crea el inodo raiz
	var InoCarpetaRaiz = newInodo()
	copy(InoCarpetaRaiz.I_type[:], "0")
	copy(InoCarpetaRaiz.I_uid[:], "1")
	copy(InoCarpetaRaiz.I_gid[:], "1")
	copy(InoCarpetaRaiz.I_size[:], "0")
	copy(InoCarpetaRaiz.I_atime[:], time.Now().Format("2006-01-02 15:04:05"))
	copy(InoCarpetaRaiz.I_ctime[:], time.Now().Format("2006-01-02 15:04:05"))
	copy(InoCarpetaRaiz.I_mtime[:], time.Now().Format("2006-01-02 15:04:05"))
	copy(InoCarpetaRaiz.I_block[0:1], "0")
	copy(InoCarpetaRaiz.I_perm[:], "664")

	// se crea el bloque carpeta de la raiz
	var carpetaRaiz = newCarpeta()
	var contenido CONTENIDO
	copy(contenido.B_name[:], ".")
	copy(contenido.B_inodo[:], "0")
	carpetaRaiz.B_content[0] = contenido
	copy(contenido.B_name[:], "..")
	carpetaRaiz.B_content[1] = contenido

	//inodo archivo user
	var InoArchivoUser = newInodo()
	copy(InoArchivoUser.I_type[:], "1")
	copy(InoArchivoUser.I_uid[:], "1")
	copy(InoArchivoUser.I_gid[:], "1")
	copy(InoArchivoUser.I_size[:], "28")
	copy(InoArchivoUser.I_atime[:], time.Now().Format("2006-01-02 15:04:05"))
	copy(InoArchivoUser.I_ctime[:], time.Now().Format("2006-01-02 15:04:05"))
	copy(InoArchivoUser.I_mtime[:], time.Now().Format("2006-01-02 15:04:05"))
	copy(InoArchivoUser.I_block[0:1], "1")
	copy(InoArchivoUser.I_perm[:], "664")

	// bloque de archivo de users.txt
	var archivoUser BLOQUE_ARCHIVO
	// contenido del archivo user
	contenidoUsers := ""
	contenidoUsers += "1, G, root\n"
	contenidoUsers += "1, U, root, root, 123"
	copy(archivoUser.B_content[:], contenidoUsers)

	// carpeta que contiene el archivo
	copy(contenido.B_name[:], "users.txt")
	copy(contenido.B_inodo[:], "1")
	carpetaRaiz.B_content[2] = contenido

	//se escribe los inodos en el archivo dks
	escribir_inodo(0, InoCarpetaRaiz, particion)
	escribir_inodo(1, InoArchivoUser, particion)

	//se escribe los bloques en el archivo
	archivo, err = os.OpenFile(particion.Path, os.O_RDWR, 0660) // Apertura del archivo
	if err != nil {
		//fmt.Println("ERROR: NO SE LOGRO ABRIR EL ARCHIVO")
		respuesta += "\nERROR MKFS: NO SE LOGRO ABRIR EL ARCHIVO"
	}
	bstar := strings.Split(string(SB.S_block_start[:]), "\x00")
	Bs, _ := strconv.Atoi(bstar[0])
	archivo.Seek(int64(Bs), 0)
	binary.Write(archivo, binary.BigEndian, carpetaRaiz)
	binary.Write(archivo, binary.BigEndian, archivoUser)
	archivo.Close()
	//fmt.Println("SE REALIZO EL MKFS CON EXITO")
	respuesta += "\nSE REALIZO EL MKFS CON EXITO"
}
