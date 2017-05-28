#version 330
uniform mat4 MVP;

in vec3 POSITION;
in vec4 COLOR;
in float SIZE;

out vec4 vs_color;

void main()
{

  gl_PointSize = SIZE;
  gl_Position = vec4(POSITION, 1.0);

  vs_color = COLOR;
}