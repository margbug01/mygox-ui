package motion

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// ============================================================================
// 29. AnimatedGroup
// ============================================================================

// AnimatedGroupOptions configures AnimatedGroup.
type AnimatedGroupOptions struct {
	Preset   Preset
	Stagger  time.Duration
	Duration time.Duration
	Delay    time.Duration
	Columns  int
	Gap      float32
	Trigger  int
}

// AnimatedGroup renders a staggered entrance group of items using fade, slide,
// scale, blur, or zoom presets. Incrementing Trigger replays the cascade.
func AnimatedGroup(c *ui.Context, items []func(), opts ...AnimatedGroupOptions) ui.Element {
	var o AnimatedGroupOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Preset == "" {
		o.Preset = PresetScale
	}
	if o.Stagger <= 0 {
		o.Stagger = 65 * time.Millisecond
	}
	if o.Duration <= 0 {
		o.Duration = 380 * time.Millisecond
	}
	if o.Gap <= 0 {
		o.Gap = 10
	}

	var host ui.Element
	if o.Columns > 1 {
		host = ui.Grid(c).Columns(o.Columns).Gap(o.Gap).WidthPercent(100)
	} else {
		host = ui.Row(c).Gap(o.Gap).Wrap().AlignItems(ui.Center)
	}

	st := ui.Local(host, "ag-state", func() triggerTimerState {
		return triggerTimerState{trigger: o.Trigger, start: c.Now()}
	})
	if st.trigger != o.Trigger {
		st.trigger = o.Trigger
		st.start = c.Now()
	}

	elapsed := c.Now().Sub(st.start)
	total := o.Delay + time.Duration(len(items))*o.Stagger + o.Duration
	if elapsed < total {
		c.AnimationFrame()
	}

	host.Children(func() {
		for i, fn := range items {
			start := o.Delay + time.Duration(i)*o.Stagger
			p := float32(easeOutCubic(float64(elapsed-start) / float64(o.Duration)))
			cell := ui.Box(c).Key(i)
			applyPreset(cell, o.Preset, p)
			cell.Children(func() {
				if fn != nil {
					fn()
				}
			})
		}
	})
	return host
}

// ============================================================================
// 30. Carousel
// ============================================================================

// CarouselOptions configures Carousel.
type CarouselOptions struct {
	SlideWidth     float32
	SlideHeight    float32
	Gap            float32
	ShowNavigation bool
	ShowIndicator  bool
	Spring         ui.SpringOptions
}

