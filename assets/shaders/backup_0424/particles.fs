#version 330
//precision highp float;

uniform sampler2D MATERIAL_TEX_DIFFUSE;

in vec4 vs_color;
out vec4 frag_color;

void main()
{
	float alph = texture(MATERIAL_TEX_DIFFUSE, gl_PointCoord.st).r;
	if (alph < 0.10)  {
		discard;
	}

	frag_color = vs_color * texture(MATERIAL_TEX_DIFFUSE, gl_PointCoord.st);
}