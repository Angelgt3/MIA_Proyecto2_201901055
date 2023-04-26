package comandos

func Rmgrp(name string) {
	if name == "" {
		//fmt.Println("ERROR RMGRP: NO SE INGRESO EL NOMBRE")
		respuesta += "\nERROR RMGRP: NO SE INGRESO EL NOMBRE"
		return
	}
	if usuario_activo.Uid != "1" { //solo root
		//fmt.Println("ERROR RMGRP: SOLO EL USUARIO ROOT PUEDE REALIZAR ESTA OPERACION")
		respuesta += "\nERROR RMGRP: SOLO EL USUARIO ROOT PUEDE REALIZAR ESTA OPERACION"
		return
	}
	contenido := leer_archivo(1, usuario_activo.Montada)
	registros := get_registros(contenido)
	newcontenido := ""
	se_elimino := false
	for _, registro := range registros { //UID [0], TIPO [1], GRUPO [2], USUARIO [3], CONTRASEÑA [4]
		if len(registro) == 3 {
			if registro[2] == name {
				newcontenido += "\n0, " + registro[1] + ", " + registro[2]
				se_elimino = true
				continue
			}
			newcontenido += "\n" + registro[0] + ", " + registro[1] + ", " + registro[2]
		} else if len(registro) == 5 {
			newcontenido += "\n" + registro[0] + ", " + registro[1] + ", " + registro[2] + ", " + registro[3] + ", " + registro[4]
		}

	}
	eliminar_bloques_archivo(1, usuario_activo.Montada)
	escribir_bloques_archivo(1, newcontenido, usuario_activo.Montada)

	if se_elimino {
		//fmt.Println("EL GRUPO " + name + " FUE ELIMINADO CON EXITO")
		respuesta += "\nEL GRUPO " + name + " FUE ELIMINADO CON EXITO"
	} else {
		//fmt.Println("ERROR RMGRP: NO SE ENCONTRO EL GRUPO")
		respuesta += "\nERROR RMGRP: NO SE ENCONTRO EL GRUPO"
	}

}
