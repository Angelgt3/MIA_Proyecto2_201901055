package comandos

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Crear_reporte(name string, path string, id string, ruta string) {
	if name == "disk" {
		reporte_disco(id, path)
	} else if name == "sb" {
		Reporte_SB(id, path)
	} else if name == "tree" {
		Reporte_tree(id, path)
		//fmt.Println("SE GENERO EL REPORTE TREE " + id)
		respuesta += "\nSE GENERO EL REPORTE TREE " + id
	} else if name == "file" {
		Reporte_file(path, ruta, id)

	}
}

// reporte file
func Reporte_file(path string, ruta string, id string) {
	disco := get_disco(id)
	indiceInodoArchivo := existe_ruta(ruta, disco, 0)
	if indiceInodoArchivo == -1 {
		//fmt.Println("ERROR REP: NO SE ENCONTRO LA RUTA PARA REALIZAR EL REPORTE FILE")
		respuesta += "\nERROR REP: NO SE ENCONTRO LA RUTA PARA REALIZAR EL REPORTE FILE"
		return
	}
	contenidoArchivo := leer_archivo(indiceInodoArchivo, disco)
	contenido := "digraph { \n"
	contenido += "rankdir = LR \n"
	contenido += "node[shape = record] \n"
	contenido += "struct [ \n"
	contenido += "label = "
	contenido += "\""
	contenido += nombre_archivo(ruta) + "|"
	contenido += contenidoArchivo
	contenido += "\""
	contenido += "] \n"
	contenido += "}"
	crear_dot(contenido, path)
	//fmt.Println("SE GENERO EL REPORTE FILE " + ruta)
	respuesta += "\nSE GENERO EL REPORTE FILE " + ruta
}

// reporte arbol
func Reporte_tree(id string, path string) {
	//encabezador del graphviz
	contenido := "digraph tree{ \n\tlabelloc=\"t\"\n\tlabel=\"REPORTE TREE\"\n\trankdir=LR\n\tnode[shape=record style=\"filled\"]"
	disco := get_disco(id)
	inodo := get_inodo(0, disco) //se obtiene el indo raiz
	contenido += recorrido_arbol(inodo, "/", disco, 0, 0)
	contenido += "\n}"
	crear_dot(contenido, path)
}

