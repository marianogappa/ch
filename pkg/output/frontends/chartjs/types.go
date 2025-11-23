package chartjs

import "time"

// ChartType represents the type of chart to render
type ChartType string

const (
	ChartTypeLine      ChartType = "line"
	ChartTypeBar       ChartType = "bar"
	ChartTypePie       ChartType = "pie"
	ChartTypeDoughnut  ChartType = "doughnut"
	ChartTypePolarArea ChartType = "polarArea"
	ChartTypeRadar     ChartType = "radar"
	ChartTypeBubble    ChartType = "bubble"
	ChartTypeScatter   ChartType = "scatter"
)

// Chart represents the complete Chart.js configuration
type Chart struct {
	Type    ChartType `json:"type"`
	Data    *Data     `json:"data"`
	Options *Options  `json:"options,omitempty"`
}

// Data holds the datasets and labels for the chart
type Data struct {
	Labels   []any     `json:"labels,omitempty"` // Can be strings, numbers, or time.Time
	Datasets []Dataset `json:"datasets"`
}

// Dataset represents a single dataset in the chart
type Dataset struct {
	// Common properties
	Label           string  `json:"label,omitempty"`
	Data            []any   `json:"data"`                      // Can be numbers, {x,y} objects, {x,y,r} objects, etc.
	BackgroundColor any     `json:"backgroundColor,omitempty"` // string, []string, or function
	BorderColor     any     `json:"borderColor,omitempty"`     // string, []string, or function
	BorderWidth     *int    `json:"borderWidth,omitempty"`
	BorderSkipped   *bool   `json:"borderSkipped,omitempty"` // For bar charts
	Order           *int    `json:"order,omitempty"`
	Stack           *string `json:"stack,omitempty"`
	IndexAxis       *string `json:"indexAxis,omitempty"` // 'x' or 'y'

	// Line chart specific
	Fill                 *bool    `json:"fill,omitempty"`
	Tension              *float64 `json:"tension,omitempty"` // Bezier curve tension (0-1)
	PointRadius          *float64 `json:"pointRadius,omitempty"`
	PointHoverRadius     *float64 `json:"pointHoverRadius,omitempty"`
	PointBackgroundColor *string  `json:"pointBackgroundColor,omitempty"`
	PointBorderColor     *string  `json:"pointBorderColor,omitempty"`
	PointStyle           *string  `json:"pointStyle,omitempty"` // 'circle', 'cross', 'crossRot', 'dash', 'line', 'rect', 'rectRounded', 'rectRot', 'star', 'triangle'
	ShowLine             *bool    `json:"showLine,omitempty"`
	SpanGaps             *bool    `json:"spanGaps,omitempty"`
	Stepped              *bool    `json:"stepped,omitempty"` // 'before', 'after', 'middle', false

	// Bar chart specific
	BarPercentage      *float64 `json:"barPercentage,omitempty"`
	CategoryPercentage *float64 `json:"categoryPercentage,omitempty"`
	BarThickness       *string  `json:"barThickness,omitempty"` // 'flex' or number

	// Pie/Doughnut specific
	Offset *float64 `json:"offset,omitempty"`
	Weight *float64 `json:"weight,omitempty"`

	// Bubble/Scatter specific
	Radius *float64 `json:"radius,omitempty"`

	// Radar specific
	LineTension *float64 `json:"lineTension,omitempty"`
}

// Options holds all chart configuration options
type Options struct {
	// Core options
	Responsive          *bool    `json:"responsive,omitempty"`
	MaintainAspectRatio *bool    `json:"maintainAspectRatio,omitempty"`
	AspectRatio         *float64 `json:"aspectRatio,omitempty"`
	ResizeDelay         *int     `json:"resizeDelay,omitempty"`

	// Interaction
	Interaction *Interaction `json:"interaction,omitempty"`
	OnHover     *string      `json:"onHover,omitempty"` // JavaScript function as string
	OnClick     *string      `json:"onClick,omitempty"` // JavaScript function as string

	// Layout
	Layout *Layout `json:"layout,omitempty"`

	// Animation
	Animation  *Animation  `json:"animation,omitempty"`
	Animations *Animations `json:"animations,omitempty"`
	Transition *Transition `json:"transition,omitempty"`

	// Scales
	Scales *Scales `json:"scales,omitempty"`

	// Plugins
	Plugins *Plugins `json:"plugins,omitempty"`

	// Elements
	Elements *Elements `json:"elements,omitempty"`

	// Device pixel ratio
	DevicePixelRatio *float64 `json:"devicePixelRatio,omitempty"`

	// Index axis (for bar charts)
	IndexAxis *string `json:"indexAxis,omitempty"` // 'x' or 'y'

	// Parsing
	Parsing *Parsing `json:"parsing,omitempty"`
}

