package comandos

func Login(user string, pass string, id string) {
	if user == "" || pass == "" || id == "" {
		//fmt.Println("ERROR LOGIN: FALTAN PAREMETROS PARA REALIZAR ESTA ACCION")
		respuesta += "\nERROR LOGIN: FALTAN PAREMETROS PARA REALIZAR ESTA ACCION"
		return
	}
	if usuario_activo.Logeado {
		//fmt.Println("ERROR LOGIN: YA HAY UN USUARIO ACTIVO")
		respuesta += "\nERROR LOGIN: YA HAY UN USUARIO ACTIVO"
		return
	}
	var mont Disco
	existe := existe_montada(id, &mont)
	if !existe {
		//fmt.Println("ERROR LOGIN: NO SE ENCUENTRA LA PARTICION MONTADA CON ESE ID")
		respuesta += "\nERROR LOGIN: NO SE ENCUENTRA LA PARTICION MONTADA CON ESE ID"
		return
	}
	contenidoUser := leer_archivo(1, mont)
	registros := get_registros(contenidoUser)
	for _, registro := range registros { //UID [0], TIPO [1], GRUPO [2], USUARIO [3], CONTRASEÑA [4]
		if len(registro) == 5 && registro[3] == user && registro[4] == pass {
			usuario_activo.Nombre = user
			usuario_activo.Cont = pass
			usuario_activo.Gid = registro[2]
			usuario_activo.Uid = registro[0]
			usuario_activo.Montada = mont
			usuario_activo.Logeado = true
		}
	}
	if !usuario_activo.Logeado {
		//fmt.Println("ERROR LOGIN: NO SE PUDO INICIAR SESION")
		respuesta += "\nERROR LOGIN: NO SE PUDO INICIAR SESION"
		return
	}
	//fmt.Println("LOGIN: SE LOGEO " + usuario_activo.Nombre + " CON EXITO")
	respuesta += "\nLOGIN: SE LOGEO " + usuario_activo.Nombre + " CON EXITO"
}

func Logout() {
	if !usuario_activo.Logeado {
		//fmt.Println("ERROR LOGOUT: NO HAY USUARIO ACTIVO")
		respuesta += "\nERROR LOGOUT: NO HAY USUARIO ACTIVO"
		return
	}
	usuario_activo.Logeado = false
	usuario_activo.Uid = "-1"
	//fmt.Println("LOGOUT REALIZADO CON EXITO")
	respuesta += "\nLOGOUT REALIZADO CON EXITO"
}
