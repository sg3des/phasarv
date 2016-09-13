#version 330
precision highp float;

uniform sampler2D MATERIAL_TEX_0;
uniform vec4 MATERIAL_DIFFUSE;

in vec2 vs_tex0_uv;
out vec4 frag_color;

void main (void) {
	frag_color = texture(MATERIAL_TEX_0, vs_tex0_uv);
	if (frag_color.a < 0.1) {
		discard;
	}

	// frag_color.a = MATERIAL_DIFFUSE.a;
	// if (MATERIAL_DIFFUSE.a < 0.9) {
	// 	discard;
	// }
}