// Interaction configures chart interaction behavior
type Interaction struct {
	Intersect        *bool   `json:"intersect,omitempty"`
	Mode             *string `json:"mode,omitempty"` // 'point', 'nearest', 'index', 'dataset', 'x', 'y'
	Axis             *string `json:"axis,omitempty"` // 'x', 'y', 'xy', 'r'
	IncludeInvisible *bool   `json:"includeInvisible,omitempty"`
}

// Layout configures padding
type Layout struct {
	Padding *Padding `json:"padding,omitempty"`
}

// Padding configures chart padding
type Padding struct {
	Top    *int `json:"top,omitempty"`
	Right  *int `json:"right,omitempty"`
	Bottom *int `json:"bottom,omitempty"`
	Left   *int `json:"left,omitempty"`
	X      *int `json:"x,omitempty"`
	Y      *int `json:"y,omitempty"`
}

// Animation configures chart animations
type Animation struct {
	Duration *int    `json:"duration,omitempty"`
	Easing   *string `json:"easing,omitempty"` // 'linear', 'easeInQuad', 'easeOutQuad', etc.
	Delay    *int    `json:"delay,omitempty"`
	Loop     *bool   `json:"loop,omitempty"`
}

// Animations configures multiple animations
type Animations struct {
	Colors *Animation `json:"colors,omitempty"`
	X      *Animation `json:"x,omitempty"`
	Y      *Animation `json:"y,omitempty"`
}

// Transition configures chart transitions
type Transition struct {
	Animation *Animation `json:"animation,omitempty"`
}

// Scales configures chart axes/scales
type Scales struct {
	X  *Scale `json:"x,omitempty"`
	Y  *Scale `json:"y,omitempty"`
	R  *Scale `json:"r,omitempty"` // For radar/polar area
	X1 *Scale `json:"x1,omitempty"`
	Y1 *Scale `json:"y1,omitempty"`
}

// Scale represents a single scale/axis configuration
type Scale struct {
	Type        *string `json:"type,omitempty"` // 'category', 'linear', 'logarithmic', 'time', 'timeseries', 'radialLinear'
	Display     *bool   `json:"display,omitempty"`
	Position    *string `json:"position,omitempty"` // 'top', 'bottom', 'left', 'right', 'center'
	Offset      *bool   `json:"offset,omitempty"`
	Reverse     *bool   `json:"reverse,omitempty"`
	BeginAtZero *bool   `json:"beginAtZero,omitempty"`

	// Title
	Title *ScaleTitle `json:"title,omitempty"`

	// Ticks
	Ticks *Ticks `json:"ticks,omitempty"`

	// Grid
	Grid *Grid `json:"grid,omitempty"`

	// Border
	Border *Border `json:"border,omitempty"`

	// Time scale specific
	Time *TimeScale `json:"time,omitempty"`

	// Logarithmic scale specific
	Min *float64 `json:"min,omitempty"`
	Max *float64 `json:"max,omitempty"`

	// Stacked
	Stacked *bool   `json:"stacked,omitempty"`
	Stack   *string `json:"stack,omitempty"`

	// Weight
	Weight *float64 `json:"weight,omitempty"`
}

// ScaleTitle configures scale title
type ScaleTitle struct {
	Display *bool    `json:"display,omitempty"`
	Text    *string  `json:"text,omitempty"`
	Color   *string  `json:"color,omitempty"`
	Font    *Font    `json:"font,omitempty"`
	Padding *Padding `json:"padding,omitempty"`
	Align   *string  `json:"align,omitempty"` // 'start', 'center', 'end'
}