func recorrido_arbol(inodo TINODOS, nombre_inodo string, disco Disco, p_inodo int, p_bloque int) string {
	contenido := ""
	for i := 0; i < 64; i = i + 4 {
		inosta := strings.Split(string(inodo.I_block[i:i+3]), "\x00")
		if inosta[0] == "-" { //si el inodo no hat bloque se lo salta
			continue
		}
		//graphviz inodo
		contenido += "\n\tinode" + strconv.Itoa(p_inodo)
		contenido += "[\n\t\tlabel=\"INODO " + strconv.Itoa(p_inodo) + " - " + nombre_inodo
		contenido += "|{i_uid|" + strings.Split(string(inodo.I_uid[:]), "\x00")[0] + "}"
		contenido += "|{i_gid|" + strings.Split(string(inodo.I_gid[:]), "\x00")[0] + "}"
		contenido += "|{i_size|" + strings.Split(string(inodo.I_size[:]), "\x00")[0] + "}"
		contenido += "|{i_atime|" + strings.Split(string(inodo.I_atime[:]), "\x00")[0] + "}"
		contenido += "|{i_ctime|" + strings.Split(string(inodo.I_ctime[:]), "\x00")[0] + "}"
		contenido += "|{i_mtime|" + strings.Split(string(inodo.I_mtime[:]), "\x00")[0] + "}"
		for a := 0; a < 64; a = a + 4 {
			contenido += "|{i_block[" + strconv.Itoa(a) + "]|"
			contenido += strings.Split(string(inodo.I_block[a:a+3]), "\x00")[0] + "}"
		}
		contenido += "|{i_type|" + strings.Split(string(inodo.I_type[:]), "\x00")[0] + "}"
		contenido += "|{i_perm|" + strings.Split(string(inodo.I_perm[:]), "\x00")[0] + "}\"\n\t];"
		if inosta[0] != "-" {
			if string(inodo.I_type[:]) == "0" { //es carpeta
				inb, _ := strconv.Atoi(inosta[0])
				bc := get_bloque_carpeta(inb, disco)
				contenido += "\n\tb" + inosta[0] + "[\n\t\tlabel=\"Bloque " + inosta[0] + " - Carpeta "

				for k := 0; k < 4; k++ {
					contenido += "|{b_content " + strconv.Itoa(k) + "}|{b_name|"
					contenido += strings.Split(string(bc.B_content[k].B_name[:]), "\x00")[0]
					contenido += "}|{b_inodo|" + strings.Split(string(bc.B_content[k].B_inodo[:]), "\x00")[0] + "}"
				}
				contenido += "\"\n\t];"
				//conexiones inodo -> bloques
				contenido += "\n\tinode" + strconv.Itoa(p_inodo) + "->b" + inosta[0]
				//ahora recorrer los inodos de los bloques
				for a := 0; a < 4; a++ {
					if strings.Split(string(bc.B_content[a].B_inodo[:]), "\x00")[0] != "-" && strings.Split(string(bc.B_content[a].B_inodo[:]), "\x00")[0] != "0" && strings.Split(string(bc.B_content[a].B_name[:]), "\x00")[0] != "." && strings.Split(string(bc.B_content[a].B_name[:]), "\x00")[0] != ".." {
						inbb, _ := strconv.Atoi(strings.Split(string(bc.B_content[a].B_inodo[:]), "\x00")[0])
						inb, _ := strconv.Atoi(inosta[0])
						sig_inodo := get_inodo(inbb, disco)
						contenido += "\n\tb" + inosta[0] + "->inode" + strings.Split(string(bc.B_content[a].B_inodo[:]), "\x00")[0]
						contenido += recorrido_arbol(sig_inodo, strings.Split(string(bc.B_content[a].B_name[:]), "\x00")[0], disco, inbb, inb)
					}
				}
			} else if string(inodo.I_type[:]) == "1" { //es archivo
				inb, _ := strconv.Atoi(inosta[0])
				ba := get_bloque_archivo(inb, disco)
				contenido += "\n\tb" + string(inosta[0]) + "[\n\t\tlabel=\"Bloque " + inosta[0] + " - Archivo |b_content|"
				contenido += strings.Split(string(ba.B_content[:]), "\x00")[0]
				contenido += "\"\n\t];"
				//conexiones inodo -> bloques
				contenido += "\n\tinode" + strconv.Itoa(p_inodo) + "->b" + strconv.Itoa(inb)
			}
		}
	}
	return contenido
}

