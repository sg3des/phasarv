package assets

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sg3des/fizzgui"
	"github.com/tbogdala/fizzle"
	"github.com/tbogdala/fizzle/graphicsprovider"
	"github.com/tbogdala/fizzle/renderer/forward"
	"github.com/tbogdala/gombz"
)

var (
	Textures = make(map[string]*texture)
	Shaders  = make(map[string]*fizzle.RenderShader)
	Models   = make(map[string]*fizzle.Renderable)
	Fonts    = make(map[string]*font)

	textureMan = fizzle.NewTextureManager()
)

type texture struct {
	Diffuse  graphicsprovider.Texture
	Normals  graphicsprovider.Texture
	Specular graphicsprovider.Texture
}

type font struct {
	Font *fizzgui.Font
	Path string
	Name string
	Size int
}

func LoadAssets(texDir, shadersDir, modelsDir, fontsDir string) error {
	textures, err := filepath.Glob(filepath.Join(texDir, "*.png"))
	if err != nil {
		return err
	}
	fmt.Println("finded", len(textures), "textures")
	for _, texture := range textures {
		if err := LoadTexture(texture); err != nil {
			return errors.New("failed to load texture `" + texture + "` reason: " + err.Error())
		}
	}

	shaders, err := filepath.Glob(filepath.Join(shadersDir, "*.vs"))
	if err != nil {
		return err
	}
	fmt.Println("finded", len(shaders), "shaders")
	for _, shader := range shaders {
		if err := LoadShader(shader); err != nil {
			return errors.New("failed to load shader `" + shader + "` reason: " + err.Error())
		}
	}
	Shaders["basic"], err = forward.CreateBasicShader()
	if err != nil {
		log.Fatalln("error with basic shader:", err)
	}
	// Shaders["basic"] = Shaders[""]

	models, err := filepath.Glob(filepath.Join(modelsDir, "*.gombz"))
	if err != nil {
		return err
	}
	fmt.Println("finded", len(models), "models")

	for _, model := range models {
		modelname, mesh, err := LoadModel(model)
		if err != nil {
			return errors.New("failed to load model `" + model + "` reason: " + err.Error())
		}
		fmt.Println(modelname, mesh.FaceCount)
		Models[modelname] = mesh
	}

	fonts, err := filepath.Glob(filepath.Join(fontsDir, "*.ttf"))
	if err != nil {
		return err
	}
	if len(fonts) == 0 {
		return fmt.Errorf("fonts by path %s not found", fontsDir)
	}

	var size int
	var name string
	for _, fontpath := range fonts {
		name = strings.TrimSuffix(filepath.Base(fontpath), filepath.Ext(fontpath))
		ss := strings.Split(name, "-")
		if len(ss) != 2 {
			return errors.New("invalid font name: " + fontpath)
		}
		name = ss[0]
		size, err = strconv.Atoi(ss[1])
		if err != nil {
			return err
		}

		Fonts[name] = &font{Path: fontpath, Name: name, Size: size}
		// _, err = fizzgui.NewFont(name, font, size, fizzgui.FontGlyphs)
		// if err != nil {
		// 	return err
		// }
	}

	return nil
}

var (
	texNorm = "_N"
	texSpec = "_S"
)

func LoadTexture(filename string) (err error) {
	name := filepath.Base(filename)
	name = strings.TrimSuffix(name, ".png")
	name = strings.TrimSuffix(name, texNorm)
	name = strings.TrimSuffix(name, texSpec)

	var t *texture
	var ok bool

	if t, ok = Textures[name]; !ok {
		t = new(texture)
		Textures[name] = t
	}

	if strings.HasSuffix(filename, texNorm+".png") {
		t.Normals, err = textureMan.LoadTexture(name+texNorm, filename)
	} else if strings.HasSuffix(filename, texSpec+".png") {
		t.Specular, err = textureMan.LoadTexture(name+texSpec, filename)
	} else {
		t.Diffuse, err = textureMan.LoadTexture(name, filename)
	}

	return
}

func LoadShader(filename string) error {
	dirname := filepath.Dir(filename)
	basename := filepath.Base(filename)
	shadername := strings.TrimSuffix(basename, filepath.Ext(basename))

	var err error
	Shaders[shadername], err = fizzle.LoadShaderProgramFromFiles(filepath.Join(dirname, shadername), nil)
	if err != nil {
		return err
	}

	return nil
}

func LoadModel(filename string) (string, *fizzle.Renderable, error) {
	gombzmesh, err := gombz.DecodeFile(filename)
	if err != nil {
		return "", nil, err
	}

	basename := filepath.Base(filename)
	modelname := strings.TrimSuffix(basename, filepath.Ext(basename))

	return modelname, fizzle.CreateFromGombz(gombzmesh), nil
}

func GetTexture(name string) *texture {
	tex, ok := Textures[name]
	if !ok {
		log.Fatalf("ERROR: texture `%s` not found", name)
	}
	return tex
}

func GetShader(name string) *fizzle.RenderShader {
	shader, ok := Shaders[name]
	if !ok {
		log.Fatalf("ERROR: shader `%s` not found", name)
	}
	return shader
}

func GetModel(name string) *fizzle.Renderable {
	mesh, ok := Models[name]
	if !ok {
		panic("ERROR: model not found! " + name)
		// log.Fatalf("ERROR: model `%s` not found!", name)
	}

	return mesh.Clone()
}

func GetFont(name string) *font {
	font, ok := Fonts[name]
	if !ok {
		log.Panicln("ERROR: font not found!", name)
	}
	return font
}

func GetFonts() (fonts []*font) {
	for _, font := range Fonts {
		fonts = append(fonts, font)
	}
	return
}
