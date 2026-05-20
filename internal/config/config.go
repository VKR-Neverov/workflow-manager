package config

const ()

var conf *Config

type Config struct {
}

func Init() error {
	//ENV := os.Getenv("ENV")
	//
	//body, err := os.ReadFile(fmt.Sprintf("./configs/values_%s.yaml", strings.ToLower(ENV)))
	//if err != nil {
	//	return err
	//}

	//err = yaml.Unmarshal(body, &conf)

	conf = &Config{}

	return nil
}

func Get(key string) interface{} {
	switch key {
	default:
		panic(ErrConfigNotFoundByKey(key))
	}
}
