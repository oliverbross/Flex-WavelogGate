package config

import "testing"

func TestMigrateV9PreservesWavelogPairingAndSelectsXCAT(t *testing.T) {
	cfg := Config{
		Version:      9,
		Profile:      0,
		ProfileNames: []string{"Maestro", "Spare"},
		Profiles: []Profile{
			{
				WavelogURL:       "https://example.test/index.php",
				WavelogKey:       "secret-key",
				WavelogID:        "1",
				WavelogRadioname: "Maestro",
				HamlibEna:        true,
				HamlibHost:       "127.0.0.1",
				HamlibPort:       "5002",
				HamlibManaged:    true,
			},
			{},
		},
	}

	got := migrate(cfg)
	p := got.Profiles[0]
	if got.Version != 10 {
		t.Fatalf("version = %d, want 10", got.Version)
	}
	if p.WavelogURL != cfg.Profiles[0].WavelogURL ||
		p.WavelogKey != cfg.Profiles[0].WavelogKey ||
		p.WavelogID != cfg.Profiles[0].WavelogID ||
		p.WavelogRadioname != cfg.Profiles[0].WavelogRadioname {
		t.Fatalf("Wavelog pairing changed during migration: %+v", p)
	}
	if !p.XCATEna || p.XCATHost != "127.0.0.1" || p.XCATPort != "4532" {
		t.Fatalf("xCAT settings = enabled:%v %s:%s", p.XCATEna, p.XCATHost, p.XCATPort)
	}
	if p.FlrigEna || p.HamlibEna || p.HamlibManaged {
		t.Fatalf("legacy radio backend remains enabled: %+v", p)
	}
}

func TestDefaultUsesXCATOnly(t *testing.T) {
	cfg := Default()
	p := cfg.ActiveProfile()
	if !p.XCATEna || p.XCATHost != "127.0.0.1" || p.XCATPort != "4532" {
		t.Fatalf("unexpected xCAT defaults: %+v", p)
	}
	if p.FlrigEna || p.HamlibEna || p.HamlibManaged {
		t.Fatalf("legacy backend enabled by default: %+v", p)
	}
}
