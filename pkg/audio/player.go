package audio

import (
	"fmt"
	"os/exec"
	"runtime"
)

type Player struct {
	enabled bool
}

func NewPlayer(enabled bool) *Player {
	return &Player{enabled: enabled}
}

func (p *Player) Play(filepath string) error {
	if !p.enabled {
		return nil
	}

	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("afplay", filepath)
	case "linux":
		if _, err := exec.LookPath("paplay"); err == nil {
			cmd = exec.Command("paplay", filepath)
		} else if _, err := exec.LookPath("aplay"); err == nil {
			cmd = exec.Command("aplay", filepath)
		} else if _, err := exec.LookPath("mpv"); err == nil {
			cmd = exec.Command("mpv", "--no-video", filepath)
		} else {
			return fmt.Errorf("no audio player found on Linux")
		}
	case "windows":
		cmd = exec.Command("powershell", "-c", fmt.Sprintf("(New-Object Media.SoundPlayer '%s').PlaySync()", filepath))
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Run()
}
