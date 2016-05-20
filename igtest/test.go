package main

import (
	"fmt"
	"github.com/Centny/dbm/mgo"
	"github.com/Centny/gfs/gfsapi"
	"github.com/Centny/gfs/gfsdb"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/tutil"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2/bson"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

func usage() {
	fmt.Println(`Uage: test -con <database name> -srv http://xxx <other options> <path>
	-srv <the file server root path>
	-args <the extern server arguments>
	-reg <the file name regex, default is .*>
	-count <the total run count, default is 100>
	-tc <max count run at the same time>
	-logf <the log file>
	-tmp	<the tmp folder, default is ./tmp>`)
}

func main() {
	runtime.GOMAXPROCS(util.CPU())
	_, args, paths := util.Args()
	if len(paths) < 1 {
		usage()
		os.Exit(1)
		return
	}
	var con, srv, sargs, reg string
	var count, tc int = 100, 3
	var logf, tmp string = "", "tmp"
	err := args.ValidF(`
		con,R|S,L:0;
		srv,R|S,L:0;
		args,O|S,L:0;
		reg,O|S,L:0;
		count,O|I,R:0;
		tc,O|I,R:0;
		logf,O|S,L:0;
		tmp,O|S,L:0;
		`, &con, &srv, &sargs, &reg, &count, &tc, &logf, &tmp)
	if err != nil {
		fmt.Println(err)
		usage()
		os.Exit(1)
		return
	}
	var name = ""
	cons := strings.Split(con, "/")
	if len(cons) < 2 {
		fmt.Println("invalid connection")
		os.Exit(1)
		return
	}
	name = cons[1]
	gfsapi.SrvAddr = func() string {
		return srv
	}
	gfsapi.SrvArgs = func() string {
		return sargs
	}
	err = mgo.AddDefault(con, name)
	if err != nil {
		log.E("connection to database error->%v", err)
		os.Exit(1)
		return
	}
	gfsdb.C = mgo.C
	files := map[string][]string{}
	shas := []string{}
	max := 0
	util.ListFunc(paths[0], reg, func(t string) string {
		ss, err := os.Stat(t)
		if err != nil {
			return t
		}
		if ss.IsDir() {
			return t
		}
		key := filepath.Ext(t)
		files[key] = append(files[key], t)
		sha, _ := util.Sha1(t)
		shas = append(shas, sha)
		if len(files[key]) > max {
			max = len(files[key])
		}
		fmt.Println(sha, "->", t)
		return t
	})
	if len(files) < 1 {
		log.E("file not found in path %v", paths[0])
		os.Exit(1)
		return
	}
	for i := 0; i < count; i++ {
		fmt.Printf("\n\n== Testing %v/%v...\n", i, count)
		err = run_all(shas, files, max, tc, logf, tmp)
		if err != nil {
			log.E("error->%v", err)
			os.Exit(1)
			return
		}
		fmt.Printf("== Test %v/%v done...\n\n", i, count)
	}
	log.D("all done...")
}

func run_all(shas []string, files map[string][]string, max, tc int, logf, tmp string) error {
	err := clear_shas(shas)
	if err != nil {
		return err
	}
	err = check_status_beg()
	if err != nil {
		return err
	}
	_, err = tutil.DoPerfV_(max, max, logf, func(i int) error {
		err := run(files, i, tmp)
		if err != nil {
			panic(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = check_status_beg()
	if err != nil {
		return err
	}
	return nil
}

func clear_shas(shas []string) error {
	log.D("do clear file by shas(%v)", shas)
	fs, err := gfsdb.ListShaF(shas)
	if err != nil {
		err = util.Err("list file by sha(%v) error->%v", shas, err)
		return err
	}
	fids := []string{}
	for _, f := range fs {
		fids = append(fids, f.Id)
	}
	log.D("do clean file by fids(%v)", fids)
	_, err = mgo.C(gfsdb.CN_F).RemoveAll(bson.M{"_id": bson.M{"$in": fids}})
	if err != nil {
		err = util.Err("remvoe all file by fid(%v)", fids)
		return err
	}
	_, err = mgo.C("ffcm_task").RemoveAll(bson.M{"_id": bson.M{"$in": fids}})
	if err != nil {
		err = util.Err("remvoe all ffcm_task by fid(%v)", fids)
		return err
	}
	return nil
}

func check_status_beg() error {
	res, err := gfsapi.DoAdmStatus()
	if err != nil {
		err = util.Err("do get adm/status error->%v", err)
		log.E("%v", err)
		return err
	}
	if res.IntValP("/total") != 0 || res.IntValP("/running_c") != 0 {
		err = util.Err("check status fail with task is not zero->%v", util.S2Json(res))
		return err
	}
	if len(res.MapValP("/proc/task_c")) < 1 {
		err = util.Err("check status fail with task_c is zero->%v", util.S2Json(res))
		return err
	}
	return nil
}

var lck = sync.RWMutex{}
var v_reg = regexp.MustCompile("(?i)^.*\\.(wmv|rm|rmvb|mpg|mpeg|mpe|3gp|mov|mp4|m4v|avi|mkv|flv|vob)$")
var v_doc_reg = regexp.MustCompile("(?i)^.*\\.(doc|docx|xps|rtf)$")
var v_pdf_reg = regexp.MustCompile("(?i)^.*\\.(pdf)$")
var v_ppt_reg = regexp.MustCompile("(?i)^.*\\.(ppt|pptx)$")
var v_img_reg = regexp.MustCompile("(?i)^.*\\.(jpg|jpeg|png|bmp)$")

func run(files map[string][]string, idx int, tmp string) error {
	marks := map[string]string{}
	shas := map[string]string{}
	pubs := map[string]string{}
	//
	log.D("test upload file...")
	for ext, fs := range files {
		if len(fs) <= idx {
			continue
		}
		sha, err := util.Sha1(fs[idx])
		if err != nil {
			err = util.Err("read sha by path(%v) error->%v", fs[idx], err)
			log.E("%v", err)
			return err
		}
		shas[ext] = sha
		mark := util.UUID()
		marks[ext] = mark
		log.D("do upload file %v", fs[idx])
		res, err := gfsapi.DoUpF(fs[idx], "", mark, "", "", "", 1, 0)
		if err != nil {
			err = util.Err("do file upload error->%v", err)
			log.E("%v", err)
			return err
		}
		if res.StrValP("/base/sha") != sha {
			err = util.Err("response sha error, expect %v, but %v found in result->%v", sha, res.StrValP("/base/sha"), util.S2Json(res))
			log.E("%v", err)
			return err
		}
		pubs[ext] = res.StrVal("data")
		//
		if v_reg.MatchString(ext) || v_doc_reg.MatchString(ext) ||
			v_pdf_reg.MatchString(ext) || v_ppt_reg.MatchString(ext) ||
			v_img_reg.MatchString(ext) {
			if res.StrValP("/base/exec") != "running" && res.StrValP("/base/exec") != "done" {
				err = util.Err("response exec status is %v, %v expect in result->%v", res.StrValP("/base/exec"), "running", util.S2Json(res))
				log.E("%v", err)
				return err
			}
		}
	}
	log.D("test upload %v file success", len(pubs))
	//
	log.D("test donwload base file by \npubs->%v\nshas->%v", util.S2Json(pubs), util.S2Json(shas))
	for ext, pub := range pubs {
		tf := tmp + "/" + util.UUID()
		err := util.DLoad(tf, "%v", pub)
		if err != nil {
			err = util.Err("do file download by pub(%v) error->%v", pub, err)
			log.E("%v", err)
			return err
		}
		sha, err := util.Sha1(tf)
		if err != nil {
			err = util.Err("read sha by path(%v) error->%v", tf, err)
			log.E("%v", err)
			return err
		}
		if sha != shas[ext] {
			err = util.Err("download file sha error, expect %v, but %v found, the file is (%v,%v)", shas[ext], sha, pub, tf)
			log.E("%v", err)
			return err
		}
	}
	log.D("test download base file success")
	//
	err := wait_v(marks, pubs, tmp)
	if err != nil {
		return err
	}
	err = wait_d(marks, pubs, v_doc_reg, "D_docx", tmp)
	if err != nil {
		return err
	}
	err = wait_d(marks, pubs, v_pdf_reg, "D_pdfx", tmp)
	if err != nil {
		return err
	}
	err = wait_d(marks, pubs, v_ppt_reg, "D_pptx", tmp)
	if err != nil {
		return err
	}
	return nil
}

func wait_v(marks, pubs map[string]string, tmp string) error {
	log.D("test wait video....")
	exts := []string{}
	for ext, _ := range marks {
		if v_reg.MatchString(ext) {
			exts = append(exts, ext)
		}
	}
	if len(exts) < 1 {
		log.D("test wait video done with not video task found")
		return nil
	}
	for i := 0; i < len(exts); {
		log.D("waiting %v/%v done", exts[i], pubs[exts[i]])
		res, err := gfsapi.DoInfo("", "", "", marks[exts[i]], "")
		if err != nil {
			err = util.Err("do get file info by mark(%v) error->%v", marks[exts[i]], err)
			log.E("%v", err)
			return err
		}
		log.D("%v/%v info is->\n%v", exts[i], pubs[exts[i]], util.S2Json(res))
		if len(res.StrValP("/base/info/V_pc/text")) < 1 || len(res.StrValP("/base/info/V_phone/text")) < 1 {
			time.Sleep(2 * time.Second)
			continue
		}
		tf := tmp + "/" + util.UUID()
		err = util.DLoad(tf, "%v/V_pc", pubs[exts[i]])
		if err != nil {
			err = util.Err("download %v/V_pc error->%v", pubs[exts[i]], err)
			return err
		}
		tf = tmp + "/" + util.UUID()
		err = util.DLoad(tf, "%v/V_phone", pubs[exts[i]])
		if err != nil {
			err = util.Err("download %v/V_pc error->%v", pubs[exts[i]], err)
			return err
		}
		log.D("%v/%v done success", exts[i], pubs[exts[i]])
		i++
	}
	log.D("test wait video done with %v success....", len(exts))
	return nil
}

func wait_d(marks, pubs map[string]string, reg *regexp.Regexp, key, tmp string) error {
	log.D("test wait %v....", key)
	exts := []string{}
	for ext, _ := range marks {
		if reg.MatchString(ext) {
			exts = append(exts, ext)
		}
	}
	if len(exts) < 1 {
		log.D("test wait %v done with not task found", key)
		return nil
	}
	for i := 0; i < len(exts); {
		log.D("waiting %v/%v done", exts[i], pubs[exts[i]])
		res, err := gfsapi.DoInfo("", "", "", marks[exts[i]], "")
		if err != nil {
			err = util.Err("do get file info by mark(%v) error->%v", marks[exts[i]], err)
			log.E("%v", err)
			return err
		}
		log.D("%v/%v info is->%v", exts[i], pubs[exts[i]], util.S2Json(res))
		count := int(res.IntValP("/base/info/" + key + "/count"))
		if count < 1 {
			time.Sleep(2 * time.Second)
			continue
		}
		for j := 0; j < count; j++ {
			tf := tmp + "/" + util.UUID()
			err = util.DLoad(tf, "%v/%v/%v", pubs[exts[i]], key, j)
			if err != nil {
				err = util.Err("download %v/%v/%v error->%v", pubs[exts[i]], key, j, err)
				return err
			}
		}
		log.D("%v/%v done success", exts[i], pubs[exts[i]])
		i++
	}
	log.D("test wait %v done with %v success....", key, len(exts))
	return nil
}
