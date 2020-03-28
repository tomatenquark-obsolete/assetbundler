package config

var _cmdWhitelist = []string{
	"ambient",
	"autograss",
	"base_1",
	"base_10",
	"base_2",
	"base_3",
	"base_4",
	"base_5",
	"base_6",
	"base_7",
	"base_8",
	"base_9",
	"blurlms",
	"blurskylight",
	"causticmillis",
	"causticscale",
	"cloudalpha",
	"cloudbox",
	"cloudboxcolour",
	"cloudcolour",
	"cloudfade",
	"cloudheight",
	"cloudlayer",
	"cloudscale",
	"cloudscrollx",
	"cloudscrolly",
	"exec",
	"fog",
	"fogcolour",
	"fogdomecap",
	"fogdomeclip",
	"fogdomeclouds",
	"fogdomecolour",
	"fogdomeheight",
	"fogdomemax",
	"fogdomemin",
	"grassalpha",
	"lightprecision",
	"lmshadows",
	"loadsky",
	"mapmodel",
	"mapmodelreset",
	"mapmsg",
	"mapsound",
	"maptitle",
	"minimapclip",
	"minimapcolour",
	"minimapheight",
	"mmodel",
	"mmodelfile",
	"setshader",
	"setshaderparam",
	"shadowmapambient",
	"shadowmapangle",
	"skybox",
	"skyboxcolour",
	"skylight",
	"skytexture",
	"skytexturelight",
	"spinclouds",
	"spinsky",
	"sunlight",
	"sunlightpitch",
	"sunlightscale",
	"sunlightyaw",
	"texalpha",
	"texcolor",
	"texffenv",
	"texlayer",
	"texoffset",
	"texrotate ",
	"texscale",
	"texscroll",
	"texture",
	"texturereset",
	"water2colour",
	"water2fog",
	"watercolour",
	"waterfallcolour",
	"waterfog",
	"waterspec",
	"yawsky",
}

// after init(), cmdWhitelist contains all ICOMMANDs allowd in a map config.
var cmdWhitelist = map[string]struct{}{}

func init() {
	for _, cmd := range _cmdWhitelist {
		cmdWhitelist[cmd] = struct{}{}
	}
}

// Contains provides a method to check wether an ICOMMAND is valid in map configs.
func ValidInMapCfg(cmd string) bool {
	_, ok := cmdWhitelist[cmd]
	return ok
}

// Maps the relation between ICOMMAND and the position of the file path in its parameter list
var filePathPosByCmd = map[string]int{
	"cloudlayer": 1,
	"exec":       1,
	"mapsound":   1,
	"mmodelfile": 1,
	"skybox":     1,
	"texture":    2,
}

func HasFilePathParam(cmd string) (bool, int) {
	pos, ok := filePathPosByCmd[cmd]
	return ok, pos
}