// Ticks configures scale ticks
type Ticks struct {
	Display           *bool    `json:"display,omitempty"`
	Color             *string  `json:"color,omitempty"`
	Font              *Font    `json:"font,omitempty"`
	Padding           *int     `json:"padding,omitempty"`
	StepSize          *float64 `json:"stepSize,omitempty"`
	MaxRotation       *float64 `json:"maxRotation,omitempty"`
	MinRotation       *float64 `json:"minRotation,omitempty"`
	Mirror            *bool    `json:"mirror,omitempty"`
	TextStrokeColor   *string  `json:"textStrokeColor,omitempty"`
	TextStrokeWidth   *float64 `json:"textStrokeWidth,omitempty"`
	Z                 *float64 `json:"z,omitempty"`
	Callback          *string  `json:"callback,omitempty"` // JavaScript function as string
	AutoSkip          *bool    `json:"autoSkip,omitempty"`
	AutoSkipPadding   *float64 `json:"autoSkipPadding,omitempty"`
	IncludeBounds     *bool    `json:"includeBounds,omitempty"`
	LabelOffset       *float64 `json:"labelOffset,omitempty"`
	MaxTicksLimit     *int     `json:"maxTicksLimit,omitempty"`
	Precision         *int     `json:"precision,omitempty"`
	ShowLabelBackdrop *bool    `json:"showLabelBackdrop,omitempty"`
	Source            *string  `json:"source,omitempty"` // 'auto', 'data', 'labels'
}

// Grid configures grid lines
type Grid struct {
	Display          *bool     `json:"display,omitempty"`
	Color            *string   `json:"color,omitempty"`
	LineWidth        *float64  `json:"lineWidth,omitempty"`
	DrawOnChartArea  *bool     `json:"drawOnChartArea,omitempty"`
	DrawTicks        *bool     `json:"drawTicks,omitempty"`
	TickLength       *float64  `json:"tickLength,omitempty"`
	Offset           *bool     `json:"offset,omitempty"`
	BorderDash       []float64 `json:"borderDash,omitempty"`
	BorderDashOffset *float64  `json:"borderDashOffset,omitempty"`
	Z                *float64  `json:"z,omitempty"`
}

// Border configures scale border
type Border struct {
	Display *bool    `json:"display,omitempty"`
	Color   *string  `json:"color,omitempty"`
	Width   *float64 `json:"width,omitempty"`
	Z       *float64 `json:"z,omitempty"`
}

// TimeScale configures time scale options
type TimeScale struct {
	DisplayFormats *map[string]string `json:"displayFormats,omitempty"`
	IsoWeekday     *bool              `json:"isoWeekday,omitempty"`
	Parser         *string            `json:"parser,omitempty"` // Date format string
	Round          *string            `json:"round,omitempty"`  // 'millisecond', 'second', 'minute', 'hour', 'day', 'week', 'month', 'quarter', 'year'
	TooltipFormat  *string            `json:"tooltipFormat,omitempty"`
	Unit           *string            `json:"unit,omitempty"` // 'millisecond', 'second', 'minute', 'hour', 'day', 'week', 'month', 'quarter', 'year'
	UnitStepSize   *float64           `json:"unitStepSize,omitempty"`
	MinUnit        *string            `json:"minUnit,omitempty"`
	Adaptive       *bool              `json:"adaptive,omitempty"`
}

// Font configures font properties
type Font struct {
	Family     *string  `json:"family,omitempty"`
	Size       *int     `json:"size,omitempty"`
	Style      *string  `json:"style,omitempty"`  // 'normal', 'italic', 'oblique'
	Weight     *string  `json:"weight,omitempty"` // 'normal', 'bold', 'lighter', 'bolder', or number
	LineHeight *float64 `json:"lineHeight,omitempty"`
}

// Plugins configures chart plugins
type Plugins struct {
	Legend     *Legend        `json:"legend,omitempty"`
	Title      *Title         `json:"title,omitempty"`
	Tooltip    *Tooltip       `json:"tooltip,omitempty"`
	Subtitle   *Subtitle      `json:"subtitle,omitempty"`
	Filler     *Filler        `json:"filler,omitempty"`
	Decimation *Decimation    `json:"decimation,omitempty"`
	Zoom       *Zoom          `json:"zoom,omitempty"`
	Custom     map[string]any `json:"-"` // For custom plugins
}

