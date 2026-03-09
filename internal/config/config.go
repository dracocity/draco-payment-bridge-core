package config

type ListenConfig struct {
	Network string `toml:"network"`
	Address string `toml:"address"`
}

type Config struct {
	Port      string              `toml:"port" envconfig:"PORT" default:"8080"`
	Listen    []ListenConfig      `toml:"listen"`
	PluginDir string              `toml:"plugin_dir" envconfig:"PLUGIN_DIR" default:"./plugins"`
	Providers map[string]PGConfig `toml:"providers"`
	Log       LogConfig           `toml:"log"`
}

// TODO checkoutMode *string `json:"checkout_mode,omitempty"` // HOSTED | SELF_HOSTED
type PGConfig map[string]string

type LogConfig struct {
	Level    string             `toml:"level" envconfig:"LOG_LEVEL" default:"info"`   // debug, info, warn, error
	Format   string             `toml:"format" envconfig:"LOG_FORMAT" default:"json"` // console or json
	Output   LogOutputConfig    `toml:"output"`
	Rotation *LogRotationConfig `toml:"rotation"`
}

type LogOutputConfig struct {
	Stdout   bool   `toml:"stdout" envconfig:"LOG_STDOUT" default:"true"`  // stdout
	Stderr   bool   `toml:"stderr" envconfig:"LOG_STDERR" default:"false"` // stderr
	FilePath string `toml:"file" envconfig:"LOG_FILE"`                     // file path
}

type LogRotationConfig struct {
	Enabled  bool `toml:"enabled" envconfig:"LOG_ROTATION_ENABLED" default:"false"`
	MaxSize  int  `toml:"max_size" envconfig:"LOG_MAX_SIZE" default:"100"`  // MB unit
	MaxAge   int  `toml:"max_age" envconfig:"LOG_MAX_AGE" default:"7"`      // day
	Compress bool `toml:"compress" envconfig:"LOG_COMPRESS" default:"true"` // gzip compression
}
