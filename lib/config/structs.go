package config

type Config struct {
	Base `yaml:"-" json:"-"` // Ignore on marshalling

	Endpoint struct { // HTTP endpoint configuration
		Address     string `default:"0.0.0.0:8680"`
		TlsEnabled  bool   `default:"true"`     // If there is no certs - will be generated
		TlsCertPath string `default:"cert.pem"` // If relative - config dir path
		TlsKeyPath  string `default:"key.pem"`  // If relative - config dir path
	}
	Sources         map[string]Source   // List of sources
	Receivers       map[string]Receiver // List of receivers to trigger playback
	StateFilePath   string              `default:"nyanshare_state.json"`   // If relative - config dir path
	AccessFilePath  string              `default:"nyanshare_access.json"`  // If relative - config dir path
	RBACFilePath    string              `default:"nyanshare_rbac.json"`    // If relative - config dir path
	CatalogFilePath string              `default:"nyanshare_catalog.json"` // If relative - config dir path

	// Used to override the gui path, if set to empty - using embedded gui resources
	GuiPath string `default:""` // If relative - current working directory
}

type Source struct {
	Uri     string       `validate:"required"` // file://, http://, https://
	Type    string       `validate:"required"` // file, directory, syncthing, glob
	Options []OptionItem // Options depends on the source type // TODO: media, video, audio, photo, xml/others common remote
}
type Receiver struct {
	Uri     string       `validate:"required"` // Address of a receiver or name and params
	Type    string       `validate:"required"` // Some type
	Options []OptionItem // Options depends on the receiver type // TODO: subtitles, change audio stream...
}
type OptionItem struct {
	Key   string `validate:"required"`
	Value string `validate:"required"`
}

func (cfg Config) SourceSet(id string, src *Source) {
	cfg.Lock()
	cfg.Sources[id] = *src
	cfg.Unlock()
	cfg.Save()
}
