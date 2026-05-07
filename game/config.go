package game

import (
	"time"
)

// step = BaseInterval - (len-1)*PerSegment, нижняя граница MinInterval.
const (
	BaseInterval = time.Second / 6
	MinInterval  = 70 * time.Millisecond
	PerSegment   = 5 * time.Millisecond
)

const (
	GridWidth       = 40
	GridHeight      = 30
	DefaultCellSize = 20
)

const WindowMonitorFraction = 0.75

const (
	ScorePerApple = 1
)

const (
	HudPadding     = 10
	HudFontMinSize = 12
	HudFontMaxSize = 28
)

// Параметры анимации счёта в HUD.
const (
	// ScoreAnimRate — скорость сближения DisplayScore со Score (1/сек).
	// Больше — быстрее догоняет; 8 даёт мягкое «докручивание» ~0.3–0.5c.
	ScoreAnimRate = 8.0
	// ScorePulseDuration — длительность «пульса» масштаба при инкременте.
	ScorePulseDuration = 280 * time.Millisecond
	// ScorePulseScale — максимальный коэффициент увеличения в пике пульса.
	ScorePulseScale = 1.35
)

var (
	ScreenWidth  = GridWidth * DefaultCellSize
	ScreenHeight = GridHeight * DefaultCellSize
)

func CellW() float32 { return float32(ScreenWidth) / float32(GridWidth) }
func CellH() float32 { return float32(ScreenHeight) / float32(GridHeight) }