// reporte del disco
func reporte_disco(id string, path string) {
	ls := buscarPath(id)
	var mbr MBR
	disco, err := os.OpenFile(ls, os.O_RDWR, 0660)
	if err != nil {
		panic(err)
	}
	disco.Seek(0, 0)
	rre := binary.Read(disco, binary.BigEndian, &mbr)
	if rre != nil {
		print(rre)
	}
	disco.Close()
	dotEx := ""
	colr := 0
	contenido := "digraph disco{ \n"
	contenido += "contenido [\n shape=plaintext \n label=< \n"
	contenido += "<table BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" color =\"blue\"> \n"
	contenido += "<tr><td>\n"
	contenido += "<table color=\"blue\" border =\"0\" cellborder=\"1\" cellpadding =\"10\"  cellspacing=\"0\">\n"
	contenido += "<tr><td rowspan =\"2\">MBR</td>\n"
	if mbr.Mbr_partition_1.Part_status[0] != '1' {
		contenido += "<td rowspan =\"2\">LIBRE<br/> 0%  del disco</td>\n"
	} else {
		if mbr.Mbr_partition_1.Part_type[0] == 'e' {
			var ebr EBR
			disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
			if err != nil {
				panic(err)
			}

			pstar, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_partition_1.Part_start[:]), "\x00")[0])
			disco.Seek(int64(pstar), 0)
			rre = binary.Read(disco, binary.BigEndian, &ebr)
			if rre != nil {
				print(rre)
			}
			disco.Close()
			if ebr.Part_status[0] == '0' {
				pnext, _ := strconv.Atoi(strings.Split(string(ebr.Part_next[:]), "\x00")[0])
				if pnext != -1 {
					disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
					if err != nil {
						panic(err)
					}
					disco.Seek(int64(pnext), 0)
					rre = binary.Read(disco, binary.BigEndian, &ebr)
					if rre != nil {
						print(rre)
					}
					disco.Close()
				}
			}
			// SE VA POR TODAS LAS PARTICIONES
			fin := false
			for !fin {
				dotEx += "<td>EBR</td>\n"
				eman := strings.Split(string(ebr.Part_name[:]), "\x00")[0]
				tt, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_partition_1.Part_size[:]), "\x00")[0])
				tp, _ := strconv.Atoi(strings.Split(string(ebr.Part_size[:]), "\x00")[0])
				dotEx += "<td>" + eman + "<br/>" + porcentaje(tt, tp) + "% de la particion</td>\n"
				colr++
				pnext, _ := strconv.Atoi(strings.Split(string(ebr.Part_next[:]), "\x00")[0])
				if pnext == -1 {
					fin = true
				} else {
					disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
					if err != nil {
						panic(err)
					}
					disco.Seek(int64(pnext), 0)
					rre = binary.Read(disco, binary.BigEndian, &ebr)
					if rre != nil {
						print(rre)
					}
					disco.Close()
				}

			}
			tt, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_tamano[:]), "\x00")[0])
			tp, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_partition_1.Part_size[:]), "\x00")[0])
			contenido += "<td colspan =\"" + strconv.Itoa(colr*2) + "\">EXTENDIDA<br/>" + porcentaje(tt, tp) + "% </td>\n"
		} else {
			eman := strings.Split(string(mbr.Mbr_partition_1.Part_name[:]), "\x00")[0]
			tt, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_tamano[:]), "\x00")[0])
			tp, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_partition_1.Part_size[:]), "\x00")[0])
			contenido += "<td rowspan =\"2\">" + eman + "<br/>" + porcentaje(tt, tp) + "% del disco</td>\n"
		}
	}

	if mbr.Mbr_partition_2.Part_status[0] != '1' {
		contenido += "<td rowspan =\"2\">LIBRE<br/>0% del disco</td>\n"
	} else {
		if mbr.Mbr_partition_2.Part_type[0] == 'e' {
			var ebr EBR
			disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
			if err != nil {
				panic(err)
			}

			pstar, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_partition_2.Part_start[:]), "\x00")[0])
			disco.Seek(int64(pstar), 0)
			rre = binary.Read(disco, binary.BigEndian, &ebr)
			if rre != nil {
				print(rre)
			}
			disco.Close()
			if ebr.Part_status[0] == '0' {
				pnext, _ := strconv.Atoi(strings.Split(string(ebr.Part_next[:]), "\x00")[0])
				if pnext != -1 {
					disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
					if err != nil {
						panic(err)
					}
					disco.Seek(int64(pnext), 0)
					rre = binary.Read(disco, binary.BigEndian, &ebr)
					if rre != nil {
						print(rre)
					}
					disco.Close()
				}
			}
			// SE VA POR TODAS LAS PARTICIONES
			fin := false
			for !fin {
				dotEx += "<td>EBR</td>\n"
				eman := strings.Split(string(ebr.Part_name[:]), "\x00")[0]
				tt, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_partition_2.Part_size[:]), "\x00")[0])
				tp, _ := strconv.Atoi(strings.Split(string(ebr.Part_size[:]), "\x00")[0])
				dotEx += "<td>" + eman + "<br/>" + porcentaje(tt, tp) + "% de la particion</td>\n"
				colr++
				pnext, _ := strconv.Atoi(strings.Split(string(ebr.Part_next[:]), "\x00")[0])
				if pnext == -1 {
					fin = true
				} else {
					disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
					if err != nil {
						panic(err)
					}
					disco.Seek(int64(pnext), 0)
					rre = binary.Read(disco, binary.BigEndian, &ebr)
					if rre != nil {
						print(rre)
					}
					disco.Close()
				}

			}
			tt, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_tamano[:]), "\x00")[0])
			tp, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_partition_2.Part_size[:]), "\x00")[0])
			contenido += "<td colspan =\"" + strconv.Itoa(colr*2) + "\">EXTENDIDA<br/>" + porcentaje(tt, tp) + "% </td>\n"
		} else {
			tt, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_tamano[:]), "\x00")[0])
			tp, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_partition_2.Part_size[:]), "\x00")[0])
			eman := strings.Split(string(mbr.Mbr_partition_2.Part_name[:]), "\x00")[0]
			contenido += "<td rowspan =\"2\">" + eman + "<br/>" + porcentaje(tt, tp) + "% del disco</td>\n"
		}
	}

	if mbr.Mbr_partition_3.Part_status[0] != '1' {
		contenido += "<td rowspan =\"2\">LIBRE<br/>0% del disco</td>\n"
	} else {
		if mbr.Mbr_partition_3.Part_type[0] == 'e' {
			var ebr EBR
			disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
			if err != nil {
				panic(err)
			}

			pstar, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_partition_3.Part_start[:]), "\x00")[0])
			disco.Seek(int64(pstar), 0)
			rre = binary.Read(disco, binary.BigEndian, &ebr)
			if rre != nil {
				print(rre)
			}
			disco.Close()
			if ebr.Part_status[0] == '0' {
				pnext, _ := strconv.Atoi(strings.Split(string(ebr.Part_next[:]), "\x00")[0])
				if pnext != -1 {
					disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
					if err != nil {
						panic(err)
					}
					disco.Seek(int64(pnext), 0)
					rre = binary.Read(disco, binary.BigEndian, &ebr)
					if rre != nil {
						print(rre)
					}
					disco.Close()
				}
			}
			// SE VA POR TODAS LAS PARTICIONES
			fin := false
			for !fin {
				dotEx += "<td>EBR</td>\n"
				eman := strings.Split(string(ebr.Part_name[:]), "\x00")[0]
				tt, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_partition_3.Part_size[:]), "\x00")[0])
				tp, _ := strconv.Atoi(strings.Split(string(ebr.Part_size[:]), "\x00")[0])
				dotEx += "<td>" + eman + "<br/>" + porcentaje(tt, tp) + "% de la particion</td>\n"
				colr++
				pnext, _ := strconv.Atoi(strings.Split(string(ebr.Part_next[:]), "\x00")[0])
				if pnext == -1 {
					fin = true
				} else {
					disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
					if err != nil {
						panic(err)
					}
					disco.Seek(int64(pnext), 0)
					rre = binary.Read(disco, binary.BigEndian, &ebr)
					if rre != nil {
						print(rre)
					}
					disco.Close()
				}

			}
			tt, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_tamano[:]), "\x00")[0])
			tp, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_partition_3.Part_size[:]), "\x00")[0])
			contenido += "<td colspan =\"" + strconv.Itoa(colr*2) + "\">EXTENDIDA<br/>" + porcentaje(tt, tp) + "% </td>\n"
		} else {
			tt, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_tamano[:]), "\x00")[0])
			tp, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_partition_3.Part_size[:]), "\x00")[0])
			eman := strings.Split(string(mbr.Mbr_partition_3.Part_name[:]), "\x00")[0]
			contenido += "<td rowspan =\"2\">" + eman + "<br/>" + porcentaje(tt, tp) + "% del disco</td>\n"
		}
	}

	if mbr.Mbr_partition_4.Part_status[0] != '1' {
		contenido += "<td rowspan =\"2\">LIBRE<br/>0% del disco</td>\n"
	} else {
		if mbr.Mbr_partition_4.Part_type[0] == 'e' {
			var ebr EBR
			disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
			if err != nil {
				panic(err)
			}

			pstar, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_partition_4.Part_start[:]), "\x00")[0])
			disco.Seek(int64(pstar), 0)
			rre = binary.Read(disco, binary.BigEndian, &ebr)
			if rre != nil {
				print(rre)
			}
			disco.Close()
			if ebr.Part_status[0] == '0' {
				pnext, _ := strconv.Atoi(strings.Split(string(ebr.Part_next[:]), "\x00")[0])
				if pnext != -1 {
					disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
					if err != nil {
						panic(err)
					}
					disco.Seek(int64(pnext), 0)
					rre = binary.Read(disco, binary.BigEndian, &ebr)
					if rre != nil {
						print(rre)
					}
					disco.Close()
				}
			}
			// SE VA POR TODAS LAS PARTICIONES
			fin := false
			for !fin {
				dotEx += "<td>EBR</td>\n"
				eman := strings.Split(string(ebr.Part_name[:]), "\x00")[0]
				tt, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_partition_4.Part_size[:]), "\x00")[0])
				tp, _ := strconv.Atoi(strings.Split(string(ebr.Part_size[:]), "\x00")[0])
				dotEx += "<td>" + eman + "<br/>" + porcentaje(tt, tp) + "% de la particion</td>\n"
				colr++
				pnext, _ := strconv.Atoi(strings.Split(string(ebr.Part_next[:]), "\x00")[0])
				if pnext == -1 {
					fin = true
				} else {
					disco, err = os.OpenFile(ls, os.O_RDWR, 0660)
					if err != nil {
						panic(err)
					}
					disco.Seek(int64(pnext), 0)
					rre = binary.Read(disco, binary.BigEndian, &ebr)
					if rre != nil {
						print(rre)
					}
					disco.Close()
				}

			}
			tt, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_tamano[:]), "\x00")[0])
			tp, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_partition_4.Part_size[:]), "\x00")[0])
			contenido += "<td colspan =\"" + strconv.Itoa(colr*2) + "\">EXTENDIDA<br/>" + porcentaje(tt, tp) + "% </td>\n"
		} else {
			tt, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_tamano[:]), "\x00")[0])
			tp, _ := strconv.Atoi(strings.Split(string(mbr.Mbr_partition_4.Part_size[:]), "\x00")[0])
			eman := strings.Split(string(mbr.Mbr_partition_4.Part_name[:]), "\x00")[0]
			contenido += "<td rowspan =\"2\">" + eman + "<br/>" + porcentaje(tt, tp) + "% del disco</td>\n"
		}
	}

	contenido += "</tr>\n"
	if dotEx != "" {
		contenido += "<tr>\n"
		contenido += dotEx
		contenido += "</tr>\n"
	}
	contenido += "</table>\n"
	contenido += "</td></tr>\n"
	contenido += "</table>\n"
	contenido += ">]\n }"

	crear_dot(contenido, path)
	//fmt.Println("SE GENERO EL REPORTE DEL DISCO " + id)
	respuesta += "\nSE GENERO EL REPORTE DEL DISCO " + id
}

