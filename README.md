# mygox-ui

A unified, coss-styled pure-Go UI suite for [MyGo](https://github.com/egoist/mygo) (`v0.10.0`), porting three best-in-class open-source ecosystems into one zero-WebView repository:

- **`mygox-ui/coss`** — All **56 primitives** and **530 particles** ported 1:1 from [coss.com/ui](https://coss.com/ui/particles), 102 lazy-parsed Lucide SVG icons, Light & Dark design tokens, and a shadcn-style source copier CLI (`cmd/mygox-ui`).
- **`mygox-ui/flow`** — Interactive GPU-accelerated node graph canvas modeled after [React Flow UI](https://reactflow.dev/ui): `BaseNode`, `DatabaseSchemaNode` (row-level foreign-key handles), `LabeledGroupNode` (drag container moves children), `PlaceholderNode` (click to promote), `AnnotationNode`, spinning conic `NodeStatusIndicator`, `ButtonEdge` midpoint delete, floating `NodeToolbar`, `ZoomSlider`, interactive `MiniMap`, Kahn auto-layout, Undo/Redo, and JSON persistence.
- **`mygox-ui/motion`** — All **33 / 33 components** ported 1:1 from [ibelick/motion-primitives](https://github.com/ibelick/motion-primitives) in restrained monochrome `coss` style (`BorderTrail`, `Spotlight`, `Tilt`, `GlowEffect`, `ProgressiveBlur`, `ImageComparison`, `AnimatedNumber`, `SlidingNumber`, `SpinningText`, `TextEffect`, `TextLoop`, `TextMorph`, `TextRoll`, `TextScramble`, `TextShimmer`, `TextShimmerWave`, `AnimatedBackground`, `Dock`, `Magnetic`, `Cursor`, `ToolbarDynamic`, `ToolbarExpandable`, `Accordion`, `Disclosure`, `Dialog`, `MorphingDialog`, `MorphingPopover`, `TransitionPanel`, `AnimatedGroup`, `Carousel`, `InView`, `InfiniteSlider`, `ScrollProgress`).

**Live Showcase & WebGL 2 WASM Playground**: https://margbug01.github.io/mygox-ui/

## Install

```bash
go get github.com/margbug01/mygox-ui
```

```go
import (
    "github.com/egoist/mygo/ui"
    "github.com/margbug01/mygox-ui/coss"
    "github.com/margbug01/mygox-ui/flow"
    "github.com/margbug01/mygox-ui/motion"
)
```

## Run the Unified Showcase

```bash
go run ./cmd/showcase
```

## License

MIT © 2026 margbug01
