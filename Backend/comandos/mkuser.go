package comandos

import (
	"fmt"
	"strconv"
)

func Mkuser(user string, pwd string, grp string) {
	if user == "" || pwd == "" || grp == "" {
		fmt.Println("ERROR MKUSER: NO SE INGRESO TODOS LOS PARAMETROS OBLIGATORIOS")
		return
	}
	if usuario_activo.Uid != "1" { //solo root
		fmt.Println("ERROR MKUSER: SOLO EL USUARIO ROOT PUEDE REALIZAR ESTA OPERACION")
		return
	}
	contenido := leer_archivo(1, usuario_activo.Montada)
	//verifica que no exista ese grupo
	registros := get_registros(contenido)
	cont := 1
	for _, registro := range registros { //UID [0], TIPO [1], GRUPO [2], USUARIO [3], CONTRASEÑA [4]
		if len(registro) == 5 {
			if registro[3] == user {
				fmt.Println("ERROR MKUSER: EL GRUPO " + user + " YA EXISTE")
				return
			}
			cont++
		}
	}
	contenido += "\n" + strconv.Itoa(cont) + ", U, " + grp + ", " + user + ", " + pwd
	eliminar_bloques_archivo(1, usuario_activo.Montada)
	escribir_bloques_archivo(1, contenido, usuario_activo.Montada)
	fmt.Println("USUARIO " + user + " CREADO CON EXITO")
}
