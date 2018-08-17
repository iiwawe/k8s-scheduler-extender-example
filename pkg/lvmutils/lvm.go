package lvmutils

import "github.com/google/lvmd/parser"

type LvmResource struct {
	Lvs []*parser.LV `json:"-"`
	Vg  *VG          `json:"vg"`
}

type VG struct {
	Name string   `json:"name"`
	UUID string   `json:"uuid"`
	Size uint64   `json:"size"`
	Free uint64   `json:"free"`
	Tags []string `json:"-"`
}