// Carousel renders a spring-sliding multi-card carousel with monochrome prev/next
// navigation buttons and clickable dot indicators.
func Carousel(c *ui.Context, items []func(), index *int, opts ...CarouselOptions) ui.Element {
	var o CarouselOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.SlideWidth <= 0 {
		o.SlideWidth = 240
	}
	if o.SlideHeight <= 0 {
		o.SlideHeight = 132
	}
	if o.Gap == 0 {
		o.Gap = 12
	}
	if len(opts) == 0 {
		o.ShowNavigation = true
		o.ShowIndicator = true
	}
	spring := resolveSpring(o.Spring, ui.SpringOptions{Stiffness: 260, Damping: 26, Mass: 0.8})

	cur := 0
	if index != nil {
		cur = *index
	}
	if len(items) > 0 {
		if cur < 0 {
			cur = 0
		}
		if cur >= len(items) {
			cur = len(items) - 1
		}
	}

	col := ui.Column(c).Width(o.SlideWidth).AlignItems(ui.Center).Gap(10)
	col.Children(func() {
		// Clipped viewport + prev/next controls overlay
		viewport := ui.Box(c).
			Width(o.SlideWidth).
			Height(o.SlideHeight).
			Radius(10).
			Clip()

		stride := o.SlideWidth + o.Gap
		targetX := -float32(cur) * stride
		trackX := viewport.Spring("car-x", targetX, spring)

		viewport.Children(func() {
			ui.Row(c).
				Height(o.SlideHeight).
				Gap(o.Gap).
				Translate(trackX, 0).
				Children(func() {
					for i, fn := range items {
						ui.Box(c).
							Key(i).
							Width(o.SlideWidth).
							Height(o.SlideHeight).
							Shrink(0).
							Children(func() {
								if fn != nil {
									fn()
								}
							})
					}
				})

			if o.ShowNavigation && len(items) > 1 {
				navBtn := func(key string, label string, enabled bool, isRight bool, onClick func()) {
					b := ui.Box(c).
						Key(key).
						Absolute().
						TopPercent(50).
						Translate(0, -13).
						Width(26).
						Height(26).
						Radius(999).
						Background(t.Surface).
						Border(1, t.Border).
						Center().
						ZIndex(2)
					if isRight {
						b.Right(8)
					} else {
						b.Left(8)
					}
					if enabled {
						b.Cursor(ui.CursorPointer)
						if b.Hovered() {
							b.Background(t.SurfaceHover).Border(1, t.BorderStrong)
						}
					} else {
						b.Opacity(0.35)
					}
					b.Children(func() {
						ui.Text(c, label).FontSize(13).FontWeight(600).TextColor(t.Text).SingleLine()
					})
					if enabled && b.Clicked() && onClick != nil {
						onClick()
					}
				}
				navBtn("prev", "‹", cur > 0, false, func() {
					if index != nil && *index > 0 {
						*index--
					}
				})
				navBtn("next", "›", cur < len(items)-1, true, func() {
					if index != nil && *index < len(items)-1 {
						*index++
					}
				})
			}
		})

		if o.ShowIndicator && len(items) > 1 {
			ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
				for i := range items {
					idx := i
					active := idx == cur
					dotCol := t.BorderStrong
					if active {
						dotCol = t.Text
					}
					dot := ui.Box(c).
						Key(idx).
						Width(6).
						Height(6).
						Radius(999).
						Background(dotCol).
						Cursor(ui.CursorPointer)
					if active {
						dot.Width(16)
					}
					if dot.Clicked() && index != nil {
						*index = idx
					}
				}
			})
		}
	})
	return col
}

// ============================================================================
// 31. InView
// ============================================================================

// InViewOptions configures InView.
type InViewOptions struct {
	Preset   Preset
	Once     bool
	Duration time.Duration
	Trigger  int
	Spring   ui.SpringOptions
}

type inViewState struct {
	seen    bool
	trigger int
	start   time.Time
}

// InView reveals its child with fade, slide, scale, or blur when entering the
// viewport (detected via Element.InView) or when Trigger is incremented to replay.
func InView(c *ui.Context, content func(), opts ...InViewOptions) ui.Element {
	var o InViewOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Preset == "" {
		o.Preset = PresetBlurSlide
	}
	if o.Duration <= 0 {
		o.Duration = 450 * time.Millisecond
	}

	box := ui.Box(c)
	now := c.Now()
	st := ui.Local(box, "iv-state", func() inViewState {
		return inViewState{trigger: o.Trigger, start: now}
	})
	if st.trigger != o.Trigger {
		st.trigger = o.Trigger
		st.seen = false
		st.start = now
	}

	visible := box.InView()
	if visible && !st.seen {
		st.seen = true
		st.start = now
	}
	active := visible || (o.Once && st.seen)

	p := float32(0)
	if active {
		elapsed := now.Sub(st.start)
		if elapsed < o.Duration {
			c.AnimationFrame()
		}
		p = float32(easeOutCubic(float64(elapsed) / float64(o.Duration)))
	}

	applyPreset(box, o.Preset, p)
	box.Children(func() {
		if content != nil {
			content()
		}
	})
	return box
}

// ============================================================================
// 32. InfiniteSlider
// ============================================================================

// InfiniteSliderOptions configures InfiniteSlider.
type InfiniteSliderOptions struct {
	Speed        float32 // pixels per second (default 48)
	SpeedOnHover float32 // pixels per second when hovered (default 14; negative disables slowdown)
	Gap          float32
	ItemWidth    float32
	Reverse      bool
	EdgeFade     bool
}

type sliderLocalState struct {
	offset   float64
	lastTime time.Time
	hasInit  bool
}

