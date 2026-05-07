package game

import (
	"bytes"
	"io"
	"snake_golang/assets"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

type Audio struct {
	ctx  *audio.Context
	Game *audio.Player
	Menu *audio.Player
	Eat  *audio.Player
}

func NewAudio(ctx *audio.Context) *Audio {
	return &Audio{
		ctx: ctx,
	}
}

func (a *Audio) Reload(t *assets.Theme) error {
	for _, p := range []**audio.Player{&a.Game, &a.Menu, &a.Eat} {
		if *p != nil {
			(*p).Pause()
			_ = (*p).Close()
			*p = nil
		}
	}
	mk := func(raw []byte, loop bool, vol float64) (*audio.Player, error) {
		s, err := vorbis.DecodeWithSampleRate(assets.SampleRate, bytes.NewReader(raw))
		if err != nil {
			return nil, err
		}
		var src io.Reader = s
		if loop {
			src = audio.NewInfiniteLoop(s, s.Length())
		}
		p, err := a.ctx.NewPlayer(src)
		if err != nil {
			return nil, err
		}
		p.SetVolume(vol)
		return p, nil
	}
	var err error
	if a.Game, err = mk(t.GameMusicOGG, true, 0.4); err != nil {
		return err
	}
	if a.Menu, err = mk(t.MenuMusicOGG, true, 0.4); err != nil {
		return err
	}
	if a.Eat, err = mk(t.EatSoundOGG, false, 0.8); err != nil {
		return err
	}
	return nil
}
