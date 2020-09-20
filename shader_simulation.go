package main

// SimulationShader defines shader sources for a simulation.
// This implements the Wireworld rules.
var SimulationShader = ShaderSource{
	Vertex: `
		#version 420

		layout(location = 0) in vec2 vertPos;
		layout(location = 1) in vec2 vertUV;
		out vec2 fragUV;

		void main() {
			gl_Position = vec4(vertPos, 0, 1);
			fragUV = vertUV;
		}
		`,
	Fragment: `
		#version 420

		$INCLUDE_SHARED$

		layout (binding = 0) uniform sampler2D input;

		in  vec2 fragUV;
		out vec4 output;

		// countHeadNeighbours checks texels surrounding fragUV and
		// counts those which have the cellHead state.
		uint countHeadNeighbours() {
			// The sampled red components are converted to uint and in
			// the process are truncated to 0 if their value is < 1.0.
			// 1.0 happens to be the value of the CellHead state we are
			// interested in. All other states are discarded.

			// top row
			uint r00 = uint(textureOffset(input, fragUV, ivec2(-1, 1)).r);
			uint r01 = uint(textureOffset(input, fragUV, ivec2( 0, 1)).r);
			uint r02 = uint(textureOffset(input, fragUV, ivec2( 1, 1)).r);

			// middle row
			uint r10 = uint(textureOffset(input, fragUV, ivec2(-1, 0)).r);
			uint r12 = uint(textureOffset(input, fragUV, ivec2( 1, 0)).r);

			// bottom row
			uint r20 = uint(textureOffset(input, fragUV, ivec2(-1,-1)).r);
			uint r21 = uint(textureOffset(input, fragUV, ivec2( 0,-1)).r);
			uint r22 = uint(textureOffset(input, fragUV, ivec2( 1,-1)).r);

			// Sum all the cell states. At this point we only have non-zero
			// values for CellHead neighbours. So the function returns the
			// total number of neighbouring CellHeads and nothing more.
			return r00 + r01 + r02 +
			       r10 +       r12 +
				   r20 + r21 + r22;
		}

		void main() {
			uint cell  = uint(texture2D(input, fragUV).r * 255);

			switch (cell) {
			case CellWire:
				uint heads = countHeadNeighbours();
				if (heads == 1 || heads == 2) {
					cell = CellHead;
				}
				break;
			case CellHead:
				cell = CellTail;
				break;
			case CellTail:
				cell = CellWire;
				break;
			}

			output = vec4(float(cell) / 255, 0, 0, 1);
		}
		`,
}
