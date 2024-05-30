package systray

import (
	"fmt"
	"os"

	"github.com/Kingfish219/PlaNet/internal/interfaces"
	"github.com/Kingfish219/PlaNet/internal/presets"
	"github.com/Kingfish219/PlaNet/network/dns"
	"github.com/getlantern/systray"
)

type SystrayUI struct {
	dnsRepository             interfaces.DnsRepository
	dnsConfigurations         []dns.Dns
	selectedDnsConfiguration  dns.Dns
	connectedDnsConfiguration dns.Dns
}

func New(dnsRepository interfaces.DnsRepository) *SystrayUI {
	return &SystrayUI{
		dnsRepository: dnsRepository,
	}
}

func (systrayUI *SystrayUI) Initialize() error {
	systray.Run(systrayUI.onReady, systrayUI.onExit)

	return nil
}

func (systrayUI *SystrayUI) onExit() {
	fmt.Println("Exiting")
}

func (systrayUI *SystrayUI) onReady() {
	systrayUI.setIcon(false)
	systrayUI.setToolTip("Not connected")

	systrayUI.addDnsConfigurations()

	systray.AddSeparator()
	menuExit := systray.AddMenuItem("Exit", "Exit the application")

	go func() {
		<-menuExit.ClickedCh
		systray.Quit()
	}()
}

func (systrayUI *SystrayUI) setIcon(status bool) {
	fileName := "idle"
	if status {
		fileName = "success"
	}

	filePath := "./assets/" + fileName + ".ico"
	ico, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Unable to read icon:", err)
	} else {
		systray.SetIcon(ico)
	}
}

func (systrayUI *SystrayUI) setToolTip(toolTip string) {
	systray.SetTooltip("PlaNet:\n" + toolTip)
}

func (console *SystrayUI) Consume(command string) error {
	return nil
}

func (systrayUI *SystrayUI) addDnsConfigurations() error {
	dnsConfigurations, err := systrayUI.dnsRepository.GetDnsConfigurations()
	if err != nil {
		return err
	}

	if len(dnsConfigurations) == 0 {
		presetDnsList := presets.GetDnsPresets()
		for _, pre := range presetDnsList {
			systrayUI.dnsRepository.ModifyDnsConfigurations(pre)
		}

		dnsConfigurations, err = systrayUI.dnsRepository.GetDnsConfigurations()
		if err != nil {
			return err
		}
	}

	systrayUI.dnsConfigurations = dnsConfigurations

	dnsConfigMenu := systray.AddMenuItem(fmt.Sprintf("DNS config: %v", systrayUI.dnsConfigurations[0].Name), "Selected DNS Configuration")
	for _, dnsConfig := range systrayUI.dnsConfigurations {
		dnsConfigSubMenu := dnsConfigMenu.AddSubMenuItem(dnsConfig.Name, dnsConfig.Name)
		localDns := dnsConfig

		go func(localDns dns.Dns) {
			for {
				<-dnsConfigSubMenu.ClickedCh
				if systrayUI.connectedDnsConfiguration.Name != localDns.Name {
					dnsService := dns.DnsService{}
					_, err := dnsService.ChangeDns(dns.ResetDns, systrayUI.connectedDnsConfiguration)
					if err != nil {
						fmt.Println(err)

						return
					}
				}

				systrayUI.setIcon(false)
				dnsConfigMenu.SetTitle(fmt.Sprintf("DNS config: %v", localDns.Name))
				systrayUI.selectedDnsConfiguration = localDns
			}

		}(localDns)
	}

	menuSet := systray.AddMenuItem("Set DNS", "Set DNS")
	menuReset := systray.AddMenuItem("Reset DNS", "Reset DNS")

	go func() {
		for {
			<-menuSet.ClickedCh
			dnsService := dns.DnsService{}
			_, err := dnsService.ChangeDns(dns.SetDns, systrayUI.selectedDnsConfiguration)
			if err != nil {
				fmt.Println(err)

				return
			}

			fmt.Println("Shecan set successfully.")

			systrayUI.setIcon(true)
			systrayUI.setToolTip("Connected to: Shecan")
		}

	}()

	go func() {
		for {
			<-menuReset.ClickedCh
			fmt.Println(systrayUI.selectedDnsConfiguration)

			dnsService := dns.DnsService{}
			_, err := dnsService.ChangeDns(dns.ResetDns, systrayUI.connectedDnsConfiguration)
			if err != nil {
				fmt.Println(err)

				return
			}

			fmt.Println("Shecan disconnected successfully.")

			systrayUI.setIcon(false)
			systrayUI.setToolTip("Not connected")
		}

	}()

	return nil
}