// Legend configures the chart legend
type Legend struct {
	Display  *bool         `json:"display,omitempty"`
	Position *string       `json:"position,omitempty"` // 'top', 'bottom', 'left', 'right', 'chartArea'
	Align    *string       `json:"align,omitempty"`    // 'start', 'center', 'end'
	FullSize *bool         `json:"fullSize,omitempty"`
	Reverse  *bool         `json:"reverse,omitempty"`
	Rtl      *bool         `json:"rtl,omitempty"`
	Text     *LegendText   `json:"text,omitempty"`
	Title    *LegendTitle  `json:"title,omitempty"`
	Labels   *LegendLabels `json:"labels,omitempty"`
	OnClick  *string       `json:"onClick,omitempty"` // JavaScript function
	OnHover  *string       `json:"onHover,omitempty"` // JavaScript function
	OnLeave  *string       `json:"onLeave,omitempty"` // JavaScript function
}

// LegendText configures legend text
type LegendText struct {
	Color         *string `json:"color,omitempty"`
	Font          *Font   `json:"font,omitempty"`
	Padding       *int    `json:"padding,omitempty"`
	UsePointStyle *bool   `json:"usePointStyle,omitempty"`
	PointStyle    *string `json:"pointStyle,omitempty"`
}

// LegendTitle configures legend title
type LegendTitle struct {
	Color   *string `json:"color,omitempty"`
	Font    *Font   `json:"font,omitempty"`
	Padding *int    `json:"padding,omitempty"`
	Text    *string `json:"text,omitempty"`
}

// LegendLabels configures legend labels
type LegendLabels struct {
	BoxWidth       *float64 `json:"boxWidth,omitempty"`
	BoxHeight      *float64 `json:"boxHeight,omitempty"`
	Padding        *int     `json:"padding,omitempty"`
	Color          *string  `json:"color,omitempty"`
	Font           *Font    `json:"font,omitempty"`
	TextAlign      *string  `json:"textAlign,omitempty"` // 'left', 'center', 'right'
	UsePointStyle  *bool    `json:"usePointStyle,omitempty"`
	PointStyle     *string  `json:"pointStyle,omitempty"`
	GenerateLabels *string  `json:"generateLabels,omitempty"` // JavaScript function
	Filter         *string  `json:"filter,omitempty"`         // JavaScript function
	Sort           *string  `json:"sort,omitempty"`           // JavaScript function
}

// Title configures the chart title
type Title struct {
	Display  *bool    `json:"display,omitempty"`
	Position *string  `json:"position,omitempty"` // 'top', 'bottom', 'left', 'right'
	Align    *string  `json:"align,omitempty"`    // 'start', 'center', 'end'
	Color    *string  `json:"color,omitempty"`
	Font     *Font    `json:"font,omitempty"`
	Padding  *Padding `json:"padding,omitempty"`
	Text     *string  `json:"text,omitempty"`
	FullSize *bool    `json:"fullSize,omitempty"`
}

// Subtitle configures the chart subtitle
type Subtitle struct {
	Display  *bool    `json:"display,omitempty"`
	Position *string  `json:"position,omitempty"`
	Align    *string  `json:"align,omitempty"`
	Color    *string  `json:"color,omitempty"`
	Font     *Font    `json:"font,omitempty"`
	Padding  *Padding `json:"padding,omitempty"`
	Text     *string  `json:"text,omitempty"`
}

