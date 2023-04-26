package comandos

import (
	"encoding/binary"
	"os"
	"strconv"
	"strings"
)

func fdisk(size int, unit string, path string, typee string, fit string, name string) {
	// verificar si existe la ruta
	existeP := archivoExiste(path)
	if !existeP {
		//fmt.Println("ERROR: NO EXISTE LA RUTA DEL DISCO")
		respuesta += "\nERROR: NO EXISTE LA RUTA DEL DISCO"
		return
	}

	//crear una particion
	//se abre el archivo
	disco, err := os.OpenFile(path, os.O_RDWR, 0660)
	if err != nil {
		panic(err)
	}
	var mbr MBR
	disco.Seek(0, 0)
	rre := binary.Read(disco, binary.BigEndian, &mbr)
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
		repetido = true
	} else if name == string(mbr.Mbr_partition_2.Part_name[:]) {
		repetido = true
	} else if name == string(mbr.Mbr_partition_3.Part_name[:]) {
		repetido = true
	} else if name == string(mbr.Mbr_partition_4.Part_name[:]) {
		repetido = true
	}
	if repetido {
		//fmt.Println("ERROR FDISK: No se puede tener particiones con el mismo nombre")
		respuesta += "\nERROR FDISK: No se puede tener particiones con el mismo nombre"
		return
	}

	//decide que tipo de paricion crear
	if typee == "p" || typee == "e" { //crea particiones primarias y extendidas
		crear_particiones(mbr, tam, path, name, typee, fit)
	} else if typee == "l" { //crea particiones logicas
		crear_logica(mbr, tam, path, name, typee, fit)
	}

}

func crear_particiones(mbr MBR, tam int, path string, name string, typee string, fit string) {
	// Buscar particiones libres
	var part [10]byte
	var p1, p2, p3, p4, ex bool

	if mbr.Mbr_partition_1.Part_status == part { //particion 1
		p1 = true
	} else if string(mbr.Mbr_partition_1.Part_type[0]) == "e" && typee == "e" {
		//fmt.Println("ERROR: YA EXISTE UNA PARTICION EXTENDIDA")
		respuesta += "\nERROR: YA EXISTE UNA PARTICION EXTENDIDA"
		return
	}
	if mbr.Mbr_partition_2.Part_status == part {
		p2 = true
	} else if string(mbr.Mbr_partition_2.Part_type[0]) == "e" && typee == "e" {
		//fmt.Println("ERROR: YA EXISTE UNA PARTICION EXTENDIDA")
		respuesta += "\nERROR: YA EXISTE UNA PARTICION EXTENDIDA"
		return
	}

	if mbr.Mbr_partition_3.Part_status == part {
		p3 = true
	} else if string(mbr.Mbr_partition_3.Part_type[0]) == "e" && typee == "e" {
		//fmt.Println("ERROR: YA EXISTE UNA PARTICION EXTENDIDA")
		respuesta += "\nERROR: YA EXISTE UNA PARTICION EXTENDIDA"
		return
	}

	if mbr.Mbr_partition_4.Part_status == part {
		p4 = true
	} else if string(mbr.Mbr_partition_4.Part_type[0]) == "e" && typee == "e" {
		//fmt.Println("ERROR: YA EXISTE UNA PARTICION EXTENDIDA")
		respuesta += "\nERROR: YA EXISTE UNA PARTICION EXTENDIDA"
		return
	}
	if !p1 && !p2 && !p3 && !p4 {
		//fmt.Println("ERROR: YA EXISTEN 4 PARTICIONES")
		respuesta += "\nERROR: YA EXISTEN 4 PARTICIONES"
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
		//fmt.Println("ERROR: NO HAY ESPACIO SUFICIENTE PARA LA PARTICION")
		respuesta += "\nERROR: NO HAY ESPACIO SUFICIENTE PARA LA PARTICION"
		return
	}
	// Busca la particion con mejor ajuste
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
		ebr := newEbr("0", "", "", "", "-1", "")
		intVar, _ := strconv.Atoi(star)
		disco, err := os.OpenFile(string(path), os.O_RDWR, 0660) // Apertura del archivo
		if err != nil {
			//fmt.Println("ERROR FDISK: NO SE LOGRO ABRIR EL ARCHIVO")
			respuesta += "\nERROR FDISK: NO SE LOGRO ABRIR EL ARCHIVO"
		}

		disco.Seek(int64(intVar), 0)
		//bf := new(bytes.Buffer)
		binary.Write(disco, binary.BigEndian, ebr) // se convierte en arreglo de byte
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
		//fmt.Println("SE CREO LA PARTICION " + name + " EN " + path + " CON EXITO")
		respuesta += "\nSE CREO LA PARTICION " + name + " EN " + path + " CON EXITO"
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
		//fmt.Println("SE CREO LA PARTICION " + name + " EN " + path + " CON EXITO")
		respuesta += "\nSE CREO LA PARTICION " + name + " EN " + path + " CON EXITO"
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
		//fmt.Println("SE CREO LA PARTICION " + name + " EN " + path + " CON EXITO")
		respuesta += "\nSE CREO LA PARTICION " + name + " EN " + path + " CON EXITO"
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
		//fmt.Println("SE CREO LA PARTICION " + name + " EN " + path + " CON EXITO")
		respuesta += "\nSE CREO LA PARTICION " + name + " EN " + path + " CON EXITO"
		return
	}
}

