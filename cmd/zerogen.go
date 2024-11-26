package main

import (
	"fmt"
	"os"

	"github.com/guonaihong/clop"
	"github.com/guonaihong/zerogen"
	"gopkg.in/yaml.v3"
)

func parseConfig(zeroGen *zerogen.ZeroGen) {

	all, err := os.ReadFile(zeroGen.ConfigFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	var zeroGenCnf zerogen.ZeroConfig
	err = yaml.Unmarshal(all, &zeroGenCnf)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, z := range zeroGenCnf.Local {
		if z.Table.Dsn == "" {
			z.Table.Dsn = zeroGenCnf.Global.Table.Dsn
		}
		if z.Table.ModelDir == "" {
			z.Table.ModelDir = zeroGenCnf.Global.Table.ModelDir
		}
		if z.Table.GoZeroApiDir == "" {
			z.Table.GoZeroApiDir = zeroGenCnf.Global.Table.GoZeroApiDir
		}
		if z.Table.CopyDir == "" {
			z.Table.CopyDir = zeroGenCnf.Global.Table.CopyDir
		}
		if z.Table.CrudLogicDir == "" {
			z.Table.CrudLogicDir = zeroGenCnf.Global.Table.CrudLogicDir
		}
		if z.Table.Table == "" {
			z.Table.Table = zeroGenCnf.Global.Table.Table
		}
		if z.Table.Home == "" {
			z.Table.Home = zeroGenCnf.Global.Table.Home
		}
		if z.Table.ModelPkgName == "" {
			z.Table.ModelPkgName = zeroGenCnf.Global.Table.ModelPkgName
		}
		if z.Table.ApiPrefix == "" {
			z.Table.ApiPrefix = zeroGenCnf.Global.Table.ApiPrefix
		}
		if z.Table.ServiceName == "" {
			z.Table.ServiceName = zeroGenCnf.Global.Table.ServiceName
		}
		if z.Table.ApiGroup == "" {
			z.Table.ApiGroup = zeroGenCnf.Global.Table.ApiGroup
		}
		if z.Table.ImportPathPrefix == "" {
			z.Table.ImportPathPrefix = zeroGenCnf.Global.Table.ImportPathPrefix
		}
		if z.Table.ApiUrlPrefix == "" {
			z.Table.ApiUrlPrefix = zeroGenCnf.Global.Table.ApiUrlPrefix
		}

		z.After.Copy = append(z.After.Copy, zeroGenCnf.Global.After.Copy...)
		z2 := zerogen.ZeroGen{ZeroGenCore: z.Table}

		fmt.Printf("%#v\n", z2)
		if err := z2.Run(); err != nil {
			fmt.Println(err)
		}
	}
}
func main() {
	var zeroGen zerogen.ZeroGen
	clop.Bind(&zeroGen)
	if len(zeroGen.ConfigFile) > 0 {
		parseConfig(&zeroGen)
		return
	}
	err := zeroGen.Run()
	if err != nil {
		fmt.Println(err)
	}
}
