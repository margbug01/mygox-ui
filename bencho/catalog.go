package bencho

import "github.com/egoist/mygo/ui"

// Interaction identifies the primary input gesture of a block.
type Interaction string

const (
	InteractionPress   Interaction = "Press"
	InteractionHover   Interaction = "Hover"
	InteractionDrag    Interaction = "Drag"
	InteractionLoading Interaction = "Loading"
	InteractionType    Interaction = "Type"
	InteractionSelect  Interaction = "Select"
	InteractionSwipe   Interaction = "Swipe"
	InteractionSlide   Interaction = "Slide"
)

// AllInteractions lists all interaction categories in display order.
var AllInteractions = []Interaction{
	InteractionPress,
	InteractionHover,
	InteractionDrag,
	InteractionLoading,
	InteractionType,
	InteractionSelect,
	InteractionSwipe,
	InteractionSlide,
}

// Block describes one entry in the bencho.dev catalog.
// A non-nil Demo marks the block as implemented.
type Block struct {
	Slug        string
	Title       string
	Interaction Interaction
	Summary     string
	Demo        func(c *ui.Context)
}

// Implemented reports whether the block has a live MyGo implementation.
func (b Block) Implemented() bool {
	return b.Demo != nil
}

var blockSpecs = []Block{
	{Slug: "thinking-orb", Title: "Thinking orb", Interaction: InteractionHover},
	{Slug: "elastic-slider", Title: "Elastic slider", Interaction: InteractionDrag},
	{Slug: "live-vote", Title: "Live vote", Interaction: InteractionPress},
	{Slug: "confirm-send", Title: "Confirm send", Interaction: InteractionPress},
	{Slug: "huddle", Title: "Huddle", Interaction: InteractionPress},
	{Slug: "stories", Title: "Stories", Interaction: InteractionHover},
	{Slug: "glide-tip", Title: "Liquid tooltip", Interaction: InteractionHover},
	{Slug: "pass-reveal", Title: "Password reveal", Interaction: InteractionPress},
	{Slug: "race-bars", Title: "Race bars", Interaction: InteractionLoading},
	{Slug: "path-panel", Title: "Path panel", Interaction: InteractionDrag},
	{Slug: "income-flow", Title: "Income flow", Interaction: InteractionHover},
	{Slug: "wave-trim", Title: "Audio trim", Interaction: InteractionDrag},
	{Slug: "dice", Title: "Dice", Interaction: InteractionDrag},
	{Slug: "pillars", Title: "3D bar chart", Interaction: InteractionHover},
	{Slug: "market-chart", Title: "Market chart", Interaction: InteractionHover},
	{Slug: "seat-map", Title: "Seat map", Interaction: InteractionPress},
	{Slug: "gauge", Title: "Score gauge", Interaction: InteractionDrag},
	{Slug: "listen", Title: "AI Assistant", Interaction: InteractionPress},
	{Slug: "hold-delete", Title: "Hold to delete", Interaction: InteractionPress},
	{Slug: "faq-frames", Title: "FAQ frames", Interaction: InteractionPress},
	{Slug: "status-picker", Title: "Status picker", Interaction: InteractionPress},
	{Slug: "card-stack", Title: "Card stack", Interaction: InteractionHover},
	{Slug: "expenses", Title: "Arc ring", Interaction: InteractionHover},
	{Slug: "format-bar", Title: "Format bar", Interaction: InteractionPress},
	{Slug: "jelly-tabs", Title: "Jelly tabs", Interaction: InteractionPress},
	{Slug: "odds", Title: "Odds", Interaction: InteractionPress},
	{Slug: "progress-tracker", Title: "Progress tracker", Interaction: InteractionPress},
	{Slug: "thinking-loader", Title: "Thinking loader", Interaction: InteractionLoading},
	{Slug: "genie", Title: "Genie effect", Interaction: InteractionPress},
	{Slug: "flap-board", Title: "Flap board", Interaction: InteractionPress},
	{Slug: "image-trail", Title: "Image trail", Interaction: InteractionHover},
	{Slug: "dust-delete", Title: "Dust delete", Interaction: InteractionPress},
	{Slug: "slingshot", Title: "Slingshot", Interaction: InteractionDrag},
	{Slug: "crumple", Title: "Crumpled receipt", Interaction: InteractionDrag},
	{Slug: "price-chart", Title: "Price chart", Interaction: InteractionDrag},
	{Slug: "standing-calendar", Title: "Standing calendar", Interaction: InteractionPress},
	{Slug: "orders-globe", Title: "Orders globe", Interaction: InteractionDrag},
	{Slug: "ai-autocomplete", Title: "AI autocomplete", Interaction: InteractionType},
	{Slug: "magnetic-letters", Title: "Magnetic letters", Interaction: InteractionHover},
	{Slug: "selection-actions", Title: "Selection actions", Interaction: InteractionPress},
	{Slug: "loading-state", Title: "Loading state", Interaction: InteractionLoading},
	{Slug: "private-balance", Title: "Private balance", Interaction: InteractionPress},
	{Slug: "spec-inspector", Title: "Spec inspector", Interaction: InteractionDrag},
	{Slug: "context-meter", Title: "Context meter", Interaction: InteractionHover},
	{Slug: "load-cards", Title: "Load cards", Interaction: InteractionPress},
	{Slug: "ink-bloom", Title: "Ink bloom", Interaction: InteractionPress},
	{Slug: "peek-rating", Title: "Peek rating", Interaction: InteractionHover},
	{Slug: "edge-light", Title: "Edge light", Interaction: InteractionHover},
	{Slug: "line-sidebar", Title: "Line sidebar", Interaction: InteractionHover},
	{Slug: "gravity-tags", Title: "Gravity tags", Interaction: InteractionHover},
	{Slug: "falling-text", Title: "Falling text", Interaction: InteractionHover},
	{Slug: "card-swap", Title: "Card swap", Interaction: InteractionDrag},
	{Slug: "dot-swarm", Title: "Dot swarm", Interaction: InteractionHover},
	{Slug: "option-wheel", Title: "Option wheel", Interaction: InteractionDrag},
	{Slug: "wake-slider", Title: "Wake slider", Interaction: InteractionDrag},
	{Slug: "line-sphere", Title: "Line sphere", Interaction: InteractionHover},
	{Slug: "page-peel", Title: "Page peel", Interaction: InteractionDrag},
	{Slug: "stacked-bars", Title: "Stacked bars", Interaction: InteractionHover},
	{Slug: "blast-radius", Title: "Network graph", Interaction: InteractionHover},
	{Slug: "connect-sources", Title: "Integrations", Interaction: InteractionPress},
	{Slug: "blueprint-loader", Title: "Blueprint loader", Interaction: InteractionLoading},
	{Slug: "ai-chat-input", Title: "AI chat input", Interaction: InteractionType},
	{Slug: "faq-accordion", Title: "FAQ accordion", Interaction: InteractionPress},
	{Slug: "avatar-group", Title: "Avatar group", Interaction: InteractionHover},
	{Slug: "profile-card", Title: "Profile hover card", Interaction: InteractionHover},
	{Slug: "resizable-panels", Title: "Resizable panels", Interaction: InteractionDrag},
	{Slug: "toast-stack", Title: "Toast stack", Interaction: InteractionHover},
	{Slug: "dot-loader", Title: "Dot loader", Interaction: InteractionLoading},
	{Slug: "image-shimmer", Title: "Image shimmer", Interaction: InteractionLoading},
	{Slug: "shimmer-text", Title: "Shimmer text", Interaction: InteractionLoading},
	{Slug: "orbit-cards", Title: "Orbit cards", Interaction: InteractionDrag},
	{Slug: "metaballs", Title: "Metaballs", Interaction: InteractionHover},
	{Slug: "tree-view", Title: "Tree view", Interaction: InteractionPress},
	{Slug: "access-request", Title: "Access request", Interaction: InteractionPress},
	{Slug: "code-reveal", Title: "Code reveal", Interaction: InteractionHover},
	{Slug: "budget-slider", Title: "Budget slider", Interaction: InteractionDrag},
	{Slug: "emoji-reactions", Title: "Emoji reactions", Interaction: InteractionHover},
	{Slug: "tag-input", Title: "Tag input", Interaction: InteractionType},
	{Slug: "country-wheel", Title: "Country wheel", Interaction: InteractionDrag},
	{Slug: "folder", Title: "Folder", Interaction: InteractionHover},
	{Slug: "foggy-word", Title: "Foggy word", Interaction: InteractionHover},
	{Slug: "coin", Title: "Coin", Interaction: InteractionDrag},
	{Slug: "text-fade", Title: "Text fade", Interaction: InteractionPress},
	{Slug: "text-rise", Title: "Text rise", Interaction: InteractionPress},
	{Slug: "text-blur", Title: "Text blur", Interaction: InteractionPress},
	{Slug: "tab-bar", Title: "Tab bar", Interaction: InteractionPress},
	{Slug: "range", Title: "Range", Interaction: InteractionDrag},
	{Slug: "tape", Title: "Tape", Interaction: InteractionDrag},
	{Slug: "pay-button", Title: "Pay button", Interaction: InteractionPress},
	{Slug: "collapsing-sidebar", Title: "Collapsing sidebar", Interaction: InteractionPress},
	{Slug: "copy-field", Title: "Copy field", Interaction: InteractionPress},
	{Slug: "fold-dropdown", Title: "Dropdown", Interaction: InteractionPress},
	{Slug: "confetti", Title: "Confetti", Interaction: InteractionPress},
	{Slug: "picker-wheel", Title: "Picker wheel", Interaction: InteractionDrag},
	{Slug: "split-view", Title: "Split view", Interaction: InteractionDrag},
	{Slug: "album-shelf", Title: "Album shelf", Interaction: InteractionHover},
	{Slug: "piano", Title: "Piano", Interaction: InteractionPress},
	{Slug: "hub-diagram", Title: "Hub diagram", Interaction: InteractionHover},
	{Slug: "billing-card", Title: "Billing card", Interaction: InteractionPress},
	{Slug: "wave-loader", Title: "Wave loader", Interaction: InteractionLoading},
	{Slug: "bounce-loader", Title: "Bounce loader", Interaction: InteractionLoading},
	{Slug: "trace-loader", Title: "Trace loader", Interaction: InteractionLoading},
	{Slug: "morph-loader", Title: "Morph loader", Interaction: InteractionLoading},
	{Slug: "pixel-loader", Title: "Pixel loader", Interaction: InteractionLoading},
	{Slug: "blocks-loader", Title: "Blocks loader", Interaction: InteractionLoading},
	{Slug: "flight-island", Title: "Flight island", Interaction: InteractionPress},
	{Slug: "music-island", Title: "Music island", Interaction: InteractionPress},
	{Slug: "call-island", Title: "Call island", Interaction: InteractionPress},
	{Slug: "pill-accordion", Title: "Pill accordion", Interaction: InteractionPress},
	{Slug: "resize-frame", Title: "Resize frame", Interaction: InteractionDrag},
	{Slug: "progressive-blur", Title: "Progressive blur", Interaction: InteractionHover},
	{Slug: "feedback-bar", Title: "Feedback bar", Interaction: InteractionPress},
	{Slug: "price-flow", Title: "Price flow", Interaction: InteractionPress},
	{Slug: "card-fan", Title: "Card fan", Interaction: InteractionPress},
	{Slug: "smudge", Title: "Smudge", Interaction: InteractionHover},
	{Slug: "swipe-deck", Title: "Swipe deck", Interaction: InteractionSwipe},
	{Slug: "radial-pop-out", Title: "Radial pop-out", Interaction: InteractionHover},
	{Slug: "form-error", Title: "Form error", Interaction: InteractionType},
	{Slug: "heart-rate", Title: "Heart rate", Interaction: InteractionHover},
	{Slug: "light", Title: "Light", Interaction: InteractionHover},
	{Slug: "segment-timer", Title: "Segment timer", Interaction: InteractionDrag},
	{Slug: "password-strength", Title: "Password strength", Interaction: InteractionType},
	{Slug: "thinking", Title: "Thinking", Interaction: InteractionPress},
	{Slug: "activity-stack", Title: "Activity stack", Interaction: InteractionSwipe},
	{Slug: "ticket", Title: "Ticket", Interaction: InteractionDrag},
	{Slug: "share-pill", Title: "Share pill", Interaction: InteractionHover},
	{Slug: "voice-chat", Title: "Voice chat", Interaction: InteractionPress},
	{Slug: "world-clock", Title: "World clock", Interaction: InteractionDrag},
	{Slug: "skeleton-morph", Title: "Skeleton morph", Interaction: InteractionPress},
	{Slug: "date-reel", Title: "Date reel", Interaction: InteractionSwipe},
	{Slug: "folding-bill", Title: "Folding bill", Interaction: InteractionDrag},
	{Slug: "magnetic-dots", Title: "Magnetic dots", Interaction: InteractionHover},
	{Slug: "dot-globe", Title: "Dot globe", Interaction: InteractionDrag},
	{Slug: "sweep-button", Title: "Sweep button", Interaction: InteractionHover},
	{Slug: "asset-swap", Title: "Asset swap", Interaction: InteractionPress},
	{Slug: "ascii-wake", Title: "ASCII wake", Interaction: InteractionHover},
	{Slug: "heat-word", Title: "Heat map", Interaction: InteractionHover},
	{Slug: "image-compare", Title: "Image compare", Interaction: InteractionDrag},
	{Slug: "like", Title: "Like", Interaction: InteractionPress},
	{Slug: "signature", Title: "Signature pad", Interaction: InteractionDrag},
	{Slug: "voice-note", Title: "Voice note", Interaction: InteractionPress},
	{Slug: "eye-tracker", Title: "Eye tracker", Interaction: InteractionHover},
	{Slug: "dynamic-island", Title: "Dynamic island", Interaction: InteractionPress},
	{Slug: "time-scrubber", Title: "Time scrubber", Interaction: InteractionSlide},
	{Slug: "upload-dropzone", Title: "Upload dropzone", Interaction: InteractionDrag},
	{Slug: "rolling-counter", Title: "Rolling counter", Interaction: InteractionDrag},
	{Slug: "particles", Title: "Particles", Interaction: InteractionHover},
	{Slug: "label-input", Title: "Label input", Interaction: InteractionType},
	{Slug: "one-time-code", Title: "One-time code", Interaction: InteractionType},
	{Slug: "generate", Title: "Generate", Interaction: InteractionPress},
	{Slug: "step-player", Title: "Step player", Interaction: InteractionPress},
	{Slug: "todo-tower", Title: "Todo tower", Interaction: InteractionPress},
	{Slug: "image-accordion", Title: "Image accordion", Interaction: InteractionHover},
	{Slug: "glass-bubble", Title: "Glass bubble", Interaction: InteractionDrag},
	{Slug: "action-node", Title: "Action node", Interaction: InteractionHover},
	{Slug: "magnet-select", Title: "Magnetic select", Interaction: InteractionSelect},
	{Slug: "slide-confirm", Title: "Slide to confirm", Interaction: InteractionDrag},
	{Slug: "picker", Title: "Assignees", Interaction: InteractionSelect},
	{Slug: "checklist", Title: "Checklist", Interaction: InteractionPress},
	{Slug: "carousel", Title: "Carousel", Interaction: InteractionSwipe},
	{Slug: "palette", Title: "Palette", Interaction: InteractionPress},
	{Slug: "aspect", Title: "Aspect ratio", Interaction: InteractionSelect},
	{Slug: "tilt", Title: "Tilt card", Interaction: InteractionHover},
	{Slug: "sound", Title: "Now playing", Interaction: InteractionPress},
	{Slug: "drag-ball", Title: "Dragging ball", Interaction: InteractionDrag},
	{Slug: "seek", Title: "Search", Interaction: InteractionPress},
	{Slug: "pull", Title: "Pull to refresh", Interaction: InteractionDrag},
	{Slug: "escape", Title: "Escape button", Interaction: InteractionHover},
	{Slug: "slosh", Title: "Slosh slider", Interaction: InteractionSlide},
	{Slug: "liq-create", Title: "Create menu", Interaction: InteractionPress},
	{Slug: "liq-arrange", Title: "Reorder list", Interaction: InteractionDrag},
	{Slug: "toolbar", Title: "Canvas toolbar", Interaction: InteractionSelect},
	{Slug: "radial", Title: "Radial menu", Interaction: InteractionPress},
	{Slug: "stepper", Title: "Drag stepper", Interaction: InteractionPress},
	{Slug: "confirm", Title: "Inline confirm", Interaction: InteractionPress},
	{Slug: "toasts", Title: "Notify", Interaction: InteractionPress},
	{Slug: "icon-bar", Title: "Icon bar", Interaction: InteractionSelect},
	{Slug: "dock", Title: "Magnifying dock", Interaction: InteractionHover},
	{Slug: "progress", Title: "Progress ticks", Interaction: InteractionHover},
	{Slug: "humidity", Title: "Wheel", Interaction: InteractionDrag},
	{Slug: "command", Title: "Command bar", Interaction: InteractionType},
	{Slug: "roster", Title: "Selection list", Interaction: InteractionSelect},
	{Slug: "sleep", Title: "Range dial", Interaction: InteractionDrag},
	{Slug: "liq-toggle", Title: "Liquid toggle", Interaction: InteractionPress},
}

type registeredDemo struct {
	summary string
	demo    func(c *ui.Context)
}

var demoRegistry = map[string]registeredDemo{}

func registerDemo(slug, summary string, demo func(c *ui.Context)) {
	demoRegistry[slug] = registeredDemo{summary: summary, demo: demo}
}

// Catalog returns all bencho.dev blocks with live demos attached.
func Catalog() []Block {
	out := make([]Block, len(blockSpecs))
	for i, spec := range blockSpecs {
		if reg, ok := demoRegistry[spec.Slug]; ok {
			spec.Summary = reg.summary
			spec.Demo = reg.demo
		}
		out[i] = spec
	}
	return out
}

// FindBlock looks up a catalog block by slug.
func FindBlock(slug string) (Block, bool) {
	for _, b := range Catalog() {
		if b.Slug == slug {
			return b, true
		}
	}
	return Block{}, false
}

// ImplementedCount returns the number of blocks with live MyGo demos.
func ImplementedCount() int {
	n := 0
	for _, spec := range blockSpecs {
		if reg, ok := demoRegistry[spec.Slug]; ok && reg.demo != nil {
			n++
		}
	}
	return n
}
