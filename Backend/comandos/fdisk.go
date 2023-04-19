package comandos

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func fdisk(size int, unit string, path string, typee string, fit string, name string) {
	// verificar si existe la ruta
	existeP := archivoExiste(path)
	if !existeP {
		fmt.Println("ERROR: NO EXISTE LA RUTA DEL DISCO")
		return
	}

	//crear una particion
	//se abre el archivo
	disco, err := os.OpenFile(path, os.O_RDWR, 0660)
	if err != nil {
		panic(err)
	}
	var mbr MBR
	disco.Seek(0, 0)                                  // se ubica al inicio del archivo
	rre := binary.Read(disco, binary.BigEndian, &mbr) // se convierte en arreglo de byte
	if rre != nil {
		print(rre)
	}
	disco.Close()
	//convierto el tam a bytes
	var tam int
	if unit == "k" {
		tam = size * 1024
	} else if unit == "m" {
		tam = size * 1024 * 1024
	} else if unit == "b" {
		tam = size
	}

	//verficar que no haya una particion con el mismo nombre
	var repetido bool = false
	if name == string(mbr.Mbr_partition_1.Part_name[:]) {
		fmt.Println(name)
		fmt.Println(string(mbr.Mbr_partition_1.Part_name[:]))
		repetido = true
	} else if name == string(mbr.Mbr_partition_2.Part_name[:]) {
		fmt.Println(name)
		fmt.Println(string(mbr.Mbr_partition_2.Part_name[:]))
		repetido = true
	} else if name == string(mbr.Mbr_partition_3.Part_name[:]) {
		fmt.Println(name)
		fmt.Println(string(mbr.Mbr_partition_3.Part_name[:]))
		repetido = true
	} else if name == string(mbr.Mbr_partition_4.Part_name[:]) {
		fmt.Println(name)
		fmt.Println(string(mbr.Mbr_partition_4.Part_name[:]))
		repetido = true
	}
	if repetido {
		fmt.Println("ERROR FDISK: No se puede tener particiones con el mismo nombre")
		return
	}

	//decide que tipo de paricion crear
	if typee == "p" { //crea particiones primarias
		crear_particiones(mbr, tam, path, name, typee, fit)
	}

}

