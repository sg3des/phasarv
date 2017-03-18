package assets

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/tbogdala/fizzle"
	"github.com/tbogdala/fizzle/graphicsprovider"
	"github.com/tbogdala/fizzle/renderer/forward"
	"github.com/tbogdala/gombz"
)

var (
	Textures = make(map[string]*texture)
	Shaders  = make(map[string]*fizzle.RenderShader)
	Models   = make(map[string]*fizzle.Renderable)

	textureMan = fizzle.NewTextureManager()
)

type texture struct {
	Path     string
	Diffuse  graphicsprovider.Texture
	Normals  graphicsprovider.Texture
	Specular graphicsprovider.Texture
}

func LoadAssets(texpath, shaderspath, modelspath string) error {
	textures, err := filepath.Glob(filepath.Join(texpath, "*_D.png"))
	if err != nil {
		return err
	}
	fmt.Println("finded", len(textures), "textures")
	for _, texture := range textures {
		if err := LoadTexture(texture); err != nil {
			return errors.New("failed to load texture `" + texture + "` reason: " + err.Error())
		}
	}

	shaders, err := filepath.Glob(filepath.Join(shaderspath, "*.vs"))
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
		log.Fatalln(err)
	}
	// Shaders["basic"] = Shaders[""]

	models, err := filepath.Glob(filepath.Join(modelspath, "*.gombz"))
	if err != nil {
		return err
	}
	fmt.Println("finded", len(models), "models")

	for _, model := range models {
		modelname, mesh, err := LoadModel(model)
		if err != nil {
			return errors.New("failed to load model `" + model + "` reason: " + err.Error())
		}
		Models[modelname] = mesh
	}

	return nil
}

var (
	textureTypeD = "_D.png"
	textureTypeN = "_N.png"
	textureTypeS = "_S.png"
)

func getTextureName(filename string) string {
	return strings.TrimSuffix(filepath.Base(filename), textureTypeD)
	// return textureNameCrop.ReplaceAllString(filename, "")
}

func LoadTexture(filename string) error {
	texturename := getTextureName(filename)

	t := &texture{Path: filename}
	var err error
	t.Diffuse, err = textureMan.LoadTexture(texturename, filename)
	if err != nil {
		return err
	}

	t.Specular, _ = textureMan.LoadTexture(texturename, strings.Replace(filename, textureTypeD, textureTypeS, 1))
	t.Normals, _ = textureMan.LoadTexture(texturename, strings.Replace(filename, textureTypeD, textureTypeN, 1))

	Textures[texturename] = t
	return nil
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
