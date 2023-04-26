package comandos

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Mount(path string, name string) {
	abc := [27]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "ñ", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	//creo un nuevo struct del nuevo mount
	eslist := newTlist()
	eslist.Name = name
	eslist.Path = path

	//se abre el archivo para obtener mbr
	disco, err := os.OpenFile(path, os.O_RDWR, 0660)
	if err != nil {
		panic(err)
	}

	var mbr MBR
	disco.Seek(0, 0)                                  // se ubica al inicio del archivo
	rre := binary.Read(disco, binary.BigEndian, &mbr) // obtiene el mbr
	if rre != nil {
		print(rre)
	}
	disco.Close()

	// ------------- encuentra la particion
	if strings.Split(string(mbr.Mbr_partition_1.Part_name[:]), "\x00")[0] == name {
		eslist.Part = mbr.Mbr_partition_1
		eslist.Size, _ = strconv.Atoi(strings.Split(string(mbr.Mbr_partition_1.Part_size[:]), "\x00")[0])
	} else if strings.Split(string(mbr.Mbr_partition_2.Part_name[:]), "\x00")[0] == name {
		eslist.Part = mbr.Mbr_partition_2
		eslist.Size, _ = strconv.Atoi(strings.Split(string(mbr.Mbr_partition_2.Part_size[:]), "\x00")[0])
	} else if strings.Split(string(mbr.Mbr_partition_3.Part_name[:]), "\x00")[0] == name {
		eslist.Part = mbr.Mbr_partition_3
		eslist.Size, _ = strconv.Atoi(strings.Split(string(mbr.Mbr_partition_3.Part_size[:]), "\x00")[0])
	} else if strings.Split(string(mbr.Mbr_partition_4.Part_name[:]), "\x00")[0] == name {
		eslist.Part = mbr.Mbr_partition_4
		eslist.Size, _ = strconv.Atoi(strings.Split(string(mbr.Mbr_partition_4.Part_size[:]), "\x00")[0])
	}

	// Si es la Primera Montada
	if Tmontadas.Len() == 0 {
		eslist.Id = "551a"
		eslist.N = 1
		Tmontadas.Add(eslist)
	} else {

		num := 0
		nn := 0
		for i := 0; i < Tmontadas.Len(); i++ {
			t := Tmontadas.GetValue(i).(Disco)
			// Revisa si esta ya montada
			if name == strings.Split(string(t.Part.Part_name[:]), "\x00")[0] {
				println("ERROR: PARTICION YA MONTADA")
				return
			}
			// revisa si hay mas montadas del mismo disco
			if path == t.Path {
				num += 1
				nn = t.N
			}
			// por si no hay montadas del mismo disco
			// obtiene el numero de disco mas grande
			if nn < t.N {
				nn = t.N
			}
		}

		// Crea el id
		if num == 0 { // no existe del mismo disco
			eslist.Id = "55" + strconv.Itoa(nn+1) + "a"
			eslist.N = nn + 1
			Tmontadas.Add(eslist)
		} else { // si hay mas del mismo disco
			eslist.Id = "55" + strconv.Itoa(nn) + abc[num] // diferentw letra
			eslist.N = nn + 1                              // mismo numero
			Tmontadas.Add(eslist)
		}
	}
	//fmt.Println("SE MONTO LA PARTICION " + name + " CON ID: " + eslist.Id)
	respuesta += "\nSE MONTO LA PARTICION " + name + " CON ID: " + eslist.Id
}

func Mostrar_mount() {
	fmt.Println("-----------------------------PARTICIONES MONTADAS------------------------------")
	respuesta += "\n-----------------------------PARTICIONES MONTADAS------------------------------\n"
	fmt.Println("")
	for i := 0; i < Tmontadas.Len(); i++ {
		mon := Tmontadas.GetValue(i).(Disco)
		fmt.Println(">id=" + mon.Id + " >path=" + mon.Path + " >name=" + mon.Name)
		respuesta += "\n>id=" + mon.Id + " >path=" + mon.Path + " >name=" + mon.Name
	}
	fmt.Println("")
	fmt.Println("-------------------------------------------------------------------------------")
	fmt.Println("")
	respuesta += "\n----------------------------------------------------------------------------\n"

}
