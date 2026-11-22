# mygox-ui

A unified, coss-styled pure-Go UI suite for [MyGo](https://github.com/egoist/mygo) (`v0.10.0`), combining three zero-WebView packages in one repository:

- **`mygox-ui/coss`** — All **56 primitives** and **530 particles** ported from [coss.com/ui](https://coss.com/ui/particles), 102 lazy-parsed Lucide SVG icons, Light & Dark design tokens, and a shadcn-style source-copying CLI (`cmd/mygox-ui`).
- **`mygox-ui/flow`** — Interactive GPU-accelerated node graph & workflow canvas (infinite pan/zoom, magnetic port wiring, Bezier/SmoothStep/Straight animated edges, alignment guides, Undo/Redo, Kahn auto-layout, interactive MiniMap, JSON persistence).
- **`mygox-ui/motion`** — **28 GPU-rendered motion & visual effects** (`BorderBeam`, `SpotlightCard`, `TiltCard`, `GlowCard`, `ShineBorder`, `NumberTicker`, `TextShimmer`, `TextScramble`, `Typewriter`, `BlurReveal`, `FlipWords`, `Dock`, `Magnetic`, `ShimmerButton`, `RippleButton`, `PulsingBadge`, `MorphTabs`, `SpringCounter`, `AuroraBackground`, `DotPattern`, `GridPattern`, `Meteors`, `Sparkles`, `RippleRings`, `Marquee`, `AnimatedBeam`, `OrbitingCircles`, `ConfettiBurst`).

**Live Showcase & WebGL 2 WASM Playground**: https://margbug01.github.io/mygox-ui/

## Install

```bash
go get github.com/margbug01/mygox-ui
```

Zero package-name collisions with `github.com/egoist/mygo/ui`:

```go
import (
    "github.com/egoist/mygo/ui"
    "github.com/margbug01/mygox-ui/coss"
    "github.com/margbug01/mygox-ui/coss/theme"
    "github.com/margbug01/mygox-ui/flow"
    "github.com/margbug01/mygox-ui/motion"
)
```

## Run the Unified Showcase

```bash
go run ./cmd/showcase
```

Switch seamlessly between **UI (530)**, **Flow (Canvas)**, and **Motion (28)** — all styled in the same Light/Dark `coss` design language, with mobile responsive drawers and WebAssembly support (`./scripts/build-site.sh dist`).

## Copy-Paste CLI (shadcn-style for Go)

```bash
go run ./cmd/mygox-ui list
go run ./cmd/mygox-ui init
go run ./cmd/mygox-ui add button dialog command
```

## License

MIT © 2026 margbug01
