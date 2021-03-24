package internal

import (
	"fmt"
	"fuegoquasar/config"
	"log"
	"math"
)

func init() {

}

type punto struct {
	X, Y float64
}

type circulo struct {
	punto
	radio float64
}

func newCirculo(x, y, r float64) *circulo {
	return &circulo{punto{x, y}, r}
}

// una funcion que a partir de dos circulos busca los puntos en comun
// puede tener hasta dos puntos o no devolver ninguno
func buscaInterseccion(a *circulo, b *circulo) (p []punto) {
	dx, dy := b.X-a.X, b.Y-a.Y        // dx y dy son las distancias horizontal y vertical entre los centros
	Sr := a.radio + b.radio           // suma de radios
	Dr := math.Abs(a.radio - b.radio) // diferencia entre radios
	Dc := math.Sqrt(dx*dx + dy*dy)    // distancia entre ambos centros
	if Dc > Sr {
		fmt.Println("Los circulos no se intersectan")

	} else if Dc < Dr {
		fmt.Println("Un circulo esta contenido dentro de otro")
	} else if Dc == 0 && a.radio == b.radio {
		fmt.Println("Ambos circulos coinciden")
	} else if Dc <= Sr && Dc > Dr {

		aa := (a.radio*a.radio - b.radio*b.radio + Dc*Dc) / (2 * Dc)
		h := math.Sqrt(a.radio*a.radio - aa*aa)
		x2 := a.X + aa*(b.X-a.X)/Dc
		y2 := a.Y + aa*(b.Y-a.Y)/Dc
		x3 := math.Round(x2 + h*(b.Y-a.Y)/Dc)
		y3 := math.Round(y2 - h*(b.X-a.X)/Dc)
		p = append(p, punto{x3, y3})
		x4 := math.Round(x2 - h*(b.Y-a.Y)/Dc)
		y4 := math.Round(y2 + h*(b.X-a.X)/Dc)
		p = append(p, punto{x4, y4})

	}
	return
}

func GetLocation(d1, d2, d3 float32) (x, y float32, ok bool) {
	var x1, x2, x3, y1, y2, y3 float64
	var result []punto
	cfg, err := config.NewConfig(config.CfgPath)
	if err != nil {
		log.Fatal(err)
	}

	x1 = float64(cfg.RebelShip1.X)
	y1 = float64(cfg.RebelShip1.Y)
	x2 = float64(cfg.RebelShip2.X)
	y2 = float64(cfg.RebelShip2.Y)
	x3 = float64(cfg.RebelShip3.X)
	y3 = float64(cfg.RebelShip3.Y)

	satelite1 := newCirculo(x1, y1, float64(d1))
	satelite2 := newCirculo(x2, y2, float64(d2))
	satelite3 := newCirculo(x3, y3, float64(d3))
	intersec1 := buscaInterseccion(satelite1, satelite2)
	intersec2 := buscaInterseccion(satelite2, satelite3)
	intersec3 := buscaInterseccion(satelite1, satelite3)
	// una vez obtenidas las intersecciones, busco una que se repita en las tres
	// de ser asi, esa es la ubicación de la nave.
	for i := 0; i < len(intersec1); i++ {
		for j := 0; j < len(intersec2); j++ {
			if intersec1[i].X == intersec2[j].X && intersec1[i].Y == intersec2[j].Y {
				slice1 := intersec1[i : i+1]
				result = append(result, slice1...)
			}
		}
		for j := 0; j < len(intersec3); j++ {
			if intersec1[i].X == intersec3[j].X && intersec1[i].Y == intersec3[j].Y {
				slice2 := intersec2[i : i+1]
				result = append(result, slice2...)
			}
		}
	}
	if len(result) == 2 {
		if result[0].X == result[1].X && result[0].Y == result[1].Y {
			return float32(result[0].X), float32(result[0].Y), true
		} else {
			return 0, 0, false
		}
	} else {
		return 0, 0, false
	}

}

func GetMessage(message1, message2, message3 []string) (msg string, ok bool) {
	var minLen int // obtengo el mensaje mas corto
	var result bool = false
	var msgpart string
	if len(message1) <= len(message2) && len(message1) <= len(message3) {
		minLen = len(message1)
	} else if len(message2) <= len(message3) {
		minLen = len(message2)
	} else {
		minLen = len(message3)
	}
	realmsg1start := len(message1) - minLen // ajusto el primer campo acorde la minima longitud
	// para quitar el desfasaje
	realmsg2start := len(message2) - minLen
	realmsg3start := len(message3) - minLen
	for i := 0; i < minLen; i++ {
		if message1[i+realmsg1start] != "" {
			msgpart = message1[i+realmsg1start]

		} else if message2[i+realmsg2start] != "" {
			msgpart = message2[i+realmsg2start]

		} else {
			msgpart = message3[i+realmsg3start]
		}
		space := ""
		if i > 0 {
			space = " "
		}
		msg = msg + space + msgpart
	}
	if msg != "" {
		result = true
	}
	return msg, result
}
