package a11y

import (
	"testing"
)

// envGetter devolve um getenv baseado em um mapa. Chaves ausentes valem "".
func envGetter(env map[string]string) func(string) string {
	return func(key string) string { return env[key] }
}

func TestDetectFromMatrix(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want Profile
	}{
		{
			name: "ambiente vazio",
			env:  map[string]string{},
			want: Profile{Unicode: true},
		},
		{
			name: "NO_COLOR vazio não desliga cor",
			env:  map[string]string{"NO_COLOR": ""},
			want: Profile{Unicode: true},
		},
		{
			name: "NO_COLOR=1",
			env:  map[string]string{"NO_COLOR": "1"},
			want: Profile{NoColor: true},
		},
		{
			name: "NO_COLOR=0 também desliga cor",
			env:  map[string]string{"NO_COLOR": "0"},
			want: Profile{NoColor: true},
		},
		{
			name: "CLICOLOR=0",
			env:  map[string]string{"CLICOLOR": "0"},
			want: Profile{NoColor: true},
		},
		{
			name: "CLICOLOR=1 não altera",
			env:  map[string]string{"CLICOLOR": "1"},
			want: Profile{Unicode: true},
		},
		{
			name: "TERM=dumb",
			env:  map[string]string{"TERM": "dumb"},
			want: Profile{NoColor: true, Dumb: true, ScreenReader: true},
		},
		{
			name: "TERM comum",
			env:  map[string]string{"TERM": "xterm-256color"},
			want: Profile{Unicode: true},
		},
		{
			name: "A11Y_HIGH_CONTRAST=1",
			env:  map[string]string{"A11Y_HIGH_CONTRAST": "1"},
			want: Profile{HighContrast: true, Unicode: true},
		},
		{
			name: "A11Y_HIGH_CONTRAST=true",
			env:  map[string]string{"A11Y_HIGH_CONTRAST": "true"},
			want: Profile{HighContrast: true, Unicode: true},
		},
		{
			name: "A11Y_HIGH_CONTRAST=TRUE sem caixa",
			env:  map[string]string{"A11Y_HIGH_CONTRAST": "TRUE"},
			want: Profile{HighContrast: true, Unicode: true},
		},
		{
			name: "A11Y_HIGH_CONTRAST=0",
			env:  map[string]string{"A11Y_HIGH_CONTRAST": "0"},
			want: Profile{Unicode: true},
		},
		{
			name: "A11Y_REDUCED_MOTION=1",
			env:  map[string]string{"A11Y_REDUCED_MOTION": "1"},
			want: Profile{ReducedMotion: true, Unicode: true},
		},
		{
			name: "A11Y_SCREEN_READER=true",
			env:  map[string]string{"A11Y_SCREEN_READER": "true"},
			want: Profile{ScreenReader: true, Unicode: true},
		},
		{
			name: "ASCII=1",
			env:  map[string]string{"ASCII": "1"},
			want: Profile{},
		},
		{
			name: "A11Y_ASCII=1",
			env:  map[string]string{"A11Y_ASCII": "1"},
			want: Profile{},
		},
		{
			name: "A11Y_ASCII=false não desliga Unicode",
			env:  map[string]string{"A11Y_ASCII": "false"},
			want: Profile{Unicode: true},
		},
		{
			name: "COLUMNS válido",
			env:  map[string]string{"COLUMNS": "120"},
			want: Profile{Unicode: true, Width: 120},
		},
		{
			name: "COLUMNS inválido",
			env:  map[string]string{"COLUMNS": "largo"},
			want: Profile{Unicode: true},
		},
		{
			name: "COLUMNS negativo",
			env:  map[string]string{"COLUMNS": "-1"},
			want: Profile{Unicode: true},
		},
		{
			name: "combinação completa",
			env: map[string]string{
				"TERM":                "dumb",
				"A11Y_HIGH_CONTRAST":  "1",
				"A11Y_REDUCED_MOTION": "true",
				"A11Y_SCREEN_READER":  "1",
				"COLUMNS":             "40",
			},
			want: Profile{
				NoColor:       true,
				Dumb:          true,
				HighContrast:  true,
				ReducedMotion: true,
				ScreenReader:  true,
				Width:         40,
			},
		},
		{
			name: "valores com espaços",
			env: map[string]string{
				"A11Y_HIGH_CONTRAST": " true ",
				"CLICOLOR":           " 0 ",
			},
			want: Profile{NoColor: true, HighContrast: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectFrom(envGetter(tt.env))
			if got != tt.want {
				t.Fatalf("DetectFrom(%v) = %#v, quer %#v", tt.env, got, tt.want)
			}
		})
	}
}

func TestDetectFromNilGetter(t *testing.T) {
	got := DetectFrom(nil)
	want := Profile{Unicode: true}
	if got != want {
		t.Fatalf("DetectFrom(nil) = %#v, quer %#v", got, want)
	}
}

func TestDetectUsesProcessEnvironment(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("CLICOLOR", "")
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("ASCII", "")
	t.Setenv("A11Y_ASCII", "")
	t.Setenv("A11Y_HIGH_CONTRAST", "")
	t.Setenv("A11Y_REDUCED_MOTION", "")
	t.Setenv("A11Y_SCREEN_READER", "")
	t.Setenv("COLUMNS", "90")

	got := Detect()
	want := Profile{Unicode: true, Width: 90}
	if got != want {
		t.Fatalf("Detect() = %#v, quer %#v", got, want)
	}
}

func TestDefaultIsStable(t *testing.T) {
	first := Default()
	second := Default()
	if first != second {
		t.Fatalf("Default() não é estável: %#v != %#v", first, second)
	}
}

func TestProfileString(t *testing.T) {
	tests := []struct {
		name    string
		profile Profile
		want    string
	}{
		{
			name:    "zero",
			profile: Profile{},
			want:    "color=on unicode=off width=unknown",
		},
		{
			name:    "sem cor",
			profile: Profile{NoColor: true, Width: -3},
			want:    "color=off unicode=off width=unknown",
		},
		{
			name:    "com largura",
			profile: Profile{Unicode: true, Width: 80, HighContrast: true},
			want:    "color=on unicode=on width=80 high-contrast",
		},
		{
			name: "todas as marcas",
			profile: Profile{
				Dumb:          true,
				HighContrast:  true,
				ReducedMotion: true,
				ScreenReader:  true,
				Width:         40,
			},
			want: "color=on unicode=off width=40 dumb high-contrast reduced-motion screen-reader",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.profile.String(); got != tt.want {
				t.Fatalf("String() = %q, quer %q", got, tt.want)
			}
			if got, want := tt.profile.Describe(), "a11y: "+tt.want; got != want {
				t.Fatalf("Describe() = %q, quer %q", got, want)
			}
		})
	}
}

func TestProfileSymbol(t *testing.T) {
	tests := []struct {
		name    string
		profile Profile
		want    string
	}{
		{name: "unicode", profile: Profile{Unicode: true}, want: "✓"},
		{name: "ascii", profile: Profile{Unicode: false}, want: "ok"},
		{name: "zero value", profile: Profile{}, want: "ok"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.profile.Symbol("✓", "ok"); got != tt.want {
				t.Fatalf("Symbol() = %q, quer %q", got, tt.want)
			}
		})
	}

	if got := (Profile{Unicode: true}).Symbol("", "x"); got != "" {
		t.Fatalf("Symbol vazio = %q, quer vazio", got)
	}
}
