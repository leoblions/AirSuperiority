package main

// embed resource files in exe as byte arrays

import (
	_ "embed"
)

// DATA

//go:embed data/charmap_letters.cfg
var CharmapLetters []byte

// IMAGES

//go:embed images/clouds1.png
var ImageClouds1 []byte

//go:embed images/clouds2.png
var ImageClouds2 []byte

//go:embed images/sky.png
var ImageOcean []byte

//go:embed images/airplanes1.png
var Airplanes1 []byte

//go:embed images/airplanes2.png
var Airplanes2 []byte

//go:embed images/airplanes3.png
var Airplanes3 []byte

//go:embed images/explosions1.png
var Explosions1 []byte

//go:embed images/explosions2.png
var Explosions2 []byte

//go:embed images/explosions3.png
var Explosions3 []byte

//go:embed images/airplanePlayer.png
var AirplanePlayer []byte

//go:embed images/rocketR.png
var RocketR []byte

//go:embed images/rocketY.png
var RocketY []byte

//go:embed images/letters10x10.png
var RasterLetters []byte

//go:embed images/icons.png
var IconsPng []byte

//go:embed images/icons2.png
var Icons2Png []byte

//go:embed images/astitle.png
var TitlePng []byte

//go:embed images/optionstitle.png
var OptionsPng []byte

// SOUND

// c.addSFXByFilename("boom.wav")    //0
// 	c.addSFXByFilename("exp1.wav")    //1
// 	c.addSFXByFilename("exp2.wav")    //2
// 	c.addSFXByFilename("hiss.wav")    //3
// 	c.addSFXByFilename("launch1.wav") //4
// 	c.addSFXByFilename("launch2.wav") //5
// 	c.addSFXByFilename("tone1.wav")   //6

//go:embed sound/boom.wav
var BoomWav []byte

//go:embed sound/exp1.wav
var Exp1Wav []byte

//go:embed sound/exp2.wav
var Exp2Wav []byte

//go:embed sound/hiss.wav
var HissWav []byte

//go:embed sound/launch1.wav
var Launch1Wav []byte

//go:embed sound/launch2.wav
var Launch2Wav []byte

//go:embed sound/tone1.wav
var Tone1Wav []byte
