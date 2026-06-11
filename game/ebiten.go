package game

import (
	"fmt"
	"image/color"
	"log"
	"snake_golang/assets"
	"snake_golang/assets/mods"
	"snake_golang/assets/skins"
	"snake_golang/game/i18n"
	gamelb "snake_golang/game/leaderboard"
	"snake_golang/game/menu"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	etxt "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	backgroundColorScale = float32(0.55)
	overlayGray          = uint8(28)
	overlayAlpha         = uint8(130)

	gameMusicVolume       = 0.4
	gameMusicPausedVolume = 0.2
)

type Screen struct {
	World      *World
	Menu       *menu.Menu
	FaceSource *etxt.GoTextFaceSource
	Audio      *Audio

	PlayerName string
	PlayerID   string
	nameDraft  []rune
	nameError  string

	codeDraft []rune
	codeError string

	LeaderboardClient *gamelb.Client

	scoreSubmitCh       chan scoreSubmitResult
	scoreSubmitState    scoreSubmitState
	scoreSubmitRank     int
	scoreSubmitImproved bool
	scoreSubmitError    string

	leaderboardCh      chan leaderboardFetchResult
	leaderboardLoading bool
	leaderboardError   string
	leaderboardEntries []gamelb.Entry

	versionInfoCh  chan versionInfoResult
	updateRequired bool
	latestVersion  string
}

func (s *Screen) StartGame() {
	if s.World == nil {
		return
	}
	s.resetScoreSubmitState()
	s.World.WrapEdges = mods.Current() == mods.Defaltyk
	s.World.State = StateWaiting
	s.World.LastMove = time.Now()
}

func (s *Screen) ScreenSize() (int, int) {
	return ScreenWidth, ScreenHeight
}

func (s *Screen) SwitchMode() {
	mods.Set(mods.Next())
}

func (s *Screen) SwitchSkin() {
	skins.Set(skins.Next())
	assets.SetActive(skins.Current().ID())
	if s.Audio != nil {
		if err := s.Audio.Reload(assets.Current()); err != nil {
			log.Printf("audio reload: %v", err)
			return
		}
		s.syncMusic()
	}
}

func (s *Screen) syncMusic() {
	if s.World == nil || s.Audio == nil {
		return
	}
	state := s.World.State

	if menuP := s.Audio.Menu; menuP != nil {
		if state == StateMenu {
			if !menuP.IsPlaying() {
				menuP.Play()
			}
		} else if menuP.IsPlaying() {
			menuP.Pause()
		}
	}

	if gameP := s.Audio.Game; gameP != nil {
		switch state {
		case StatePlaying:
			gameP.SetVolume(gameMusicVolume)
			if !gameP.IsPlaying() {
				gameP.Play()
			}
		case StatePaused:
			gameP.SetVolume(gameMusicPausedVolume)
			if !gameP.IsPlaying() {
				gameP.Play()
			}
		default:
			if gameP.IsPlaying() {
				gameP.Pause()
			}
			gameP.SetVolume(gameMusicVolume)
		}
	}
}

func (s *Screen) playEatSound() {
	if s.Audio == nil || s.Audio.Eat == nil {
		return
	}
	if err := s.Audio.Eat.Rewind(); err != nil {
		return
	}
	s.Audio.Eat.Play()
}