func crear_particiones(mbr MBR, tam int, path string, name string, typee string, fit string) {
	/*
		fmt.Println(string(mbr.Dsk_fit[:]))
		fmt.Println(string(mbr.Mbr_dsk_signature[:]))
		fmt.Println(string(mbr.Mbr_fecha_creacion[:]))
	*/
	// Buscar particiones libres
	var part [10]byte
	var p1, p2, p3, p4, ex bool

	if mbr.Mbr_partition_1.Part_status == part { //particion 1
		p1 = true
	} else if string(mbr.Mbr_partition_1.Part_type[0]) == "e" && typee == "e" {
		fmt.Println("ERROR: YA EXISTE UNA PARTICION EXTENDIDA")
		return
	}
	if mbr.Mbr_partition_2.Part_status == part {
		p2 = true
	} else if string(mbr.Mbr_partition_2.Part_type[0]) == "e" && typee == "e" {
		fmt.Println("ERROR: YA EXISTE UNA PARTICION EXTENDIDA")
		return
	}

	if mbr.Mbr_partition_3.Part_status == part {
		p3 = true
	} else if string(mbr.Mbr_partition_3.Part_type[0]) == "e" && typee == "e" {
		fmt.Println("ERROR: YA EXISTE UNA PARTICION EXTENDIDA")
		return
	}

	if mbr.Mbr_partition_4.Part_status == part {
		p4 = true
	} else if string(mbr.Mbr_partition_4.Part_type[0]) == "e" && typee == "e" {
		fmt.Println("ERROR: YA EXISTE UNA PARTICION EXTENDIDA")
		return
	}
	if !p1 && !p2 && !p3 && !p4 {
		fmt.Println("ERROR: YA EXISTEN 4 PARTICIONES")
		return
	}
	if typee == "e" {
		ex = true
	}

	calc := newAjuste()
	num := 0
	calc.Inicio[0] = 0 // 0
	calc.Fin[0] = 280  //len(bf.Bytes())

	if !p1 {

		res1 := strings.Split(string(mbr.Mbr_partition_1.Part_start[:]), "\x00")
		num, _ = strconv.Atoi(res1[0])

		calc.Inicio[1] = num

		res2 := strings.Split(string(mbr.Mbr_partition_1.Part_size[:]), "\x00")
		num, _ = strconv.Atoi(res2[0])

		calc.Fin[1] = num + calc.Inicio[1]
	}
	if !p2 {

		res1 := strings.Split(string(mbr.Mbr_partition_2.Part_start[:]), "\x00")
		num, _ = strconv.Atoi(res1[0])

		calc.Inicio[2] = num

		res2 := strings.Split(string(mbr.Mbr_partition_2.Part_size[:]), "\x00")
		num, _ = strconv.Atoi(res2[0])

		calc.Fin[2] = num + calc.Inicio[2]
	}
	if !p3 {

		res1 := strings.Split(string(mbr.Mbr_partition_3.Part_start[:]), "\x00")
		num, _ = strconv.Atoi(res1[0])

		calc.Inicio[3] = num

		res2 := strings.Split(string(mbr.Mbr_partition_3.Part_size[:]), "\x00")
		num, _ = strconv.Atoi(res2[0])

		calc.Fin[3] = num + calc.Inicio[3]
	}
	if !p4 {
		res1 := strings.Split(string(mbr.Mbr_partition_4.Part_start[:]), "\x00")
		num, _ = strconv.Atoi(res1[0])

		calc.Inicio[4] = num

		res2 := strings.Split(string(mbr.Mbr_partition_4.Part_size[:]), "\x00")
		num, _ = strconv.Atoi(res2[0])

		calc.Fin[4] = num + calc.Inicio[4]
	}

	tamMbr := strings.Split(string(mbr.Mbr_tamano[:]), "\x00")
	num, _ = strconv.Atoi(tamMbr[0])
	calc.Inicio[5] = num
	calc.Fin[5] = num

	// Ordenar de menor a mayor
	var in = 0
	var fn = 0

	for i := 1; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if calc.Fin[j] > calc.Fin[j+1] {
				in = calc.Inicio[j]
				fn = calc.Fin[j]
				calc.Inicio[j] = calc.Inicio[j+1]
				calc.Fin[j] = calc.Fin[j+1]
				calc.Inicio[j+1] = in
				calc.Fin[j+1] = fn
			}
		}
	}

	//Calcula el espacio entre particiones
	// Y el espacio que ocupa la nueva entre los espacios

	var hayS = false
	for i := 0; i < 5; i++ {
		calc.Inter[i] = calc.Inicio[i+1] - calc.Fin[i]
		calc.Resta[i] = calc.Inter[i] - tam
		if !hayS && calc.Resta[i] > 1 {
			hayS = true
		}
	}
	if !hayS {
		fmt.Println("ERROR: NO HAY ESPACIO SUFICIENTE PARA LA PARTICION")
		return
	}

	// ----- Busca la particion con mejor ajuste
	star := ""
	if fit == "bf" || fit == "wf" {
		// Ordenar de menor a mayor
		r := 0
		for i := 0; i < 5; i++ {
			for j := 0; j < 4; j++ {
				if calc.Resta[j] > calc.Resta[j+1] {
					in = calc.Inicio[j]
					fn = calc.Fin[j]
					r = calc.Resta[j]
					calc.Inicio[j] = calc.Inicio[j+1]
					calc.Fin[j] = calc.Fin[j+1]
					calc.Resta[j] = calc.Resta[j+1]
					calc.Inicio[j+1] = in
					calc.Fin[j+1] = fn
					calc.Resta[j+1] = r
				}
			}
		}
		if fit == "bf" {
			for i := 0; i < 5; i++ {
				if calc.Resta[i] > 1 {
					star = strconv.Itoa(calc.Fin[i] + 1)
					break
				}
			}
		} else {
			for i := 4; i >= 0; i++ {
				if calc.Resta[i] > 1 {
					star = strconv.Itoa(calc.Fin[i] + 1)
					break
				}
			}
		}
	}

	if fit == "f" {
		for i := 0; i < 5; i++ {
			if calc.Resta[i] > 1 {
				star = strconv.Itoa(calc.Fin[i] + 1)
				break
			}
		}
	}
	// Por si es extendida
	if ex {
		rbe := newEbr("0", "", "", "", "-1", "")
		intVar, _ := strconv.Atoi(star)
		disco, err := os.OpenFile(string(path), os.O_RDWR, 0660) // Apertura del archivo
		if err != nil {
			fmt.Println("ERROR: NO SE LOGRO ABRIR EL ARCHIVO")
		}

		disco.Seek(int64(intVar), 0)
		//bf := new(bytes.Buffer)
		binary.Write(disco, binary.BigEndian, rbe) // se convierte en arreglo de byte
		disco.Close()
	}
	if p1 {
		// para guardar solo en una particion la info
		p2 = false
		p3 = false
		p4 = false
		cad := strconv.Itoa(tam)
		copy(mbr.Mbr_partition_1.Part_start[:], star)
		copy(mbr.Mbr_partition_1.Part_fit[:], fit)
		copy(mbr.Mbr_partition_1.Part_name[:], name)
		copy(mbr.Mbr_partition_1.Part_size[:], cad)
		copy(mbr.Mbr_partition_1.Part_status[:], "1")
		copy(mbr.Mbr_partition_1.Part_type[:], typee)
		guardaMBR(mbr, path)
		fmt.Println("SE CREO LA PARTICION " + name + " EN " + path + " CON EXITO")
		return
	}

	if p2 {
		// para guardar solo en una particion la info
		p1 = false
		p3 = false
		p4 = false
		cad := strconv.Itoa(tam)
		copy(mbr.Mbr_partition_2.Part_start[:], star)
		copy(mbr.Mbr_partition_2.Part_fit[:], fit)
		copy(mbr.Mbr_partition_2.Part_name[:], name)
		copy(mbr.Mbr_partition_2.Part_size[:], cad)
		copy(mbr.Mbr_partition_2.Part_status[:], "1")
		copy(mbr.Mbr_partition_2.Part_type[:], typee)
		guardaMBR(mbr, path)
		fmt.Println("SE CREO LA PARTICION " + name + " EN " + path + " CON EXITO")
		return
	}

	if p3 {
		// para guardar solo en una particion la info
		p2 = false
		p1 = false
		p4 = false
		cad := strconv.Itoa(tam)
		copy(mbr.Mbr_partition_3.Part_start[:], star)
		copy(mbr.Mbr_partition_3.Part_fit[:], fit)
		copy(mbr.Mbr_partition_3.Part_name[:], name)
		copy(mbr.Mbr_partition_3.Part_size[:], cad)
		copy(mbr.Mbr_partition_3.Part_status[:], "1")
		copy(mbr.Mbr_partition_3.Part_type[:], typee)
		guardaMBR(mbr, path)
		fmt.Println("SE CREO LA PARTICION " + name + " EN " + path + " CON EXITO")
		return
	}

	if p4 {
		// para guardar solo en una particion la info
		p2 = false
		p3 = false
		p1 = false
		cad := strconv.Itoa(tam)
		copy(mbr.Mbr_partition_4.Part_start[:], star)
		copy(mbr.Mbr_partition_4.Part_fit[:], fit)
		copy(mbr.Mbr_partition_4.Part_name[:], name)
		copy(mbr.Mbr_partition_4.Part_size[:], cad)
		copy(mbr.Mbr_partition_4.Part_status[:], "1")
		copy(mbr.Mbr_partition_4.Part_type[:], typee)
		guardaMBR(mbr, path)
		fmt.Println("SE CREO LA PARTICION " + name + " EN " + path + " CON EXITO")
		return
	}

}
func guardaMBR(mbr MBR, path string) {
	// -------------------------- GUARDA EL MBR
	disco, err := os.OpenFile(string(path), os.O_RDWR, 0660) // Apertura del archivo
	if err != nil {
		fmt.Println("ERROR: NO SE LOGRO ABRIR EL ARCHIVO")
	}
	disco.Seek(0, 0) // Posicion inicial
	//bf := new(bytes.Buffer)
	rre := binary.Write(disco, binary.BigEndian, mbr) // se convierte en arreglo de byte
	if rre != nil {
		print(rre)
	}
	disco.Close()
}
