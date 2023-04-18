package comandos

// ------------------------------ESTRUCTURAS---------------------------------------------
type MBR struct {
	mbr_tamano         [10]byte
	mbr_fecha_creacion [10]byte
	mbr_dsk_signature  [10]byte
	dsk_fit            [10]byte
	mbr_partition_1    PARTITION
	mbr_partition_2    PARTITION
	mbr_partition_3    PARTITION
	mbr_partition_4    PARTITION
}

type PARTITION struct {
	part_status [10]byte
	part_type   [10]byte
	part_fit    [10]byte
	part_start  [10]byte
	part_size   [10]byte
	part_name   [10]byte
}

type EBR struct {
	part_status [10]byte
	part_fit    [10]byte
	part_start  [10]byte
	part_size   [10]byte
	part_next   [10]byte
	part_name   [10]byte
}

func newMBR(t, f, d, ft string) MBR {
	e := MBR{}
	copy(e.mbr_tamano[:], t)
	copy(e.mbr_fecha_creacion[:], f)
	copy(e.mbr_dsk_signature[:], d)
	copy(e.dsk_fit[:], ft)
	return e
}
