package config

import "testing"

const CONFIG = `{
	"timeout": 30,
	"settings_dir": "/media/NAS/Downloads/torrents",
	"sites": [
		{
			"name": "ttg",
			"rss": "https://totheglory.im/putrssmc.php?par=passkey",
			"download_dir": "/media/NAS/Downloads/torrents",
			"interval": 60
		},
		{
			"name": "hdc",
			"rss": "https://hdchina.club/torrentrss.php?rows=50&linktype=dl&passkey=passkey",
			"download_dir": "/media/NAS/Downloads/torrents",
			"interval": 60

		}
	]
}`

func TestNewConfig(t *testing.T) {
	c, err := NewConfig([]byte(CONFIG))
	if err != nil {
		t.Fatal("NewConfig failed")
	} else {
		t.Log("NewConfig success")
	}

	if c.Timeout != 30 {
		t.Fatal("timeout != 30")
	} else if c.SettingsDir != "/media/NAS/Downloads/torrents" {
		t.Fatal("settings_dir error")
	} else {
		t.Log("load timeout success")
	}

	if len(c.Sites) != 2 ||
		c.Sites[0].Name != "ttg" ||
		c.Sites[1].Rss != "https://hdchina.club/torrentrss.php?rows=50&linktype=dl&passkey=passkey" ||
		c.Sites[1].DownloadDir != "/media/NAS/Downloads/torrents" ||
		c.Sites[1].Interval != 60 {
		t.Fatal("load sites failed")
	} else {
		t.Log("load sites success")
	}
}
