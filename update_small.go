package main

import (
	"fmt"
	"os"

	"github.com/Centny/dbm/mgo"
	_ "github.com/Centny/ffcm"
	"github.com/Centny/gfs/gfsdb"
	"github.com/Centny/gwf/log"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("update_small <db connection>")
		return
	}
	mgo.AddDefault2(os.Args[1])
	gfsdb.C = mgo.C
	var fs, err = gfsdb.ListFv(bson.M{
		"info.small.text": bson.M{
			"$exists": 1,
		},
	})
	if err != nil {
		log.E("list file error(%v)", err)
		return
	}
	var skipped, updated int
	for _, f := range fs {
		var text = f.Info.StrValP("/small/text")
		if len(text) < 1 {
			skipped += 1
			continue
		}
		err = gfsdb.UpdateF(f.Id, bson.M{
			"info.small": bson.M{
				"count": 1,
				"files": []string{text},
			},
		})
		if err != nil {
			log.E("update file error(%v)", err)
			return
		}
		updated += 1
		if updated%100 == 1 {
			log.D("update file skipped(%v),updated(%v),total(%v) success", skipped, updated, len(fs))
		}
	}
	log.D("update all done by skipped(%v),updated(%v),total(%v)", skipped, updated, len(fs))
}