// reporte del super bloque
func Reporte_SB(id string, path string) {
	ls := get_disco(id)
	var sb SUPER_BLOQUE
	archivo, err := os.OpenFile(ls.Path, os.O_RDWR, 0660)
	if err != nil {
		panic(err)
	}
	pps, _ := strconv.Atoi(strings.Split(string(ls.Part.Part_start[:]), "\x00")[0])
	archivo.Seek(int64(pps), 0)
	rre := binary.Read(archivo, binary.BigEndian, &sb)
	if rre != nil {
		print(rre)
	}
	archivo.Close()

	contenido := "digraph {\n"
	contenido += "node [ shape=none fontname=Helvetica]"
	contenido += "n1 [ label = <\n"
	contenido += "<table>\n"
	contenido += "<tr><td colspan=\"2\" bgcolor=\"#ccccff\">REPORTE SB</td></tr>\n"
	contenido += "<tr><td>s_filesystem_type</td><td>" + strings.Split(string(sb.S_filesystem_type[:]), "\x00")[0] + "</td></tr>\n"
	contenido += "<tr><td>s_inodes_count</td><td>" + strings.Split(string(sb.S_inodes_count[:]), "\x00")[0] + "</td></tr>\n"
	contenido += "<tr><td>s_blocks_count</td><td>" + strings.Split(string(sb.S_blocks_count[:]), "\x00")[0] + "</td></tr>\n"
	contenido += "<tr><td>s_free_blocks_count</td><td>" + strings.Split(string(sb.S_free_blocks_count[:]), "\x00")[0] + "</td></tr>\n"
	contenido += "<tr><td>s_free_inodes_count</td><td>" + strings.Split(string(sb.S_free_inodes_count[:]), "\x00")[0] + "</td></tr>\n"
	contenido += "<tr><td>s_mtime</td><td>" + strings.Split(string(sb.S_mtime[:]), "\x00")[0] + "</td></tr>\n"
	contenido += "<tr><td>s_mnt_count</td><td>" + strings.Split(string(sb.S_mnt_count[:]), "\x00")[0] + "</td></tr>\n"
	contenido += "<tr><td>s_magic</td><td>" + strings.Split(string(sb.S_magic[:]), "\x00")[0] + "</td></tr>\n"
	contenido += "<tr><td>s_inode_size</td><td>" + strings.Split(string(sb.S_inode_size[:]), "\x00")[0] + "</td></tr>\n"
	contenido += "<tr><td>s_block_size</td><td>" + strings.Split(string(sb.S_block_size[:]), "\x00")[0] + "</td></tr>\n"
	contenido += "<tr><td>s_first_ino</td><td>" + strings.Split(string(sb.S_first_ino[:]), "\x00")[0] + "</td></tr>\n"
	contenido += "<tr><td>s_first_blo</td><td>" + strings.Split(string(sb.S_first_blo[:]), "\x00")[0] + "</td></tr>\n"
	contenido += "<tr><td>s_bm_inode_start</td><td>" + strings.Split(string(sb.S_bm_inode_start[:]), "\x00")[0] + "</td></tr>\n"
	contenido += "<tr><td>s_bm_block_start</td><td>" + strings.Split(string(sb.S_bm_block_start[:]), "\x00")[0] + "</td></tr>\n"
	contenido += "<tr><td>s_inode_start</td><td>" + strings.Split(string(sb.S_inode_start[:]), "\x00")[0] + "</td></tr>\n"
	contenido += "<tr><td>s_block_start</td><td>" + strings.Split(string(sb.S_block_start[:]), "\x00")[0] + "</td></tr>\n"
	contenido += "</table>\n > ];\n  }"
	crear_dot(contenido, path)
	//fmt.Println("SE GENERO EL REPORTE DEL SB " + id)
	respuesta += "\nSE GENERO EL REPORTE DEL SB " + id
}

