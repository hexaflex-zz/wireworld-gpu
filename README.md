## Wireworld-gpu

Wireworld implements the data and rules for the [Wireworld cellular automata](https://en.wikipedia.org/wiki/Wireworld).

This particular version is an experiment whereby the simulation is run
entirely on the GPU using multiple render passes whereby a fragment shader
alternates between two framebuffers for input and output. Meaning the output
from one render pass becomes the input for the next. The framebuffers contain
the simulation state. The fragment shader reads from the input buffer,
applies the Wireworld rules and then writes the new state to the output buffer.

Program state can be retrieved as an image and saved to disk.

Programs to be run can be provided by passing the path to an image file as a
command line parameter.

The program uses OpenGL `v4.2-core` with and `GLFW 3.3`.
It has been tested on a `GeForce GTX 750 Ti` with driver version `NVIDIA 436.02`.
On this system, the simulation runs at a speed of up to ~100KHz.


## Usage

    $ wireworld-gpu mysim.png

Use the `-help` flag for an overview of supported options.

The input image is meant to be drawn using a recognized color palette.
The fragment shader uses this palette to determine what kind of cell a
specific fragment represents.

The default palette is as follows:

 Cell State    | RGB Color
 --------------|------------
 Empty         | #000000
 Wire          | #015B96
 Electron head | #ffffff
 Electron tail | #99ff00

---

The `testdata/palette.gpl` file contains a GIMP Palette with the default
colors recognized by this program, along with two extra colors you can use
to draw annotations.

The color palette can be changed by providing custom RGB values through
the respective `-pal-???` flags in the command line. These should match
the colors used in the input image.

Pixels with unrecognized colors in the input image are ignored and treated
as an Empty cell. This allows you to add drawings or text annotations to
the image, without it affecting the simulation.

Refer to the `testdata` directory for examples of images with Wireworld
simulations.


## Keyboard shortcuts

  Key               | Description
 -------------------|------------------------------------
  Escape            | Close the program.
  Q                 | Start/Stop the simulation.
  E                 | Perform a single simulation step.
  W                 | Increase the simulation speed by 10x.
  S                 | Decrease the simulation speed by 10x.
  F1                | Saves the current simulation state in `<timestamp>.<inputfile>.png`
  F2                | Loads latest simulation state from `<timestamp>.<inputfile>.png` where it picks the highest timestamp if more than one such file exists. If no such file is available, this does the same as F5.
  F5                | Reset the simulation (reloads the original input image).
  Space + Mousemove | Pan the camera left/right/up/down. 
  Mouse Scroll      | Zoom in/out. 
  C                 | Center the simulation in the window.

---

## License

Unless otherwise stated, this project and its contents are provided under a
3-Clause BSD license. Refer to the LICENSE file for its contents.