func crear_logica(mbr MBR, tam int, path string, name string, typee string, fit string) {
	// verificar si hay particion extendida
	var p1, p2, p3, p4 bool

	if string(mbr.Mbr_partition_1.Part_type[0]) == "e" {
		p1 = true
	} else if string(mbr.Mbr_partition_2.Part_type[0]) == "e" {
		p2 = true
	} else if string(mbr.Mbr_partition_3.Part_type[0]) == "e" {
		p3 = true
	} else if string(mbr.Mbr_partition_4.Part_type[0]) == "e" {
		p4 = true
	}
	// Guardan los valores del inicio y el size de la particion ext
	s := ""
	sz := ""
	if p1 {

		res1 := strings.Split(string(mbr.Mbr_partition_1.Part_start[:]), "\x00")
		s = res1[0]

		res2 := strings.Split(string(mbr.Mbr_partition_1.Part_size[:]), "\x00")
		sz = res2[0]

	} else if p2 {

		res1 := strings.Split(string(mbr.Mbr_partition_2.Part_start[:]), "\x00")
		s = res1[0]

		res2 := strings.Split(string(mbr.Mbr_partition_2.Part_size[:]), "\x00")
		sz = res2[0]

	} else if p3 {

		res1 := strings.Split(string(mbr.Mbr_partition_3.Part_start[:]), "\x00")
		s = res1[0]

		res2 := strings.Split(string(mbr.Mbr_partition_3.Part_size[:]), "\x00")
		sz = res2[0]

	} else if p4 {

		res1 := strings.Split(string(mbr.Mbr_partition_4.Part_start[:]), "\x00")
		s = res1[0]

		res2 := strings.Split(string(mbr.Mbr_partition_4.Part_size[:]), "\x00")
		sz = res2[0]

	} else {
		//fmt.Println("ERROR FDISK: NO EXISTE PARTICION EXTENDIDA")
		respuesta += "\nERROR FDISK: NO EXISTE PARTICION EXTENDIDA"
		return
	}

	disco, err := os.OpenFile(path, os.O_RDWR, 0660)
	if err != nil {
		panic(err)
	}

	var ebr EBR
	sta, _ := strconv.Atoi(s)
	total, _ := strconv.Atoi(sz)
	disco.Seek(int64(sta), 0)
	rre := binary.Read(disco, binary.BigEndian, &ebr)
	if rre != nil {
		print(rre)
	}
	disco.Close()

	if string(ebr.Part_status[0]) == "0" { // No hay ebr inicial
		if tam < total { // Si hay espacio disponible
			// Guarda todo
			tl := strconv.Itoa(tam)
			copy(ebr.Part_status[:], "1")
			copy(ebr.Part_fit[:], fit)
			copy(ebr.Part_start[:], s)
			copy(ebr.Part_size[:], tl)
			copy(ebr.Part_next[:], "-1")
			copy(ebr.Part_name[:], []byte(name))
			// Lo guarda en el archivo
			disco, err := os.OpenFile(string(path), os.O_RDWR, 0660) // Apertura del archivo
			if err != nil {
				//fmt.Println("ERROR: NO SE LOGRO ABRIR EL ARCHIVO")
				respuesta += "\nERROR: NO SE LOGRO ABRIR EL ARCHIVO"
			}

			disco.Seek(int64(sta), 0) // Se corre al inicio del ebr
			//bf := new(bytes.Buffer)
			binary.Write(disco, binary.BigEndian, ebr) // se convierte en arreglo de byte
			disco.Close()
			//fmt.Println("SE CREO LA PARTICION LOGICA " + name + " CON EXITO")
			respuesta += "\nSE CREO LA PARTICION LOGICA " + name + " CON EXITO"
		} else {
			//fmt.Println("ERROR FDISK: NO HAY ESPACIO SUFICIENTE PARA CREAR PARTICION LOGICA ")
			respuesta += "\nERROR FDISK: NO HAY ESPACIO SUFICIENTE PARA CREAR PARTICION LOGICA "
		}
	} else {
		var inicio List
		var fin List
		var inter List
		nxt := 0

		disco, err := os.OpenFile(path, os.O_RDWR, 0660)
		if err != nil {
			panic(err)
		}

		for nxt == 0 {
			res1 := strings.Split(string(ebr.Part_start[:]), "\x00")
			res2 := strings.Split(string(ebr.Part_size[:]), "\x00")
			nmi, _ := strconv.Atoi(res1[0])
			nmf, _ := strconv.Atoi(res2[0])
			inicio.Add(nmi)
			fin.Add(nmi + nmf)

			res3 := strings.Split(string(ebr.Part_next[:]), "\x00")
			sig, _ := strconv.Atoi(res3[0])

			if sig != -1 { //ya no hay mas ebr

				disco.Seek(int64(sig), 0)
				rre := binary.Read(disco, binary.BigEndian, &ebr)
				if rre != nil {
					print(rre)
				}

			} else {
				nxt = -1
			}

		}
		disco.Close()

		inicio.Add(sta + total) // Inicio
		fin.Add(sta + total)    // Fin

		for i := 1; i < inicio.Len(); i++ {
			ini := inicio.GetValue(i)
			nif := fin.GetValue(i - 1)
			inter.Add(ini.(int) - nif.(int)) // inicio - fin
		}

		for i := 0; i < inter.Len(); i++ {
			in := inicio.GetValue(i).(int)
			fi := fin.GetValue(i).(int)
			if inter.GetValue(i).(int) > tam {
				emp := fi + 1
				// se apunta a la nueva particion y se guarda
				disco, err = os.OpenFile(path, os.O_RDWR, 0660)
				if err != nil {
					panic(err)
				}
				disco.Seek(int64(in), 0)
				rre := binary.Read(disco, binary.BigEndian, &ebr)
				if rre != nil {
					print(rre)
				}
				cad := strconv.Itoa(emp)
				copy(ebr.Part_next[:], cad)

				disco.Seek(int64(in), 0)
				rre = binary.Write(disco, binary.BigEndian, ebr)
				if rre != nil {
					print(rre)
				}

				tm := strconv.Itoa(tam)
				nuevo := newEbr("1", fit, cad, tm, "-1", name)
				disco.Seek(int64(emp), 0)
				rre = binary.Write(disco, binary.BigEndian, nuevo)
				if rre != nil {
					print(rre)
				}
				//fmt.Println("SE CREO LA PARTICION LOGICA " + name + " CON EXITO")
				respuesta += "\nSE CREO LA PARTICION LOGICA " + name + " CON EXITO"
				disco.Close()
			}
		}

	}

}

func guardaMBR(mbr MBR, path string) {
	// -------------------------- GUARDA EL MBR
	disco, err := os.OpenFile(string(path), os.O_RDWR, 0660) // Apertura del archivo
	if err != nil {
		//fmt.Println("ERROR FDISK: NO SE LOGRO ABRIR EL ARCHIVO")
		respuesta += "\nERROR FDISK: NO SE LOGRO ABRIR EL ARCHIVO"
	}
	disco.Seek(0, 0) // Posicion inicial
	//bf := new(bytes.Buffer)
	rre := binary.Write(disco, binary.BigEndian, mbr) // se convierte en arreglo de byte
	if rre != nil {
		print(rre)
	}
	disco.Close()
}