func (s *Screen) Update() error {
	w := s.World
	if w == nil {
		return nil
	}

	if s.Menu == nil {
		s.Menu = menu.NewMenu(s)
	}

	s.pollLeaderboardAsync()

	if w.State == StateNameInput {
		s.syncMusic()
		return s.UpdateNameInput()
	}
	if w.State == StateCodeInput {
		s.syncMusic()
		return s.UpdateCodeInput()
	}

	if w.State == StateMenu {
		s.syncMusic()
		if inpututil.IsKeyJustPressed(ebiten.KeyN) {
			s.BeginNameInput()
			return nil
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			s.BeginCodeInput()
			return nil
		}
		return s.Menu.Update()
	}

	if w.State == StateLeaderboard {
		s.syncMusic()
		return s.UpdateLeaderboard()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) &&
		(w.State == StatePaused || w.State == StateGameOver) {
		s.resetScoreSubmitState()
		s.World = NewWorld()
		w = s.World
		s.syncMusic()
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		switch w.State {
		case StateWaiting:
			w.State = StatePlaying
			w.LastMove = time.Now()
			w.StartedAt = time.Now()
		case StatePlaying:
			w.State = StatePaused
		case StatePaused:
			w.State = StatePlaying
			w.LastMove = time.Now()
		case StateGameOver:
			s.resetScoreSubmitState()
			s.World = newWorldPlaying()
			w = s.World
		}
	}

	s.syncMusic()

	if w.State == StateGameOver {
		return nil
	}
	if w.State != StatePlaying {
		return nil
	}

	nextDirection := w.NextDirection

	if inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		nextDirection = Up
	} else if inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		nextDirection = Left
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) || inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		nextDirection = Down
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		nextDirection = Right
	}

	if !IsOppositeDirection(w.Direction, nextDirection) {
		w.NextDirection = nextDirection
	}

	if time.Since(w.LastMove) < w.moveInterval() {
		return nil
	}
	w.LastMove = time.Now()
	previousState := w.State
	prevScore := w.Score
	Step(w)
	if w.Score > prevScore {
		s.playEatSound()
	}
	if previousState != StateGameOver && w.State == StateGameOver {
		s.onGameOver()
	}
	return nil
}

func (s *Screen) Draw(screen *ebiten.Image) {
	if s.World != nil && s.World.State == StateNameInput {
		DrawNameInput(screen, s.FaceSource, string(s.nameDraft), s.nameError, s.PlayerName != "")
		return
	}
	if s.World != nil && s.World.State == StateCodeInput {
		DrawCodeInput(screen, s.FaceSource, string(s.codeDraft), s.codeError)
		return
	}
	if s.World != nil && s.World.State == StateMenu {
		menu.Draw(screen, s.FaceSource, s.Menu)
		DrawMenuPlayerName(screen, s.FaceSource, s.PlayerName)
		DrawMenuCodeHint(screen, s.FaceSource)
		if s.updateRequired {
			DrawUpdateBanner(screen, s.FaceSource, s.latestVersion)
		}
		return
	}
	if s.World != nil && s.World.State == StateLeaderboard {
		DrawLeaderboardScreen(screen, s.FaceSource, s)
		return
	}
	DrawWorld(s.World, s.FaceSource, screen)
	if s.World.State == StatePlaying || s.World.State == StatePaused {
		DrawScoreHud(screen, s.FaceSource, s.World)
	}
	if s.World != nil && s.World.State == StateGameOver {
		DrawGameOverSubmitStatus(screen, s.FaceSource, s)
		if s.updateRequired {
			DrawUpdateBanner(screen, s.FaceSource, s.latestVersion)
		}
	}
}

func DrawUpdateBanner(dst *ebiten.Image, face *etxt.GoTextFaceSource, latest string) {
	if dst == nil || face == nil {
		return
	}
	msg := i18n.T("update.required")
	if strings.TrimSpace(latest) != "" {
		msg = fmt.Sprintf(i18n.T("update.required_version"), latest)
	}

	bannerH := float32(ScreenHeight) * 0.09
	if bannerH < 28 {
		bannerH = 28
	}
	vector.FillRect(dst, 0, 0, float32(ScreenWidth), bannerH,
		color.RGBA{R: 178, G: 34, B: 34, A: 235}, true)

	size := FitFontSize(face, msg, UiLabelFontSize()*0.42, float64(ScreenWidth)*0.94, 8)
	fo := &etxt.GoTextFace{Source: face, Size: size}
	op := &etxt.DrawOptions{}
	op.GeoM.Translate(float64(ScreenWidth)/2, float64(bannerH)/2)
	op.PrimaryAlign = etxt.AlignCenter
	op.SecondaryAlign = etxt.AlignCenter
	op.ColorScale.ScaleWithColor(color.White)
	etxt.Draw(dst, msg, fo, op)
}

