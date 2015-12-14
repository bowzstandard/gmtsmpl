package main

import(

	"time"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/gl"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/exp/sprite/glsprite"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	

	"gmtsmpl/textutil"
)

var (
	startTime = time.Now()
	images *glutil.Images
	eng sprite.Engine
	scene *sprite.Node
	sz size.Event
	txList =  [...]string{
		"話をしよう。",
		"あれは今から36万・・・いや、1万4000年前だったか。",
		"まぁいい。",
		"私にとってはつい昨日の出来事だが…君達にとっては多分明日の出来事だ。",
		"彼には72通りの名前があるから、なんて呼べばいいのか・・・。",
		"確か、最初に会った時は・・・Gopher！",
		"そう。あいつは最初から言うことを聞かなかった。",
		"私の言うとおりにしていればな。",
		"まあ・・・いい奴だったよ。",
	}
)

func main(){

	//entry point
	app.Main(func(a app.App){
		var glctx gl.Context
		//var sz size.Event

		for e:=range a.Events(){
			switch e:=a.Filter(e).(type){
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible){
				case lifecycle.CrossOn:
					
					glctx,_ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop()
					glctx = nil
				}
			case size.Event:
				sz = e
			case paint.Event:
				if glctx == nil || e.External{
					continue
				}
				onPaint(glctx,sz)
				//flush draw command
				a.Publish()
				a.Send(paint.Event{})
			}
		}
	})

}

func onStart(glctx gl.Context){
	images = glutil.NewImages(glctx)
	eng = glsprite.Engine(images)
	loadScene()
}

func onStop(){
	eng.Release()
	images.Release()
}

func onPaint(glctx gl.Context, sz size.Event) {
	glctx.ClearColor(1, 1, 1, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
	now := clock.Time(time.Since(startTime) * 60 / time.Second)
	eng.Render(scene, now, sz)
}

func newNode() *sprite.Node {
	n := &sprite.Node{}
	eng.Register(n)
	scene.AppendChild(n)
	return n
}

func loadScene(){
	scene = &sprite.Node{}
	eng.Register(scene)
	eng.SetTransform(scene,f32.Affine{
		{1,0,0},
		{0,1,0},
	})
	
	n := newNode()
	tf:=&textutil.TextField{}
	tf.New(eng,images)
	st:=make([]*textutil.TextureData,len(txList))
	for i:=0;i<len(txList);i++{
		st[i]=tf.Generate(sz,txList[i])
	}

	index:=0

	n.Arranger = arrangerFunc(func(eng sprite.Engine,n *sprite.Node,t clock.Time){

		t0 := uint32(t) % 150

		eng.SetSubTex(n,st[index].Get("body").(sprite.SubTex))
		
		eng.SetTransform(n,f32.Affine{
			{st[index].Get("width").(float32),0,0},
			{0,st[index].Get("height").(float32),20},
		})
		
		
		if t0 == 0{
			if index >= len(txList)-1{
				index = 0
			}else{
				index++
			}
		}
	})
}

type arrangerFunc func(e sprite.Engine, n *sprite.Node,t clock.Time)

func (a arrangerFunc) Arrange(e sprite.Engine,n *sprite.Node,t clock.Time){
	a(e,n,t)
}
