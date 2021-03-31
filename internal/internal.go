package internal

import (
	"fuegoquasar/config"
	"log"
	"math"
)

var cfg *config.Config
var err error

func init() { //cargo el archivo de configuracion
	cfg, err = config.NewConfig(config.CfgPath)
	if err != nil {
		log.Fatal(err)
	}
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
		log.Println("Los circulos no se intersectan")
	} else if Dc < Dr {
		log.Println("Un circulo esta contenido dentro de otro")
	} else if Dc == 0 && a.radio == b.radio {
		log.Println("Ambos circulos coinciden")
	} else if Dc <= Sr && Dc > Dr {

		aa := (a.radio*a.radio - b.radio*b.radio + Dc*Dc) / (2 * Dc)
		h := math.Sqrt(a.radio*a.radio - aa*aa)
		x2 := a.X + aa*(b.X-a.X)/Dc
		y2 := a.Y + aa*(b.Y-a.Y)/Dc
		x3 := math.Round(x2 + h*(b.Y-a.Y)/Dc)
		y3 := math.Round(y2 - h*(b.X-a.X)/Dc)
		p = append(p, punto{x3, y3})
		x4 := math.Round((x2-h*(b.Y-a.Y)/Dc)*10) / 10 //redondeo a 1 cifra decimal
		y4 := math.Round((y2+h*(b.X-a.X)/Dc)*10) / 10
		p = append(p, punto{x4, y4})
	}
	return
}

func GetLocation(distances ...float32) (x, y float32) {
	var result []punto
	var satelites []*circulo
	satelites = make([]*circulo, len(distances))
	for i := 0; i < len(cfg.RebelShips) && i < len(distances); i++ {
		satelites[i] = newCirculo(float64(cfg.RebelShips[i].X), float64(cfg.RebelShips[i].Y), float64(distances[i]))
	}

	intersec1 := buscaInterseccion(satelites[0], satelites[1])
	intersec2 := buscaInterseccion(satelites[1], satelites[2])
	intersec3 := buscaInterseccion(satelites[0], satelites[2])
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
	if len(result) == 2 { // el punto coincide en los tres circulos
		if result[0].X == result[1].X && result[0].Y == result[1].Y {
			return float32(result[0].X), float32(result[0].Y)
		} else {
			return -0.09, -0.09 //defini -0.09 como un valor para "no encontrado", asumiendo tambien que todas las respuestas, tienen maximo 1 decimal
		}
	} else {
		return -0.09, -0.09
	}
}

func GetMessage(messages ...[]string) (msg string) {
	var minLen int // obtengo el mensaje mas corto
	var result bool = true
	var completemsg []string
	minLen = len(messages[0])
	for i := 0; i < len(messages); i++ {
		if len(messages[i]) < minLen {
			minLen = len(messages[i])
		}
	}
	completemsg = make([]string, minLen)
	for i := 0; i < len(messages); i++ {
		realstart := len(messages[i]) - minLen // para quitar el desfasaje
		for j := realstart; j < len(messages[i]); j++ {
			if messages[i][j] != "" && completemsg[j-realstart] == "" {
				completemsg[j-realstart] = messages[i][j] //asigno los elementos no vacios
			}
		}
	}
	for k := 0; k < len(completemsg); k++ {
		if completemsg[k] != "" { //valido si complete el mensaje
			space := ""
			if k > 0 {
				space = " "
			}
			msg = msg + space + completemsg[k]
		} else {
			result = false
			break
		}
	}
	if result == false {
		log.Println("No se pudo determinar el mensaje. Uno de los elementos esta vacio")
		msg = ""
	}
	return msg
}
