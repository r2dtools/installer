package utils

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	SYSTEMD_SERVICE_FILE = "/etc/systemd/system/r2dtools.service"
)

type SystemD struct {
	Logger *log.Logger
	Sh     *SH
}

func (sd *SystemD) CreateService(binPath, userName, groupName string) error {
	serviceFile, err := os.OpenFile(SYSTEMD_SERVICE_FILE, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer serviceFile.Close()

	content := fmt.Sprintf(`[Unit]
	Description=R2DTools agent service
	After=network.target
	StartLimitIntervalSec=0
	
	[Service]
	Type=simple
	Restart=always
	RestartSec=1
	User=%s
	Group=%s
	ExecStart=%s
	
	[Install]
	WantedBy=multi-user.target`, userName, groupName, binPath+" serve")
	if _, err := io.WriteString(serviceFile, content); err != nil {
		return err
	}

	return nil
}

func (sd *SystemD) StartService(serviceName string) error {
	sd.Logger.Println("starting R2DTools agent service ...")
	sd.Sh.Exec("systemctl daemon-reload")

	if err := sd.Sh.Exec(fmt.Sprintf("systemctl start \"%s\"", serviceName)); err != nil {
		return fmt.Errorf("could not start '%s' service: %v", serviceName, err)
	}

	if err := sd.Sh.Exec(fmt.Sprintf("systemctl status \"%s\"", serviceName)); err != nil {
		return fmt.Errorf("could not start '%s' service: %v", serviceName, err)
	}

	if err := sd.Sh.Exec(fmt.Sprintf("systemctl enable \"%s\"", serviceName)); err != nil {
		return fmt.Errorf("could not start '%s' service: %v", serviceName, err)
	}
	sd.Logger.Println("R2DTools agent service successfully started")

	return nil
}

func (sd *SystemD) StopService(serviceName string) error {
	sd.Logger.Println("stopping R2DTools agent service ...")
	if err := sd.Sh.Exec(fmt.Sprintf("systemctl stop \"%s\"", serviceName)); err != nil {
		return fmt.Errorf("could not stop '%s' service: %v", serviceName, err)
	}

	return nil
}

func (sd *SystemD) RemoveService(serviceName string) error {
	sd.Logger.Println("disabling R2DTools agent systemd service ...")
	if err := sd.Sh.Exec(fmt.Sprintf("systemctl disable \"%s\"", serviceName)); err != nil {
		return fmt.Errorf("could not disable '%s' service: %v", serviceName, err)
	}
	os.Remove(SYSTEMD_SERVICE_FILE)
	sd.Sh.Exec("systemctl daemon-reload")

	return nil
}