// Tooltip configures tooltips
type Tooltip struct {
	Enabled           *bool             `json:"enabled,omitempty"`
	External          *string           `json:"external,omitempty"` // JavaScript function
	Position          *string           `json:"position,omitempty"` // 'average', 'nearest', or function
	BackgroundColor   *string           `json:"backgroundColor,omitempty"`
	TitleColor        *string           `json:"titleColor,omitempty"`
	TitleFont         *Font             `json:"titleFont,omitempty"`
	TitleAlign        *string           `json:"titleAlign,omitempty"` // 'left', 'center', 'right'
	TitleSpacing      *float64          `json:"titleSpacing,omitempty"`
	TitleMarginBottom *float64          `json:"titleMarginBottom,omitempty"`
	BodyColor         *string           `json:"bodyColor,omitempty"`
	BodyFont          *Font             `json:"bodyFont,omitempty"`
	BodyAlign         *string           `json:"bodyAlign,omitempty"`
	BodySpacing       *float64          `json:"bodySpacing,omitempty"`
	FooterColor       *string           `json:"footerColor,omitempty"`
	FooterFont        *Font             `json:"footerFont,omitempty"`
	FooterAlign       *string           `json:"footerAlign,omitempty"`
	FooterSpacing     *float64          `json:"footerSpacing,omitempty"`
	FooterMarginTop   *float64          `json:"footerMarginTop,omitempty"`
	Padding           *Padding          `json:"padding,omitempty"`
	PaddingX          *int              `json:"paddingX,omitempty"`
	PaddingY          *int              `json:"paddingY,omitempty"`
	CaretPadding      *float64          `json:"caretPadding,omitempty"`
	CaretSize         *float64          `json:"caretSize,omitempty"`
	CornerRadius      *float64          `json:"cornerRadius,omitempty"`
	DisplayColors     *bool             `json:"displayColors,omitempty"`
	BoxWidth          *float64          `json:"boxWidth,omitempty"`
	BoxHeight         *float64          `json:"boxHeight,omitempty"`
	BoxPadding        *float64          `json:"boxPadding,omitempty"`
	UsePointStyle     *bool             `json:"usePointStyle,omitempty"`
	BorderColor       *string           `json:"borderColor,omitempty"`
	BorderWidth       *float64          `json:"borderWidth,omitempty"`
	Rtl               *bool             `json:"rtl,omitempty"`
	TextDirection     *string           `json:"textDirection,omitempty"` // 'ltr', 'rtl'
	XAlign            *string           `json:"xAlign,omitempty"`        // 'left', 'center', 'right'
	YAlign            *string           `json:"yAlign,omitempty"`        // 'top', 'center', 'bottom'
	Filter            *string           `json:"filter,omitempty"`        // JavaScript function
	ItemSort          *string           `json:"itemSort,omitempty"`      // JavaScript function
	Callbacks         *TooltipCallbacks `json:"callbacks,omitempty"`
	Intersect         *bool             `json:"intersect,omitempty"`
	Mode              *string           `json:"mode,omitempty"` // 'point', 'nearest', 'index', 'dataset', 'x', 'y'
	Axis              *string           `json:"axis,omitempty"` // 'x', 'y', 'xy', 'r'
}

// TooltipCallbacks configures tooltip callbacks
type TooltipCallbacks struct {
	BeforeTitle    *string `json:"beforeTitle,omitempty"`    // JavaScript function
	Title          *string `json:"title,omitempty"`          // JavaScript function
	AfterTitle     *string `json:"afterTitle,omitempty"`     // JavaScript function
	BeforeBody     *string `json:"beforeBody,omitempty"`     // JavaScript function
	BeforeLabel    *string `json:"beforeLabel,omitempty"`    // JavaScript function
	Label          *string `json:"label,omitempty"`          // JavaScript function
	LabelColor     *string `json:"labelColor,omitempty"`     // JavaScript function
	LabelTextColor *string `json:"labelTextColor,omitempty"` // JavaScript function
	AfterLabel     *string `json:"afterLabel,omitempty"`     // JavaScript function
	AfterBody      *string `json:"afterBody,omitempty"`      // JavaScript function
	BeforeFooter   *string `json:"beforeFooter,omitempty"`   // JavaScript function
	Footer         *string `json:"footer,omitempty"`         // JavaScript function
	AfterFooter    *string `json:"afterFooter,omitempty"`    // JavaScript function
}

// Filler configures area fill
type Filler struct {
	Propagate *bool `json:"propagate,omitempty"`
}

// Decimation configures data decimation
type Decimation struct {
	Enabled   *bool   `json:"enabled,omitempty"`
	Algorithm *string `json:"algorithm,omitempty"` // 'lttb', 'min-max'
	Samples   *int    `json:"samples,omitempty"`
	Threshold *int    `json:"threshold,omitempty"`
}

// Zoom configures zoom plugin
type Zoom struct {
	Zoom *ZoomOptions `json:"zoom,omitempty"`
	Pan  *PanOptions  `json:"pan,omitempty"`
}

// ZoomOptions configures zoom behavior
type ZoomOptions struct {
	Wheel *WheelOptions `json:"wheel,omitempty"`
	Pinch *PinchOptions `json:"pinch,omitempty"`
	Drag  *DragOptions  `json:"drag,omitempty"`
	Mode  *string       `json:"mode,omitempty"` // 'x', 'y', 'xy'
}

// PanOptions configures pan behavior
type PanOptions struct {
	Enabled     *bool   `json:"enabled,omitempty"`
	Mode        *string `json:"mode,omitempty"`        // 'x', 'y', 'xy'
	ModifierKey *string `json:"modifierKey,omitempty"` // 'ctrl', 'alt', 'shift', 'meta'
}

