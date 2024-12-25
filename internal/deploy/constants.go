package deploy

var (
	DefaultConfigName   = "wb-go-deploy.yaml"
	SshPort             = 22
	SshUser             = "root"
	SshPassword         = "wirenboard"
	DefaultAppDir       = "/mnt/data/"
	systemdTemplateFile = "templates/systemd.gotmpl"
)
