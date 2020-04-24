package config

type Config struct {
	Endpoint struct { // HTTP endpoint configuration
		Address     string `cfgDefault:"0.0.0.0:8680"`
		TlsEnabled  bool   `cfgDefault:"true"`     // If there is no certs - will be generated
		TlsCertPath string `cfgDefault:"cert.pem"` // If relative - config dir path
		TlsKeyPath  string `cfgDefault:"key.pem"`  // If relative - config dir path
	}
	Sources []Source
	Receivers []struct { // List of supported receivers to trigger playback
		Id      string       `cfgRequired:"true"` // Just an unique name
		Url     string       `cfgRequired:"true"` // Address of a receiver or name and params
		Options []OptionItem // Options depends on the receiver type // TODO: subtitles, change audio stream...
	}
	StateFilePath   string `cfgDefault:"nyansync_state.json"`   // If relative - config dir path
	AccessFilePath  string `cfgDefault:"nyansync_access.json"`  // If relative - config dir path
	CatalogFilePath string `cfgDefault:"nyansync_catalog.json"` // If relative - config dir path

	// Used to override the gui path, if set to empty - using embedded gui resources
	GuiPath string `cfgDefault:""` // If relative - current working directory
}

type Source struct { // List of sources
	Id      string       `cfgRequired:"true"` // Just an unique name
	Url     string       `cfgRequired:"true"` // file://, http://, https://
	Type    string       `cfgRequired:"true"` // file, directory, syncthing, glob
	Options []OptionItem // Options depends on the source type // TODO: media, video, audio, photo, xml/others common remote
}
type OptionItem struct {
	Key   string `cfgRequired:"true"`
	Value string `cfgRequired:"true"`
}
