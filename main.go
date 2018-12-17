// main is main.
package main

import (
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
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	SVG()
}

func SVG() {
	http.HandleFunc("/svg", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "image/svg+xml")

		f, err := ioutil.ReadFile("test.html")
		if err != nil {
			panic(err)
		}
		w.Write(f)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func handler1(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	params := r.URL.Query()
	params.Add("ok", "mate")
	params.Add("?\\''tobe", "ornottobe")
	r.URL.RawQuery = params.Encode()
	log.Println("handler1", r.URL.RawQuery)

	http.Redirect(w, r, "/redirect?"+params.Encode(), http.StatusSeeOther)
}

func handler2(w http.ResponseWriter, r *http.Request) {
	s, _ := url.QueryUnescape(r.URL.Query().Encode())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}
func Gif() {
	http.HandleFunc("/gif", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			panic(err)
		}
		for k, v := range r.Form {
			if k == "cycles" {
				c, err := strconv.Atoi(v[0])
				if err != nil {
					panic(err)
				}
				cycles = c
			}
		}
		lissajous(w)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func getHttp() {
	httpPrefix := "http://www"
	url := os.Args[1]
	if !strings.HasPrefix(url, httpPrefix) {
		url = strings.Join([]string{httpPrefix, url}, ".")
	}
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	n, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	log.Println(resp.Status, n)
}

var palette = color.Palette{
	color.RGBA{0, 0, 0, 255},
	color.RGBA{255, 255, 255, 255},
	color.RGBA{255, 0, 0, 255},
	color.RGBA{0, 255, 0, 255},
	color.RGBA{0, 0, 255, 255},
}

var (
	cycles  = 5 // number of complete x oscillator revolutions
	randSrc = rand.New(rand.NewSource(time.Now().UnixNano()))
)

const (
	res     = 0.001 // angular resolution
	size    = 100   // image canvas covers [-size..+size]
	nframes = 120   // number of animation frames
	delay   = 1     // delay between frames in 10ms units
)

func lissajous(out io.Writer) {

	freq := randSrc.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, int(2*size+1), int(2*size+1))
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				uint8(randSrc.Intn(4)+1))
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	err := gif.EncodeAll(out, &anim)
	if err != nil {
		panic(err)
	}
}
