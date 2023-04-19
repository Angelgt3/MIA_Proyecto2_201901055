package comandos

import (
	"fmt"
)

func Crear_reporte(name string, path string, id string, ruta string) {
	if name == "disk" {
		fmt.Println("REPORTE DE DISCO")
		//reporte_disco()
	}
}

/*
func reporte_disco() {
	ls := BuscarPath(d.id)
	var rbm expresion.MBR
	// ----- Apertura de archivo para obtener mbr
	disco, err := os.OpenFile(ls, os.O_RDWR, 0660)
	if err != nil {
		msg_error(err)
	}
	disco.Seek(0, 0)
	rre := binary.Read(disco, binary.BigEndian, &rbm) // obtiene el mbr
	if rre != nil {
		print(rre)
	}
	disco.Close()
	dotEx := ""
	colr := 0
	dotText := "digraph disco{ \n"
	dotText += "contenido [\n shape=plaintext \n label=< \n"
	dotText += "<table BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" color =\"blue\"> \n"
	dotText += "<tr><td>\n"
	dotText += "<table color=\"blue\" border =\"0\" cellborder=\"1\" cellpadding =\"10\"  cellspacing=\"0\">\n"
	dotText += "<tr><td rowspan =\"2\">MBR</td>\n"
	if rbm.Mbr_partition_1.Part_status[0] != '1' {
		dotText += "<td rowspan =\"2\">LIBRE<br/>0% del disco</td>\n"
	} else {
		if rbm.Mbr_partition_1.Part_type[0] == 'E' {
			var rbe expresion.EBR
			disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
			if err != nil {
				msg_error(err)
			}

			pstar, _ := strconv.Atoi(strings.Split(string(rbm.Mbr_partition_1.Part_start[:]), "\x00")[0])
			disco.Seek(int64(pstar), 0)
			rre = binary.Read(disco, binary.BigEndian, &rbe) // obtiene el mbr
			if rre != nil {
				print(rre)
			}
			disco.Close()
			if rbe.Part_status[0] == '0' {
				pnext, _ := strconv.Atoi(strings.Split(string(rbe.Part_next[:]), "\x00")[0])
				if pnext != -1 {
					disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
					if err != nil {
						msg_error(err)
					}
					disco.Seek(int64(pnext), 0)
					rre = binary.Read(disco, binary.BigEndian, &rbe) // obtiene el mbr
					if rre != nil {
						print(rre)
					}
					disco.Close()
				}
			}
			// SE VA POR TODAS LAS PARTICIONES
			fin := false
			for fin == false {
				dotEx += "<td>EBR</td>\n"
				eman := strings.Split(string(rbe.Part_name[:]), "\x00")[0]
				dotEx += "<td>" + eman + "<br/>" + "50" + "% del disco</td>\n"
				colr++
				pnext, _ := strconv.Atoi(strings.Split(string(rbe.Part_next[:]), "\x00")[0])
				if pnext == -1 {
					fin = true
				} else {
					disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
					if err != nil {
						msg_error(err)
					}
					disco.Seek(int64(pnext), 0)
					rre = binary.Read(disco, binary.BigEndian, &rbe) // obtiene el mbr
					if rre != nil {
						print(rre)
					}
					disco.Close()
				}

			}
			dotText += "<td colspan =\"" + strconv.Itoa(colr*2) + "\">EXTENDIDA</td>\n"
		} else {
			eman := strings.Split(string(rbm.Mbr_partition_1.Part_name[:]), "\x00")[0]
			dotText += "<td rowspan =\"2\">" + eman + "<br/>" + "0" + "% del disco</td>\n"
		}
	}

	if rbm.Mbr_partition_2.Part_status[0] != '1' {
		dotText += "<td rowspan =\"2\">LIBRE<br/>0% del disco</td>\n"
	} else {
		if rbm.Mbr_partition_2.Part_type[0] == 'E' {
			var rbe expresion.EBR
			disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
			if err != nil {
				msg_error(err)
			}

			pstar, _ := strconv.Atoi(strings.Split(string(rbm.Mbr_partition_2.Part_start[:]), "\x00")[0])
			disco.Seek(int64(pstar), 0)
			rre = binary.Read(disco, binary.BigEndian, &rbe) // obtiene el mbr
			if rre != nil {
				print(rre)
			}
			disco.Close()
			if rbe.Part_status[0] == '0' {
				pnext, _ := strconv.Atoi(strings.Split(string(rbe.Part_next[:]), "\x00")[0])
				if pnext != -1 {
					disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
					if err != nil {
						msg_error(err)
					}
					disco.Seek(int64(pnext), 0)
					rre = binary.Read(disco, binary.BigEndian, &rbe) // obtiene el mbr
					if rre != nil {
						print(rre)
					}
					disco.Close()
				}
			}
			// SE VA POR TODAS LAS PARTICIONES
			fin := false
			for fin == false {
				dotEx += "<td>EBR</td>\n"
				eman := strings.Split(string(rbe.Part_name[:]), "\x00")[0]
				dotEx += "<td>" + eman + "<br/>" + "50" + "% del disco</td>\n"
				colr++
				pnext, _ := strconv.Atoi(strings.Split(string(rbe.Part_next[:]), "\x00")[0])
				if pnext == -1 {
					fin = true
				} else {
					disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
					if err != nil {
						msg_error(err)
					}
					disco.Seek(int64(pnext), 0)
					rre = binary.Read(disco, binary.BigEndian, &rbe) // obtiene el mbr
					if rre != nil {
						print(rre)
					}
					disco.Close()
				}

			}
			dotText += "<td colspan =\"" + strconv.Itoa(colr*2) + "\">EXTENDIDA</td>\n"
		} else {
			eman := strings.Split(string(rbm.Mbr_partition_2.Part_name[:]), "\x00")[0]
			dotText += "<td rowspan =\"2\">" + eman + "<br/>" + "0" + "% del disco</td>\n"
		}
	}

	if rbm.Mbr_partition_3.Part_status[0] != '1' {
		dotText += "<td rowspan =\"2\">LIBRE<br/>0% del disco</td>\n"
	} else {
		if rbm.Mbr_partition_3.Part_type[0] == 'E' {
			var rbe expresion.EBR
			disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
			if err != nil {
				msg_error(err)
			}

			pstar, _ := strconv.Atoi(strings.Split(string(rbm.Mbr_partition_3.Part_start[:]), "\x00")[0])
			disco.Seek(int64(pstar), 0)
			rre = binary.Read(disco, binary.BigEndian, &rbe) // obtiene el mbr
			if rre != nil {
				print(rre)
			}
			disco.Close()
			if rbe.Part_status[0] == '0' {
				pnext, _ := strconv.Atoi(strings.Split(string(rbe.Part_next[:]), "\x00")[0])
				if pnext != -1 {
					disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
					if err != nil {
						msg_error(err)
					}
					disco.Seek(int64(pnext), 0)
					rre = binary.Read(disco, binary.BigEndian, &rbe) // obtiene el mbr
					if rre != nil {
						print(rre)
					}
					disco.Close()
				}
			}
			// SE VA POR TODAS LAS PARTICIONES
			fin := false
			for fin == false {
				dotEx += "<td>EBR</td>\n"
				eman := strings.Split(string(rbe.Part_name[:]), "\x00")[0]
				dotEx += "<td>" + eman + "<br/>" + "50" + "% del disco</td>\n"
				colr++
				pnext, _ := strconv.Atoi(strings.Split(string(rbe.Part_next[:]), "\x00")[0])
				if pnext == -1 {
					fin = true
				} else {
					disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
					if err != nil {
						msg_error(err)
					}
					disco.Seek(int64(pnext), 0)
					rre = binary.Read(disco, binary.BigEndian, &rbe) // obtiene el mbr
					if rre != nil {
						print(rre)
					}
					disco.Close()
				}

			}
			dotText += "<td colspan =\"" + strconv.Itoa(colr*2) + "\">EXTENDIDA</td>\n"
		} else {
			eman := strings.Split(string(rbm.Mbr_partition_3.Part_name[:]), "\x00")[0]
			dotText += "<td rowspan =\"2\">" + eman + "<br/>" + "0" + "% del disco</td>\n"
		}
	}

	if rbm.Mbr_partition_4.Part_status[0] != '1' {
		dotText += "<td rowspan =\"2\">LIBRE<br/>0% del disco</td>\n"
	} else {
		if rbm.Mbr_partition_4.Part_type[0] == 'E' {
			var rbe expresion.EBR
			disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
			if err != nil {
				msg_error(err)
			}

			pstar, _ := strconv.Atoi(strings.Split(string(rbm.Mbr_partition_4.Part_start[:]), "\x00")[0])
			disco.Seek(int64(pstar), 0)
			rre = binary.Read(disco, binary.BigEndian, &rbe) // obtiene el mbr
			if rre != nil {
				print(rre)
			}
			disco.Close()
			if rbe.Part_status[0] == '0' {
				pnext, _ := strconv.Atoi(strings.Split(string(rbe.Part_next[:]), "\x00")[0])
				if pnext != -1 {
					disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
					if err != nil {
						msg_error(err)
					}
					disco.Seek(int64(pnext), 0)
					rre = binary.Read(disco, binary.BigEndian, &rbe) // obtiene el mbr
					if rre != nil {
						print(rre)
					}
					disco.Close()
				}
			}
			// SE VA POR TODAS LAS PARTICIONES
			fin := false
			for fin == false {
				dotEx += "<td>EBR</td>\n"
				eman := strings.Split(string(rbe.Part_name[:]), "\x00")[0]
				dotEx += "<td>" + eman + "<br/>" + "50" + "% del disco</td>\n"
				colr++
				pnext, _ := strconv.Atoi(strings.Split(string(rbe.Part_next[:]), "\x00")[0])
				if pnext == -1 {
					fin = true
				} else {
					disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
					if err != nil {
						msg_error(err)
					}
					disco.Seek(int64(pnext), 0)
					rre = binary.Read(disco, binary.BigEndian, &rbe) // obtiene el mbr
					if rre != nil {
						print(rre)
					}
					disco.Close()
				}

			}
			dotText += "<td colspan =\"" + strconv.Itoa(colr*2) + "\">EXTENDIDA</td>\n"
		} else {
			eman := strings.Split(string(rbm.Mbr_partition_4.Part_name[:]), "\x00")[0]
			dotText += "<td rowspan =\"2\">" + eman + "<br/>" + "0" + "% del disco</td>\n"
		}
	}

	dotText += "</tr>\n"
	if dotEx != "" {
		dotText += "<tr>\n"
		dotText += dotEx
		dotText += "</tr>\n"
	}
	dotText += "</table>\n"
	dotText += "</td></tr>\n"
	dotText += "</table>\n"
	dotText += ">]\n }"

	CrearRep(dotText, d.path)
}

// REPORTE FILE
func RepFile(d Rep) {
	ls := BuscarLista(d.id)
	indiceInodoArchivo := expresion.ExisteRuta(d.ruta, ls, 0)
	if indiceInodoArchivo == -1 {
		fmt.Println("Error: No se encontro el archivo en esa ruta.")
		return
	}
	contenidoArchivo := expresion.Leerarchivo(indiceInodoArchivo, ls)

	dot := "digraph { \n"
	dot += "rankdir = LR \n"
	dot += "node[shape = record] \n"
	dot += "struct [ \n"
	dot += "label = "
	dot += "\""
	dot += expresion.NombreArchivo(d.ruta) + "|"
	dot += contenidoArchivo
	dot += "\""
	dot += "] \n"
	dot += "}"

	CrearRep(dot, d.path)
}

// Busca el path por el id
func BuscarPath(id string) string {
	for i := 0; i < expresion.Tlista.Len(); i++ {
		r := expresion.Tlista.GetValue(i).(expresion.Tlist)
		if r.Id == id {
			return r.Path
		}
	}
	return ""
}
*/
