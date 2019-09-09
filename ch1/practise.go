package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	practise_1_12()
}

//练习 1.1： 修改 echo 程序，使其能够打印 os.Args[0] ，即被执行命令本身的名字。
func practise_1_1() {
	var s, sep string
	for i := 0; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}

//练习 1.2： 修改 echo 程序，使其打印每个参数的索引和值，每个一行。
func practise_1_2() {
	for i := 0; i < len(os.Args); i++ {
		fmt.Println(i, " ", os.Args[i])
	}
}

//练习 1.3： 做实验测量潜在低效的版本和使用了 strings.Join 的版本的运行时间差异。
func practise_1_3() {
	ret := make([]string, len(os.Args))
	t1 := time.Now().UnixNano() / 1e9
	practise_1_1()
	t2 := time.Now().UnixNano() / 1e9
	fmt.Printf("not use strings.Join:%v;\n", t2-t1)
	t3 := time.Now().UnixNano() / 1e9
	for i := 0; i < len(os.Args); i++ {
		ret = append(ret, os.Args[i])
	}
	fmt.Println(strings.Join(ret, " "))
	t4 := time.Now().UnixNano() / 1e9
	fmt.Printf("use strings.Join:%v;\n", t4-t3)
}

//练习 练习 1.4： 修改 dup2 ，出现重复的行时打印文件名称。 ch1/practise.go ch1/dup2/main.go
type fileNameLine struct {
	content  string
	fileName string
}

func countLines(f *os.File, counts map[fileNameLine]int) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		var fnl fileNameLine
		fnl.content = input.Text()
		fnl.fileName = f.Name()
		counts[fnl]++
	}
	// NOTE: ignoring potential errors from input.Err()
}
func practise_1_4() {
	counts := make(map[fileNameLine]int)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}
	for fnl, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\t%s\n", n, fnl.content, fnl.fileName)
		}
	}
}

//练习 1.5： 修改前面的Lissajous程序里的调色板，由黑色改为绿色。我们可以
//用 color.RGBA{0xRR, 0xGG, 0xBB, 0xff} 来得到 #RRGGBB 这个色值，三个十六进制的字符串分
//别代表红、绿、蓝像素。
//练习1.6： 修改Lissajous程序，修改其调色板来生成更丰富的颜色，然后修改 SetColorIndex
//的第三个参数，看看显示结果吧。

var palette = []color.Color{color.White, color.RGBA{0xAA, 0xB4, 0x00, 0xff}}

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

func lissajous(out io.Writer) {
	const (
		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.9),
				blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
func practise_1_5_6() {
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) > 1 && os.Args[1] == "web" {
		//!+http
		handler := func(w http.ResponseWriter, r *http.Request) {
			lissajous(w)
		}
		http.HandleFunc("/", handler)
		//!-http
		log.Fatal(http.ListenAndServe("localhost:8000", nil))
		return
	}
	//!+main
	lissajous(os.Stdout)
}

//练习 1.7： 函数调用io.Copy(dst, src)会从src中读取内容，并将读到的结果写入到dst中，使用
//这个函数替代掉例子中的ioutil.ReadAll来拷贝响应结构体到os.Stdout，避免申请一个缓冲区
//（ 例子中的b） 来存储。记得处理io.Copy返回结果中的错误。
//练习 1.8： 修改fetch这个范例，如果输入的url参数没有 http:// 前缀的话，为这个url加上该
//前缀。你可能会用到strings.HasPrefix这个函数。
//练习 1.9： 修改fetch打印出HTTP协议的状态码，可以从resp.Status变量得到该状态码。
//parm:	www.baidu.com
func practise_1_7_8_9() {
	for _, url := range os.Args[1:] {
		if !strings.HasPrefix(url, "http://") {
			url = "http://" + url
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		//b, err := ioutil.ReadAll(resp.Body)
		//fmt.Printf("%s", b)
		_, err = io.Copy(os.Stdout, resp.Body)
		fmt.Printf("%s", resp.Status)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
	}
}

//练习 1.10： 找一个数据量比较大的网站，用本小节中的程序调研网站的缓存策略，对每个
//URL执行两遍请求，查看两次时间是否有较大的差别，并且每次获取到的响应内容是否一
//致，修改本节中的程序，将响应结果输出，以便于进行对比。
//练习 1.11： 在fatchall中尝试使用长一些的参数列表，比如使用在alexa.com的上百万网站里
//排名靠前的。如果一个网站没有回应，程序将采取怎样的行为？（ Section8.9 描述了在这种
//情况下的应对机制）
//参数：https://golang.org http://gopl.io http://gopl.io https://godoc.org
func practise_1_10_11() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch) // start a goroutine
	}
	for range os.Args[1:] {
		fmt.Println(<-ch) // receive from channel ch
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs  %7d  %s  %s", secs, nbytes, url, resp.Status)
}

//练习 1.12： 修改Lissajour服务，从URL读取变量，比如你可以访问 http://localhost:8000/?
//cycles=20 这个URL，这样访问可以将程序里的cycles默认的5修改为20。字符串转换为数字
//可以调用strconv.Atoi函数。你可以在godoc里查看strconv.Atoi的详细说明。
func lissajous_1_12(out io.Writer, cycle string) {
	cls, _ := strconv.Atoi(cycle)
	const (
		cycles  = 1     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < float64(cls)*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.9),
				blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
func practise_1_12() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var cycle string
		if err := r.ParseForm(); err != nil {
			log.Print(err)
		}
		for _, v := range r.Form {
			cycle = v[0]
		}
		lissajous_1_12(w, cycle)
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
