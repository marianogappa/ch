package chartjs

// Helper functions for common chart configurations

// WithTitle adds a title to chart options
func WithTitle(text string) func(*Options) {
	return func(opts *Options) {
		if opts.Plugins == nil {
			opts.Plugins = &Plugins{}
		}
		opts.Plugins.Title = &Title{
			Display: boolPtr(true),
			Text:    stringPtr(text),
		}
	}
}

// WithSubtitle adds a subtitle to chart options
func WithSubtitle(text string) func(*Options) {
	return func(opts *Options) {
		if opts.Plugins == nil {
			opts.Plugins = &Plugins{}
		}
		opts.Plugins.Subtitle = &Subtitle{
			Display: boolPtr(true),
			Text:    stringPtr(text),
		}
	}
}

// WithLegend configures the legend
func WithLegend(display bool, position string) func(*Options) {
	return func(opts *Options) {
		if opts.Plugins == nil {
			opts.Plugins = &Plugins{}
		}
		opts.Plugins.Legend = &Legend{
			Display:  boolPtr(display),
			Position: stringPtr(position),
		}
	}
}

// WithXAxisLabel adds an X-axis label
func WithXAxisLabel(label string) func(*Options) {
	return func(opts *Options) {
		if opts.Scales == nil {
			opts.Scales = &Scales{}
		}
		if opts.Scales.X == nil {
			opts.Scales.X = &Scale{}
		}
		opts.Scales.X.Title = &ScaleTitle{
			Display: boolPtr(true),
			Text:    stringPtr(label),
		}
	}
}

// WithYAxisLabel adds a Y-axis label
func WithYAxisLabel(label string) func(*Options) {
	return func(opts *Options) {
		if opts.Scales == nil {
			opts.Scales = &Scales{}
		}
		if opts.Scales.Y == nil {
			opts.Scales.Y = &Scale{}
		}
		opts.Scales.Y.Title = &ScaleTitle{
			Display: boolPtr(true),
			Text:    stringPtr(label),
		}
	}
}

// WithYAxisBeginAtZero sets Y-axis to begin at zero
func WithYAxisBeginAtZero(beginAtZero bool) func(*Options) {
	return func(opts *Options) {
		if opts.Scales == nil {
			opts.Scales = &Scales{}
		}
		if opts.Scales.Y == nil {
			opts.Scales.Y = &Scale{}
		}
		opts.Scales.Y.BeginAtZero = boolPtr(beginAtZero)
	}
}

// WithLogarithmicScale sets Y-axis to logarithmic scale
func WithLogarithmicScale() func(*Options) {
	return func(opts *Options) {
		if opts.Scales == nil {
			opts.Scales = &Scales{}
		}
		if opts.Scales.Y == nil {
			opts.Scales.Y = &Scale{}
		}
		opts.Scales.Y.Type = stringPtr("logarithmic")
	}
}

// WithTimeScale configures X-axis as time scale
func WithTimeScale(unit string) func(*Options) {
	return func(opts *Options) {
		if opts.Scales == nil {
			opts.Scales = &Scales{}
		}
		if opts.Scales.X == nil {
			opts.Scales.X = &Scale{}
		}
		opts.Scales.X.Type = stringPtr("time")
		if unit != "" {
			opts.Scales.X.Time = &TimeScale{
				Unit: stringPtr(unit),
			}
		}
	}
}

// WithResponsive sets responsive option
func WithResponsive(responsive bool) func(*Options) {
	return func(opts *Options) {
		opts.Responsive = boolPtr(responsive)
	}
}

// WithAspectRatio sets aspect ratio
func WithAspectRatio(ratio float64) func(*Options) {
	return func(opts *Options) {
		opts.AspectRatio = float64Ptr(ratio)
	}
}

// WithAnimation configures animation
func WithAnimation(duration int, easing string) func(*Options) {
	return func(opts *Options) {
		opts.Animation = &Animation{
			Duration: intPtr(duration),
			Easing:   stringPtr(easing),
		}
	}
}

// WithTooltipCallback sets a custom tooltip label callback
func WithTooltipCallback(callback string) func(*Options) {
	return func(opts *Options) {
		if opts.Plugins == nil {
			opts.Plugins = &Plugins{}
		}
		if opts.Plugins.Tooltip == nil {
			opts.Plugins.Tooltip = &Tooltip{}
		}
		if opts.Plugins.Tooltip.Callbacks == nil {
			opts.Plugins.Tooltip.Callbacks = &TooltipCallbacks{}
		}
		opts.Plugins.Tooltip.Callbacks.Label = stringPtr(callback)
	}
}

// WithStacked sets stacked mode for bar/line charts
func WithStacked(stacked bool) func(*Options) {
	return func(opts *Options) {
		if opts.Scales == nil {
			opts.Scales = &Scales{}
		}
		if opts.Scales.Y == nil {
			opts.Scales.Y = &Scale{}
		}
		opts.Scales.Y.Stacked = boolPtr(stacked)
	}
}

// ApplyOptions applies multiple option functions to Options
func ApplyOptions(opts *Options, funcs ...func(*Options)) {
	for _, f := range funcs {
		f(opts)
	}
}

// Helper pointer functions
func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