// WheelOptions configures wheel zoom
type WheelOptions struct {
	Enabled     *bool    `json:"enabled,omitempty"`
	Speed       *float64 `json:"speed,omitempty"`
	ModifierKey *string  `json:"modifierKey,omitempty"`
}

// PinchOptions configures pinch zoom
type PinchOptions struct {
	Enabled *bool `json:"enabled,omitempty"`
}

// DragOptions configures drag zoom
type DragOptions struct {
	Enabled     *bool    `json:"enabled,omitempty"`
	ModifierKey *string  `json:"modifierKey,omitempty"`
	Threshold   *float64 `json:"threshold,omitempty"`
}

// Elements configures chart elements
type Elements struct {
	Point *PointElement `json:"point,omitempty"`
	Line  *LineElement  `json:"line,omitempty"`
	Bar   *BarElement   `json:"bar,omitempty"`
	Arc   *ArcElement   `json:"arc,omitempty"`
}

// PointElement configures point elements
type PointElement struct {
	Radius           *float64 `json:"radius,omitempty"`
	PointStyle       *string  `json:"pointStyle,omitempty"`
	Rotation         *float64 `json:"rotation,omitempty"`
	BackgroundColor  *string  `json:"backgroundColor,omitempty"`
	BorderWidth      *float64 `json:"borderWidth,omitempty"`
	BorderColor      *string  `json:"borderColor,omitempty"`
	HitRadius        *float64 `json:"hitRadius,omitempty"`
	HoverRadius      *float64 `json:"hoverRadius,omitempty"`
	HoverBorderWidth *float64 `json:"hoverBorderWidth,omitempty"`
}

// LineElement configures line elements
type LineElement struct {
	Tension          *float64  `json:"tension,omitempty"`
	BackgroundColor  *string   `json:"backgroundColor,omitempty"`
	BorderWidth      *float64  `json:"borderWidth,omitempty"`
	BorderColor      *string   `json:"borderColor,omitempty"`
	BorderCapStyle   *string   `json:"borderCapStyle,omitempty"` // 'butt', 'round', 'square'
	BorderDash       []float64 `json:"borderDash,omitempty"`
	BorderDashOffset *float64  `json:"borderDashOffset,omitempty"`
	BorderJoinStyle  *string   `json:"borderJoinStyle,omitempty"` // 'bevel', 'round', 'miter'
	CapBezierPoints  *bool     `json:"capBezierPoints,omitempty"`
	Fill             *bool     `json:"fill,omitempty"`
	Stepped          *bool     `json:"stepped,omitempty"`
}

// BarElement configures bar elements
type BarElement struct {
	BackgroundColor *string  `json:"backgroundColor,omitempty"`
	BorderWidth     *float64 `json:"borderWidth,omitempty"`
	BorderColor     *string  `json:"borderColor,omitempty"`
	BorderSkipped   *bool    `json:"borderSkipped,omitempty"`
	BorderRadius    *float64 `json:"borderRadius,omitempty"`
	InflateAmount   *string  `json:"inflateAmount,omitempty"` // 'auto' or number
}

// ArcElement configures arc elements (for pie/doughnut)
type ArcElement struct {
	Angle           *float64 `json:"angle,omitempty"`
	BackgroundColor *string  `json:"backgroundColor,omitempty"`
	BorderWidth     *float64 `json:"borderWidth,omitempty"`
	BorderColor     *string  `json:"borderColor,omitempty"`
	BorderAlign     *string  `json:"borderAlign,omitempty"` // 'center', 'inner'
	Circular        *bool    `json:"circular,omitempty"`
	Offset          *float64 `json:"offset,omitempty"`
	Spacing         *float64 `json:"spacing,omitempty"`
	Weight          *float64 `json:"weight,omitempty"`
}

// Parsing configures data parsing
type Parsing struct {
	XAxisKey *string `json:"xAxisKey,omitempty"`
	YAxisKey *string `json:"yAxisKey,omitempty"`
}

// DataPoint represents a single data point (for scatter/bubble charts)
type DataPoint struct {
	X float64  `json:"x"`
	Y float64  `json:"y"`
	R *float64 `json:"r,omitempty"` // For bubble charts
}

// TimeDataPoint represents a time-based data point
type TimeDataPoint struct {
	X time.Time `json:"x"`
	Y float64   `json:"y"`
	R *float64  `json:"r,omitempty"`
}
