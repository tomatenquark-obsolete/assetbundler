package config

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

type expectedLoadResult struct {
	files []FileToLoad
	err   error
}

func TestReadResources(t *testing.T) {
	tests := []struct {
		cfgFileName string
		expectedLoadResult
	}{
		{
			cfgFileName: "collide.cfg",
			expectedLoadResult: expectedLoadResult{
				files: []FileToLoad{
					{Command: "texture", Path: "textures/sky.png"},
					{Command: "texture", Path: "textures/default.png"},
					{Command: "texture", Path: "textures/default.png"},
					{Command: "texture", Path: "textures/nieb/sand01.jpg"},
					{Command: "exec", Path: "packages/textures/yves_allaire/ex512/package.cfg"},
					{Command: "texture", Path: "philipk/pk01_panel01a_d.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel01_local.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel01_s.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel01a_add.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel01b_d.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel01_local.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel01_s.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel01b_add.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel02_d.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel02_local.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel02_s.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel02a_add.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel03a_d.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel03_local.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel03_s.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel03a_add.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel03b_d.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel03_local.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel03_s.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel03b_add.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel_small01_d.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel_small01_local.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel_small01_s.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel_small01_add.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel_small02_d.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel_small02_local.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel_small02_s.jpg"},
					{Command: "texture", Path: "philipk/pk01_panel_small02_add.jpg"},
					{Command: "texture", Path: "caustics/caust00.png"},
					{Command: "texture", Path: "caustics/caust00.png"},
					{Command: "texture", Path: "fohlen/black.jpg"},
					{Command: "texture", Path: "fohlen/white.jpg"},
					{Command: "texture", Path: "fohlen/black.jpg"},
					{Command: "texture", Path: "fohlen/white.jpg"},
					{Command: "mmodelfile", Path: "models/dcp/blade_y/collide/obj.cfg"},
					{Command: "mapsound", Path: "ambience/doomish/rumble1.ogg"},
					{Command: "mapsound", Path: "ambience/hum.ogg"},
					{Command: "mapsound", Path: "jdagenet/hum1.wav"},
					{Command: "mapsound", Path: "owlish/BLEEPINGCOMPUTER2.wav"},
					{Command: "mapsound", Path: "jdagenet/hum1.wav"},
					{Command: "mapsound", Path: "jdagenet/hum1.wav"},
				},
				err: nil,
			},
		},
		{
			cfgFileName: "corrupted.cfg",
			expectedLoadResult: expectedLoadResult{
				files: []FileToLoad{},
				err:   tooFewArgs("texture", 8),
			},
		},
	}

	for _, test := range tests {
		t.Run("get files to load from "+test.cfgFileName, func(t *testing.T) {
			cfg, err := os.Open(path.Join(findModuleRoot(), "internal", "config", "testdata", test.cfgFileName))
			if err != nil {
				t.Errorf("failed to open file '%s'", test.cfgFileName)
			}
			defer cfg.Close()

			files, err := FilesToLoadFromCfg(cfg)
			if err != test.err {
				if (err == nil && test.err != nil) ||
					(err != nil && test.err == nil) ||
					(err.Error() != test.err.Error()) {
					t.Errorf("expected err to be '%s', but got '%s'", test.err, err)
				}
			}

			if len(files) != len(test.files) {
				t.Errorf("expected %d file paths (%#v), but got %d (%#v)",
					len(test.files), test.files,
					len(files), files)
			}

			for i, file := range files {
				if file != test.files[i] {
					t.Errorf("expected file '%v' at position %d, but got '%v'", test.files[i], i, file)
				}
			}
		})
	}
}

type expectedFilterResult struct {
	output string
	err    error
}

func TestFilter(t *testing.T) {
	tests := []struct {
		cfgFileName string
		expectedFilterResult
	}{
		{
			cfgFileName: "corrupted.cfg",
			expectedFilterResult: expectedFilterResult{
				output: "// corrupted map by maleficent author\n" +
					"\n" +
					"skybox \"this/is/actually/valid.png\"\n" +
					"texture bla \"this/is/actually/valid.png\"\n" +
					"maptitle \"this/is/actually/valid.png\"\n" +
					"\n" +
					"texture \"this/is/actually/valid.png\"\n" +
					"\n" +
					"",
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run("filter config "+test.cfgFileName, func(t *testing.T) {
			cfg, err := ioutil.ReadFile(path.Join(findModuleRoot(), "internal", "config", "testdata", test.cfgFileName))
			if err != nil {
				t.Errorf("failed to read file '%s'", test.cfgFileName)
			}
			output, err := Filter(string(cfg))
			if err != test.err {
				t.Errorf("expected err to be '%s', but got '%s'", test.err, err)
			}

			if output != test.output {
				t.Errorf("expected output to be '%s', but got '%s'", test.output, output)
			}
		})
	}
}

// mostly from https://golang.org/src/cmd/go/internal/modload/init.go
func findModuleRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir = filepath.Clean(dir)

	// Look for enclosing go.mod.
	for {
		if fi, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !fi.IsDir() {
			return dir
		}
		d := filepath.Dir(dir)
		if d == dir {
			break
		}
		dir = d
	}

	panic("can't find module root")
}