// Busca el path por el id
func buscarPath(id string) string {
	for i := 0; i < Tmontadas.Len(); i++ {
		mont := Tmontadas.GetValue(i).(Disco)
		if mont.Id == id {
			return mont.Path
		}
	}
	return ""
}

func crear_dot(contenido string, path string) {
	fmt.Println("----------------------------Dot------------------------------")
	fmt.Println(contenido)
	fmt.Println("----------------------------fin Dot------------------------------")
	directorio := directorio(path)
	name := nombre_archivo(path)
	ruta := directorio + "/" + name + ".dot"
	b := []byte(contenido)
	err := ioutil.WriteFile(ruta, b, 0755)
	if err != nil {
		panic(err)
	}
	err = exec.Command("dot", "-Tjpg", ruta, "-o", directorio+"/"+name+".jpg").Run()
	if err != nil {
		panic(err)
	}
}

// retorna el directorio de un path
func directorio(path string) string {
	ruta := strings.Split(path, "/")
	carpets := ""
	for i := 0; i < len(ruta)-1; i++ {
		carpets += ruta[i]
		carpets += "/"
	}
	carpets = carpets[0 : len(carpets)-1]
	return carpets
}

// retorna el nombre del archivo
func nombre_archivo(path string) string {
	ruta := strings.Split(path, "/")
	nomext := ruta[len(ruta)-1]
	sepNom := strings.Split(nomext, ".")
	nom := sepNom[0]
	return nom
}

// retorna la extension del archivo
func extension(path string) string {
	ruta := strings.Split(path, "/")
	nomext := ruta[len(ruta)-1]
	sepNom := strings.Split(nomext, ".")
	ext := sepNom[1]
	return ext
}

func porcentaje(tam_disco int, tam int) string {
	return strconv.Itoa((tam * 100) / tam_disco)
}

// busca el disco en la lista por el id
func get_disco(id string) Disco {
	var mdisco Disco
	for i := 0; i < Tmontadas.Len(); i++ {
		tempdisco := Tmontadas.GetValue(i).(Disco)
		if tempdisco.Id == id {
			return tempdisco
		}
	}
	return mdisco
}
