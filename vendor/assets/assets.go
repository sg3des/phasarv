package assets

import (
	"errors"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tbogdala/fizzle"
	"github.com/tbogdala/fizzle/graphicsprovider"
	"github.com/tbogdala/gombz"
)

var (
	Textures = make(map[string]texture)
	Shaders  = make(map[string]*fizzle.RenderShader)
	Models   = make(map[string]*fizzle.Renderable)

	textureMan = fizzle.NewTextureManager()
)

type texture struct {
	Diffuse  graphicsprovider.Texture
	Normals  graphicsprovider.Texture
	Specular graphicsprovider.Texture
}

func LoadAssets(texpath, shaderspath, modelspath string) error {
	textures, err := filepath.Glob(filepath.Join(texpath, "*_D.png"))
	if err != nil {
		return err
	}
	for _, texture := range textures {
		if err := LoadTexture(texture); err != nil {
			return errors.New("failed to load texture `" + texture + "` reason: " + err.Error())
		}
	}

	shaders, err := filepath.Glob(filepath.Join(shaderspath, "*.vs"))
	if err != nil {
		return err
	}
	for _, shader := range shaders {
		if err := LoadShader(shader); err != nil {
			return errors.New("failed to load shader `" + shader + "` reason: " + err.Error())
		}
	}

	models, err := filepath.Glob(filepath.Join(modelspath, "*.gombz"))
	if err != nil {
		return err
	}
	log.Println("finded", len(models), " models")

	for _, model := range models {
		modelname, mesh, err := LoadModel(model)
		if err != nil {
			return errors.New("failed to load model `" + model + "` reason: " + err.Error())
		}
		Models[modelname] = mesh
	}
	log.Println(Models)

	return nil
}

var (
	// textureNameCrop = regexp.MustCompile("_D\\.png")
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

	var t texture
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
	// return nil
}

func GetModel(modelname string) *fizzle.Renderable {
	mesh, ok := Models[modelname]
	if !ok {
		log.Fatal("model not found:", modelname, " all models", Models)
	}

	return mesh.Clone()
}

func GetImage(filename string) image.Image {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln("failed load image")
		return nil
	}

	image, _, err := image.Decode(f)
	if err != nil {
		log.Fatalln("failed decode image", err)
		return nil
	}

	return image
}
