package generator

import (
	"fmt"
	"os"
)

func GenerateCaddyfile(config *Config) error {
	var content string

	if config.Domain != "" {
		content = fmt.Sprintf(`{
    auto_https on
}

%s {
    reverse_proxy /api* http://api:1337
    reverse_proxy /*    http://web:5173
    encode gzip
}

:80 {
    reverse_proxy /api* http://api:1337
    reverse_proxy /*    http://web:5173
}
`, config.Domain)
	} else {
		content = `{
    auto_https off
}

:80 {
    reverse_proxy /api* http://api:1337
    reverse_proxy /*    http://web:5173
}
`
	}

	return os.WriteFile("Caddyfile", []byte(content), 0644)
}
