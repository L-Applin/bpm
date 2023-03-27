package java

import "bpm"

type JavaVersion interface {
	GetVersion() string
}

type javaVersion struct {
	version string
}

func (jv javaVersion) GetVersion() string {
	return jv.version
}

type JavaFlavor interface {
	GetFlavor() string
}

type javaFlavor struct {
	flavor string
}

func (jf javaFlavor) GetFlavor() string {
	return jf.flavor
}

var (
	Jdk11 JavaVersion = javaVersion{version: "jdk8"}
	Jdk17             = javaVersion{version: "jdk17:"}
	Jdk21             = javaVersion{version: "jdk21:"}
)

var (
	Maven  JavaFlavor = javaFlavor{flavor: "maven"}
	Gradle            = javaFlavor{flavor: "gradle"}
)

func Agent(version JavaVersion, flavor JavaFlavor) bpm.Agent {
	return bpm.Agent("bpm." + version.GetVersion() + "." + flavor.GetFlavor())
}
