package gocon

type Spec struct {
	AppName   string `json:"appName"`
	Namespace string `json:"namespace"`
	Port      string `json:"port"`
}

type Config struct {
	Spec    Spec   `json:"spec"`
	Content []byte `json:"content"`
}
