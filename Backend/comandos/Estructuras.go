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

type Ajuste struct {
	Inicio [6]int
	Fin    [6]int
	Inter  [6]int
	Resta  [6]int
}

type Tlist struct {
	Path string
	Name string
	Id   string
	Size int
	Part PARTITION
	N    int
}

// ------------------------------VARIABLES---------------------------------------------
var Tmontadas List

// ------------------------------FUNCIONES---------------------------------------------
func newAjuste() Ajuste {
	var l [6]int
	e := Ajuste{l, l, l, l}
	return e
}

func newEbr(st, f, s, tm, tms, n string) EBR {
	e := EBR{}
	copy(e.Part_status[:], st)
	copy(e.Part_fit[:], f)
	copy(e.Part_start[:], s)
	copy(e.Part_size[:], tm)
	copy(e.Part_next[:], tms)
	copy(e.Part_name[:], n)
	return e
}

func newMBR(t, f, d, ft string) MBR {
	e := MBR{}
	copy(e.Mbr_tamano[:], t)
	copy(e.Mbr_fecha_creacion[:], f)
	copy(e.Mbr_dsk_signature[:], d)
	copy(e.Dsk_fit[:], ft)
	return e
}

func newPartition(s, t, f, st, si, n []byte) PARTITION {
	e := PARTITION{}
	return e
}

func newTlist() Tlist {
	e := Tlist{}
	return e
}
