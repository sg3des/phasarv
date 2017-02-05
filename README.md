# PHASARV - WIP

This game project developed on the GO language. 

The game is a topdown dynamic battle.

![screenshot](screenshots/screen.png)


## About code and project structure

All code is in the **vendor** directory - as the used internal packages can not be used elsewhere.

Idea about project structure is, use one code for client and server.

Game engine have two separately loop: in first only render opengl, in second, everything else - physics, game logic, callbacks etc. As opengl should be in main process, second loop launched with goroutine.

For network interaction is used only **UDP** packages. For the packaging network data used **encoding/gob**.


## Engines

Opengl engine - [tbogdala/fizzle](https://github.com/tbogdala/fizzle)

2dPhys engine - chipmunk [TheZeroSlave/chipmunk](https://github.com/TheZeroSlave/chipmunk) forked from vova616/chipmunk


## TODO

- Models:
	- Airplanes
	- Trees
	- Environment objects

- Lightining:
	- Directional shadows
	- Blur shadows
	- [?] Ambient Occlusion

- Effects:
	- Explosions
	- Fog & clouds
	- [?] Water

- UI

- Map Editor

- Network