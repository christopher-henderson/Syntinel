package PAO

const (
	DockerCommand      = "docker"
	DockerBuild        = "build"
	DockerBuildForceRM = "--force-rm" // Force deletion of the temporary image.
	DockerBuildTag     = "-t"         // Give a tag to the built image.
	DockerRun          = "run"
	DockerRunRM        = "--rm"   // Delete container upon completion of run.
	DockerRunName      = "--name" // Name the container.
	DockerStop         = "stop"   // Stop the named container. Used if the container is killed prematurely.
	DockerRM           = "rm"     // Destroy the named container. Used if the container is killed prematurely.
	DockerScriptName   = "script.sh"
	DockerFile         = "Dockerfile"
	ContainerPrefix    = "executor_" // Namespace containers created by the executor.
)
