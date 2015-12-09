package textutil

import(
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"github.com/golang/freetype/truetype"
	"image"
	"image/draw"
	"strings"
	"io/ioutil"

	"log"
)

const (
	suffix = ".ttf"
	default_font = "mplus-1p-light"
	default_fontsize = 24.0
	default_spacing = 1.5
	default_align = "left"
	dpi = 72.0
	scale = 4
	hinting = "none"
)

var (
	fontList = make(map[string]*truetype.Font)
)

type TextField struct{
	engine sprite.Engine
	size size.Event
	images *glutil.Images
	image *glutil.Image
	format *textFormat
}

func (tf *TextField)New(eng sprite.Engine,images *glutil.Images){
	tf.engine = eng
	tf.images = images
	tf.format = &textFormat{}
	tf.format.new()
	tf.setup()
}

func (tf *TextField)SetFormat(fm map[string]interface{}){
	tf.format.setFormat(fm)
	tf.setup()
}

func (tf *TextField)Generate(sz size.Event,text string)*TextureData{
	fz:=tf.format.get("fontsize").(float64)
	ft:=tf.format.get("font").(string)

	imgW, imgH := int(fz)*(len(text)+10),int(fz)+10*dpi/72

	tf.size = sz
	if tf.image != nil{
		tf.image.Release()
	}

	tf.image = tf.images.NewImage(imgW,imgH)
	
	fg,bg := image.Black,image.White
	draw.Draw(tf.image.RGBA,tf.image.RGBA.Bounds(),bg,image.Point{},draw.Src)

	// Draw the text.
	h := font.HintingNone
	switch hinting {
	case "full":
		h = font.HintingFull
	}

	d := &font.Drawer{
		Dst: tf.image.RGBA,
		Src: fg,
		Face: truetype.NewFace(fontList[ft], &truetype.Options{
			Size:    fz,
			DPI:     dpi,
			Hinting: h,
		}),
	}
	
	d.Dot = fixed.Point26_6{
		X:fixed.I(10),
		Y:fixed.I(int(fz*dpi/72)),
	}
	d.DrawString(text)

	tf.image.Upload()
	tf.image.Draw(
		sz,
		geom.Point{0, (sz.HeightPt - geom.Pt(imgH)/scale)},
		geom.Point{geom.Pt(imgW)/scale, (sz.HeightPt - geom.Pt(imgH)/scale)},
		geom.Point{0, (sz.HeightPt - geom.Pt(imgH)/scale)},
		tf.image.RGBA.Bounds().Inset(1),
	)

	t, err := tf.engine.LoadTexture(tf.image.RGBA)
	if err != nil {
		log.Fatal(err)
	}
	
	st := sprite.SubTex{t,tf.image.RGBA.Bounds().Inset(1)}
	
	td := &TextureData{}
	td.new(st,float32(imgH/scale),float32(imgW/scale))

	return td
}

func (tf *TextField)setup(){
	tf.loadFont()
}

func (tf *TextField)loadFont(){
	ft:=tf.format.get("font").(string)
	if _,ok:=fontList[ft];ok{
		return
	}

	a,err:=asset.Open(ft+suffix)
	if err != nil{
		log.Fatal(err)
	}

	fontBytes,err2:=ioutil.ReadAll(a)
	if err2 != nil{
		log.Fatal(err2)
	}

	tmp,err3 := truetype.Parse(fontBytes)
	if err3 != nil{
		log.Fatal(err3)
	}
	fontList[ft] = tmp

}

type TextureData struct{
	body sprite.SubTex
	height float32
	width float32
}

func (td *TextureData)new(st sprite.SubTex,h float32,w float32){
	td.body = st
	td.height = h
	td.width = w
}

func (td *TextureData)Get(key string)interface{}{
	var val interface{}

	switch key{
	case "body":
		val = td.body
	case "height":
		val = td.height
	case "width":
		val = td.width
	}

	return val
}

type textFormat struct{
	font string
	fontsize float64
	spacing float64
	align string
}

func (tm *textFormat)new(){
	tm.font = default_font
	tm.fontsize = default_fontsize
	tm.spacing = default_spacing
	tm.align = default_align
}

func (tm *textFormat)setFormat(fm map[string]interface{}){
	for key,val:=range fm{
		switch key{
		case "font":
			if v,ok := val.(string);ok{
				f := v
				if strings.HasSuffix(v,suffix){
					f = strings.TrimSuffix(v,suffix)
				}
				tm.font = f
			}
		case "fontsize":
			if v,ok := val.(float64);ok{
				tm.fontsize = v
			}
		case "spacing":
			if v,ok := val.(float64);ok{
				tm.spacing = v
			}
		case "align":
			if v,ok := val.(string);ok{
				tm.align = v
			}
		}	
	}
}

func (tm *textFormat)get(key string)interface{}{
	var val interface{}
	
	switch key{
	case "font":
		val = tm.font
	case "fontsize":
		val = tm.fontsize
	case "spacing":
		val = tm.spacing
	case "align":
		val = tm.align
	}

	return val

}