func (s *Screen) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if outsideWidth <= 0 || outsideHeight <= 0 {
		return ScreenWidth, ScreenHeight
	}
	ScreenWidth, ScreenHeight = outsideWidth, outsideHeight
	return ScreenWidth, ScreenHeight
}

func DrawWorld(w *World, face *etxt.GoTextFaceSource, dst *ebiten.Image) {
	if w == nil {
		return
	}
	active := assets.Current()
	if active == nil {
		return
	}
	if bg := active.Background; bg != nil {
		op := &ebiten.DrawImageOptions{}
		bw := float64(bg.Bounds().Dx())
		bh := float64(bg.Bounds().Dy())
		if bw > 0 && bh > 0 {
			op.GeoM.Scale(float64(ScreenWidth)/bw, float64(ScreenHeight)/bh)
		}
		op.ColorScale.Scale(backgroundColorScale, backgroundColorScale, backgroundColorScale, 1)
		dst.DrawImage(bg, op)
	}

	vector.FillRect(dst, 0, 0, float32(ScreenWidth), float32(ScreenHeight),
		color.RGBA{R: overlayGray, G: overlayGray, B: overlayGray + 4, A: overlayAlpha}, true)

	if w.State != StateWaiting {
		cw, ch := CellW(), CellH()
		drawSnake(dst, w, cw, ch)
		ax, ay := float32(w.Apple.X)*cw, float32(w.Apple.Y)*ch
		if apple := active.Apple; apple != nil {
			iw := float64(apple.Bounds().Dx())
			ih := float64(apple.Bounds().Dy())
			op := &ebiten.DrawImageOptions{}
			if iw > 0 && ih > 0 {
				op.GeoM.Scale(float64(cw)/iw, float64(ch)/ih)
			}
			op.GeoM.Translate(float64(ax), float64(ay))
			dst.DrawImage(apple, op)
		} else {
			vector.FillRect(dst, ax, ay, cw, ch, color.RGBA{R: 255, G: 0, B: 0, A: 255}, true)
		}
	}

	if face == nil {
		return
	}
	switch w.State {
	case StateWaiting:
		DrawCenteredLabel(dst, face, UiLabelFontSize(), i18n.T("game.press_space"))
	case StatePaused:
		DrawCenteredLabel(dst, face, UiLabelFontSize(), i18n.T("game.pause"))
	case StateGameOver:
		DrawCenteredLabel(dst, face, GameOverLabelFontSize(),
			fmt.Sprintf(i18n.T("game.game_over"), w.Score))
	}
}

func UiLabelFontSize() float64 {
	m := min(ScreenWidth, ScreenHeight)
	s := float64(m) * 0.075
	if s < 12 {
		return 12
	}
	if s > 72 {
		return 72
	}
	return s
}

func GameOverLabelFontSize() float64 {
	s := UiLabelFontSize() * 0.48
	if s < 11 {
		return 11
	}
	if s > 24 {
		return 24
	}
	return s
}

func DrawCenteredLabel(dst *ebiten.Image, face *etxt.GoTextFaceSource, size float64, text string) {
	maxW := float64(ScreenWidth) * 0.92
	size = FitFontSize(face, text, size, maxW, 8)
	fo := &etxt.GoTextFace{Source: face, Size: size}
	op := &etxt.DrawOptions{}
	op.GeoM.Translate(float64(ScreenWidth)/2, float64(ScreenHeight)/2)
	op.PrimaryAlign = etxt.AlignCenter
	op.SecondaryAlign = etxt.AlignCenter
	op.ColorScale.ScaleWithColor(color.White)
	etxt.Draw(dst, text, fo, op)
}

func FitFontSize(face *etxt.GoTextFaceSource, s string, size, maxW, minSize float64) float64 {
	if face == nil || s == "" || maxW <= 0 {
		return size
	}
	fo := &etxt.GoTextFace{Source: face, Size: size}
	w, _ := etxt.Measure(s, fo, 0)
	if w <= maxW {
		return size
	}
	scaled := size * maxW / w
	if scaled < minSize {
		scaled = minSize
	}
	return scaled
}
