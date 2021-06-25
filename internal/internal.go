package internal

import (
	"locator/config"
	"log"
	"math"
)

var cfg *config.Config
var err error

func init() {
	cfg, err = config.NewConfig(config.CfgPath)
	if err != nil {
		log.Fatal(err)
	}
}

type point struct {
	X, Y float64
}

type circle struct {
	point
	radio float64
}

func newCircle(x, y, r float64) *circle {
	return &circle{point{x, y}, r}
}

func lookIntersection(a *circle, b *circle) (p []point) {
	dx, dy := b.X-a.X, b.Y-a.Y
	Sr := a.radio + b.radio
	Dr := math.Abs(a.radio - b.radio)
	Dc := math.Sqrt(dx*dx + dy*dy)
	if Dc > Sr {
		log.Println("Circles do not intersect")
	} else if Dc < Dr {
		log.Println("One circle is inside the other")
	} else if Dc == 0 && a.radio == b.radio {
		log.Println("Both are the same circle")
	} else if Dc <= Sr && Dc > Dr {

		aa := (a.radio*a.radio - b.radio*b.radio + Dc*Dc) / (2 * Dc)
		h := math.Sqrt(a.radio*a.radio - aa*aa)
		x2 := a.X + aa*(b.X-a.X)/Dc
		y2 := a.Y + aa*(b.Y-a.Y)/Dc
		x3 := math.Round((x2+h*(b.Y-a.Y)/Dc)*10) / 10
		y3 := math.Round((y2-h*(b.X-a.X)/Dc)*10) / 10
		p = append(p, point{x3, y3})
		x4 := math.Round((x2-h*(b.Y-a.Y)/Dc)*10) / 10
		y4 := math.Round((y2+h*(b.X-a.X)/Dc)*10) / 10
		p = append(p, point{x4, y4})
	}
	return
}

func GetLocation(distances ...float32) (x, y float32) {
	var result []point
	var circles []*circle
	circles = make([]*circle, len(distances))
	for i := 0; i < len(cfg.Ships) && i < len(distances); i++ {
		circles[i] = newCircle(float64(cfg.Ships[i].X), float64(cfg.Ships[i].Y), float64(distances[i]))
	}

	intersec1 := lookIntersection(circles[0], circles[1])
	intersec2 := lookIntersection(circles[1], circles[2])
	intersec3 := lookIntersection(circles[0], circles[2])

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
			return float32(result[0].X), float32(result[0].Y)
		} else {
			return -0.09, -0.09 // -0.09 is defaulted as not found value
		}
	} else {
		return -0.09, -0.09
	}
}

func GetMessage(messages ...[]string) (msg string) {
	var minLen int
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
		realstart := len(messages[i]) - minLen
		for j := realstart; j < len(messages[i]); j++ {
			if messages[i][j] != "" && completemsg[j-realstart] == "" {
				completemsg[j-realstart] = messages[i][j]
			}
		}
	}
	for k := 0; k < len(completemsg); k++ {
		if completemsg[k] != "" {
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
		log.Println("Message cannot be recovered. At least one of array elements is empty")
		msg = ""
	}
	return msg
}
