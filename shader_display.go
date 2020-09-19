package main

// DisplayShader defines shader sources for the simulation display.
var DisplayShader = ShaderSource{
	Vertex: `
		#version 420

		$INCLUDE_SHARED$

		uniform mat4 Model;

		in  vec2 vertPos;
		in  vec2 vertUV;
		out vec2 fragUV;

		void main() {
			gl_Position = Projection * View * Model * vec4(vertPos, 0, 1);
			fragUV = vertUV;
		}
		`,
	Fragment: `
		#version 420

		$INCLUDE_SHARED$

		layout (binding = 0) uniform sampler2D input;

		uniform vec4 PalEmpty;
		uniform vec4 PalWire;
		uniform vec4 PalHead;
		uniform vec4 PalTail;

		in  vec2 fragUV;
		out vec4 output;

		void main() {
			uint cell = uint(texture2D(input, fragUV).r * 255);

			switch (cell) {
			case CellWire:
				output = PalWire;
				break;
			case CellHead:
				output = PalHead;
				break;
			case CellTail:
				output = PalTail;
				break;
			default:
				output = PalEmpty;
				break;
			}
		}
		`,
}
