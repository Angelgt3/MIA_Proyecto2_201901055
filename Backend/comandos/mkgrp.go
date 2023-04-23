package comandos

import (
	"fmt"
	"strconv"
)

func Mkgrp(name string) {
	if name == "" {
		fmt.Println("ERROR MKGRP: NO SE INGRESO EL NOMBRE")
		return
	}
	if usuario_activo.Uid != "1" { //solo root
		fmt.Println("ERROR MKGRP: SOLO EL USUARIO ROOT PUEDE REALIZAR ESTA OPERACION")
		return
	}
	contenido := leer_archivo(1, usuario_activo.Montada)
	//verifica que no exista ese grupo
	registros := get_registros(contenido)
	cont := 1
	for _, registro := range registros { //UID [0], TIPO [1], GRUPO [2], USUARIO [3], CONTRASEÑA [4]
		if len(registro) == 3 {
			if registro[2] == name {
				fmt.Println("ERROR MKGRP: EL GRUPO " + name + " YA EXISTE")
				return
			}
			cont++
		}
	}
	contenido += "\n" + strconv.Itoa(cont) + ", G, " + name
	eliminar_bloques_archivo(1, usuario_activo.Montada)
	escribir_bloques_archivo(1, contenido, usuario_activo.Montada)
	fmt.Println("GRUPO " + name + " CREADO CON EXITO")
}
