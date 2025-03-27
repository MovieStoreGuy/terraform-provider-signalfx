package feature

const (
	PreviewNameProviderTags = "provider.tags"
)

var (
	_ = GetGlobalRegistry().MustRegister(
		PreviewNameProviderTags,
		WithPreviewDescription("When enabled, allows for tags to be set a provider that will apply to all applicable resources."),
		WithPreviewAddInVersion("v9.9.0"),
	)
)
