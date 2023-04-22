package comandos

// ------------------------------ESTRUCTURAS---------------------------------------------
type MBR struct {
	Mbr_tamano         [10]byte  //tamaño del disco
	Mbr_fecha_creacion [10]byte  //fecha de creacion disco
	Mbr_dsk_signature  [10]byte  //id del disco
	Dsk_fit            [10]byte  //ajuste de particion
	Mbr_partition_1    PARTITION //particiones
	Mbr_partition_2    PARTITION
	Mbr_partition_3    PARTITION
	Mbr_partition_4    PARTITION
}

type PARTITION struct {
	Part_status [10]byte //0=libre , 1=ocupada
	Part_type   [10]byte // p=primaria e=extendida
	Part_fit    [10]byte // b=best f=first w=worst
	Part_start  [10]byte //inicio de la particion
	Part_size   [10]byte //tamaño total de la particion
	Part_name   [10]byte //nombre de la particion
}

type EBR struct {
	Part_status [10]byte // 0=libre, 1=activa
	Part_fit    [10]byte // b=best f=first w=worst
	Part_start  [10]byte //inicio de la particion
	Part_size   [10]byte //tamaño total de la particion
	Part_next   [10]byte //byte del proximo ebr
	Part_name   [10]byte //nombre de la particion
}

type Ajuste struct { //estructura para controlar el ajuste de las particiones
	Inicio [6]int
	Fin    [6]int
	Inter  [6]int
	Resta  [6]int
}

type Disco struct { //Mantiene toda la info de una particion montada
	Path string
	Name string
	Id   string
	Size int
	Part PARTITION
	N    int
}

type SUPER_BLOQUE struct {
	S_filesystem_type   [10]byte //Guarda el número que identifica el sistema de archivos utilizado
	S_inodes_count      [10]byte //Guarda el número total de inodos
	S_blocks_count      [10]byte //Guarda el número total de bloques
	S_free_blocks_count [10]byte //Contiene el número de bloques libres
	S_free_inodes_count [10]byte //Contiene el número de inodos libres
	S_mtime             [10]byte //Última fecha en el que el sistema fue montado
	S_mnt_count         [10]byte //ndica cuantas veces se ha montado el sistema
	S_magic             [10]byte //Valor que identifica al sistema de archivos,
	S_inode_size        [10]byte //Tamaño del inodo
	S_block_size        [10]byte //Tamaño del bloque
	S_first_ino         [10]byte //Primer inodo libre
	S_first_blo         [10]byte //Primer bloque libre
	S_bm_inode_start    [10]byte //Guardará el inicio del bitmap de inodos
	S_bm_block_start    [10]byte //Guardará el inicio del bitmap de bloques
	S_inode_start       [10]byte //Guardará el inicio de la tabla de inodos
	S_block_start       [10]byte //Guardará el inicio de la tabla de bloques
}

type TINODOS struct {
	I_uid   [4]byte  //UID del usuario propietario del archivo o carpeta
	I_gid   [4]byte  //GID del grupo al que pertenece el archivo o carpeta.
	I_size  [4]byte  //Tamaño del archivo en bytes
	I_atime [10]byte //Última fecha en que se leyó el inodo sin modificarlo
	I_ctime [10]byte //Fecha en la que se creó el inodo
	I_mtime [10]byte //Última fecha en la que se modifica el inodo
	I_block [16]byte //Array en los que los primeros 16 registros son bloques directos.
	I_type  [1]byte  //Indica si es archivo o carpeta. 1 = ARCHIVO, 0=CARPETA
	I_perm  [4]byte  //permisos
}

type CONTENIDO struct {
	B_name  [12]byte // Nombre de carpeta o archivo
	B_inodo [4]byte  // Apuntador hacia un inodo asociado al archivo o carpeta
}

type BLOQUE_CARPETA struct {
	B_content [4]CONTENIDO //Array con el contenido de la carpeta
}

type BLOQUE_ARCHIVO struct {
	B_content [64]byte // array con el contenido del archivo
}
