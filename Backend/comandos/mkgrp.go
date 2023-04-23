package comandos

import "fmt"

func Mkgrp(name string) {
	if name == "" {
		fmt.Println("ERROR MKGRP: NO SE INGRESO EL NOMBRE")
		return
	}
	contenido := leer_archivo(1, usuario_activo.Montada)
	contenido += "\n1, G, " + name
	eliminar_bloques_archivo(1, usuario_activo.Montada)
	escribir_bloques_archivo(1, contenido, usuario_activo.Montada)
	fmt.Println("GRUPO " + name + " CREADO CON EXITO")
}
