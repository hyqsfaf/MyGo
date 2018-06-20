package main

import (
	"strings"
	"fmt"
	"os"
	"bufio"
	"io"
	"time"
)

type Reader interface {
	Read(rc chan string)

}

type Writer interface {
	Write(wc chan string)
}


type LogProcess struct {
	rc chan string
	wc chan string

	read Reader
	write Writer

}

type ReadFromFile struct {
	path string //读取文件路径
}

type WriteToInfluxDB struct {
	influxDBDsn string
}

//读取模块
func (r *ReadFromFile)Read(rc chan string){
	f,err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file error :%s",err.Error()))
	}


	f.Seek(0,2)
	rd := bufio.NewReader(f)
	for {
		line,err := rd.ReadBytes('\n')
		//fmt.Print("1\n")
		if err == io.EOF {
			//fmt.Print("2\n")
			time.Sleep(500*time.Microsecond)
			continue
		}else if err != nil {
			fmt.Print("3")
			panic(fmt.Sprintf("ReadBytes error :%s",err.Error()))

		}

		rc <- string(line[:len(line)-1])
	}

}

//写入模块
func (w *WriteToInfluxDB)Write(wc chan string)  {
	for v := range wc {
		fmt.Println(v)
	}
}


//解析模块
func (l *LogProcess) Process()  {

	for v := range l.rc {
		l.wc <- strings.ToUpper(v)
	}
}




func main()  {
	r := &ReadFromFile{
		path:"./access.log",
	}

	w := &WriteToInfluxDB{
		influxDBDsn:"username..",
	}

	lp := &LogProcess{
		rc:make(chan string),
		wc:make(chan string),
		write:w,
		read:r,
	}

	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)

	var i int64
	for {
		i++
	}
}