// InfiniteSlider renders a seamless infinite horizontal marquee strip that
// smoothly decelerates on hover (speedOnHover) without position jumps.
func InfiniteSlider(c *ui.Context, items []func(), opts ...InfiniteSliderOptions) ui.Element {
	var o InfiniteSliderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Speed <= 0 {
		o.Speed = 48
	}
	if o.SpeedOnHover == 0 {
		o.SpeedOnHover = 14
	}
	if o.Gap <= 0 {
		o.Gap = 14
	}
	if o.ItemWidth <= 0 {
		o.ItemWidth = 104
	}

	viewport := ui.Box(c).WidthPercent(100).Clip()
	if len(items) == 0 {
		return viewport
	}

	c.AnimationFrame()
	now := c.Now()
	st := ui.Local(viewport, "is-state", func() sliderLocalState {
		return sliderLocalState{lastTime: now, hasInit: true}
	})

	dt := now.Sub(st.lastTime).Seconds()
	if dt < 0 || dt > 0.25 {
		dt = 0.016
	}
	st.lastTime = now

	targetSpeed := o.Speed
	if viewport.Hovered() && o.SpeedOnHover >= 0 {
		targetSpeed = o.SpeedOnHover
	}
	curSpeed := viewport.Spring("is-spd", targetSpeed, ui.SpringOptions{Stiffness: 180, Damping: 26, Mass: 1})

	setWidth := float64(len(items)) * float64(o.ItemWidth+o.Gap)
	if setWidth < 1 {
		setWidth = 1
	}
	if o.Reverse {
		st.offset -= float64(curSpeed) * dt
	} else {
		st.offset += float64(curSpeed) * dt
	}
	for st.offset >= setWidth {
		st.offset -= setWidth
	}
	for st.offset < 0 {
		st.offset += setWidth
	}

	if o.EdgeFade {
		viewport.MaskLinearGradient(ui.LinearGradient{
			Angle: 90,
			Stops: []ui.GradientStop{
				{Color: ui.Transparent, At: 0},
				{Color: ui.Hex("#000000"), At: 0.1},
				{Color: ui.Hex("#000000"), At: 0.9},
				{Color: ui.Transparent, At: 1},
			},
		})
	}

	viewport.Children(func() {
		ui.Row(c).
			AlignItems(ui.Center).
			Gap(o.Gap).
			Translate(-float32(st.offset), 0).
			Children(func() {
				// Render 3 contiguous sets so wide viewports never show empty seams.
				for rep := 0; rep < 3; rep++ {
					for i, fn := range items {
						ui.Box(c).
							Key(rep*1000 + i).
							Width(o.ItemWidth).
							Shrink(0).
							Children(func() {
								if fn != nil {
									fn()
								}
							})
					}
				}
			})
	})
	return viewport
}

// ============================================================================
// 33. ScrollProgress
// ============================================================================

// ScrollProgressOptions configures ScrollProgress.
type ScrollProgressOptions struct {
	Height     float32
	Radius     float32
	Color      ui.Color
	TrackColor ui.Color
	Spring     ui.SpringOptions
}

// ScrollProgress renders a live spring-smoothed progress bar bound to a
// scrollable container's ui.ScrollState (scroll.Y / scroll.MaxY).
func ScrollProgress(c *ui.Context, scroll *ui.ScrollState, opts ...ScrollProgressOptions) ui.Element {
	var o ScrollProgressOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Height <= 0 {
		o.Height = 3
	}
	if o.Radius == 0 {
		o.Radius = 999
	}
	barCol := colorOr(o.Color, t.Text)
	trackCol := colorOr(o.TrackColor, t.SurfaceHover)
	spring := resolveSpring(o.Spring, ui.SpringOptions{Stiffness: 280, Damping: 24, Mass: 0.5})

	var target float32
	if scroll != nil && scroll.MaxY > 0 {
		target = clamp32(scroll.Y/scroll.MaxY, 0, 1)
	}

	track := ui.Box(c).
		WidthPercent(100).
		Height(o.Height).
		Radius(o.Radius).
		Background(trackCol).
		Clip()

	p := clamp32(track.Spring("sp-val", target, spring), 0, 1)

	track.Children(func() {
		ui.Box(c).
			WidthPercent(p * 100).
			HeightPercent(100).
			Radius(o.Radius).
			Background(barCol)
	})
	return track
}
