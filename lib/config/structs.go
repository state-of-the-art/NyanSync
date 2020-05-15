package config

type Config struct {
	Base `yaml:",omitempty"`

	Endpoint struct { // HTTP endpoint configuration
		Address     string `cfgDefault:"0.0.0.0:8680"`
		TlsEnabled  bool   `cfgDefault:"true"`     // If there is no certs - will be generated
		TlsCertPath string `cfgDefault:"cert.pem"` // If relative - config dir path
		TlsKeyPath  string `cfgDefault:"key.pem"`  // If relative - config dir path
	}
	Sources         map[string]Source   // List of sources
	Receivers       map[string]Receiver // List of receivers to trigger playback
	StateFilePath   string              `cfgDefault:"nyanshare_state.json"`   // If relative - config dir path
	AccessFilePath  string              `cfgDefault:"nyanshare_access.json"`  // If relative - config dir path
	CatalogFilePath string              `cfgDefault:"nyanshare_catalog.json"` // If relative - config dir path

	// Used to override the gui path, if set to empty - using embedded gui resources
	GuiPath string `cfgDefault:""` // If relative - current working directory
}

type Source struct {
	Uri     string       `cfgRequired:"true"` // file://, http://, https://
	Type    string       `cfgRequired:"true"` // file, directory, syncthing, glob
	Options []OptionItem // Options depends on the source type // TODO: media, video, audio, photo, xml/others common remote
}
type Receiver struct {
	Uri     string       `cfgRequired:"true"` // Address of a receiver or name and params
	Type    string       `cfgRequired:"true"` // Some type
	Options []OptionItem // Options depends on the receiver type // TODO: subtitles, change audio stream...
}
type OptionItem struct {
	Key   string `cfgRequired:"true"`
	Value string `cfgRequired:"true"`
}

func (cfg Config) SourceSet(id string, src *Source) {
	cfg.Lock()
	cfg.Sources[id] = *src
	cfg.Unlock()
	cfg.Save()
}
