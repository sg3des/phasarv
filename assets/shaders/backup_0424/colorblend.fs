#version 330
precision highp float;

uniform sampler2D MATERIAL_TEX_DIFFUSE;
uniform vec4 MATERIAL_DIFFUSE;

in vec2 vs_tex0_uv;
out vec4 frag_color;

void main (void) {
	// Modify the color's alpha based on the alpha map used in
	// in sample2D channel 0

	vec4 tex0 = texture(MATERIAL_TEX_DIFFUSE, vs_tex0_uv);// * diff;
	if (tex0.a > MATERIAL_DIFFUSE.a) {
		tex0.a = MATERIAL_DIFFUSE.a;
	}

	// if (tex0.r > MATERIAL_DIFFUSE.r) {
		tex0.r = MATERIAL_DIFFUSE.r;
	// }

		// if (tex0.g > MATERIAL_DIFFUSE.g) {
		tex0.g = MATERIAL_DIFFUSE.g;
	// }

		// if (tex0.b > MATERIAL_DIFFUSE.b) {
		tex0.b = MATERIAL_DIFFUSE.b;
	// }
	


	frag_color = tex0;
	//mix(diff.rgb, tex0.rgb, tex0.a);
	// frag_color = vec4((tex0.r + diff.r)/2, (tex0.g + diff.g)/2, (tex0.b + diff.b)/2, (tex0.a + diff.a)/2 );
	

	// alpha = texture(alpha, MATERIAL_DIFFUSE.a).r;
	// if (alph < 0.10) {
	//     discard;
	// }
	// frag_color = vec4(diff.rgb, tex0.a);
}
