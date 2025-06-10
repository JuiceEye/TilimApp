package assets

import (
	"embed"
	"encoding/base64"
	"fmt"
)

//go:embed audio/*.mp3
var audioFiles embed.FS

func LoadAudioByUUID(uuid string) (string, error) {
	path := fmt.Sprintf("audio/%s.mp3", uuid)
	data, err := audioFiles.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to load audio %s: %w", uuid, err)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}
