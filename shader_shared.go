package main

// ShaderShared defines shader code which is shared and imported by other programs.
const ShaderShared = `
	layout(std140, binding = 0) uniform Shared {
		mat4 View;
		mat4 Projection;
	};

	// Simulation cell states.
	//
	// These need to stay in sync with the Cell constants in palette.go
	const uint CellEmpty = 0;
	const uint CellWire  = 50;
	const uint CellTail  = 100;
	const uint CellHead  = 255;
	`
