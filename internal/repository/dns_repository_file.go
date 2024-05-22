package repository

import (
	"encoding/json"
	"os"

	"github.com/Kingfish219/PlaNet/network/dns"
)

type DnsRepositoryFile struct {
	filePath string
}

func NewDnsRepository(filePath string) *DnsRepositoryFile {
	return &DnsRepositoryFile{
		filePath: filePath,
	}
}

func (repo *DnsRepositoryFile) ReadDb() (*FileDb, error) {
	file, err := os.ReadFile(repo.filePath)
	if err != nil {
		return nil, err
	}

	if len(file) == 0 {
		return &FileDb{}, nil
	}

	fileDb := FileDb{}
	err = json.Unmarshal(file, &fileDb)
	if err != nil {
		return nil, err
	}

	return &fileDb, nil
}

func (repo *DnsRepositoryFile) WriteDb(fileDb *FileDb) error {
	var jsonData []byte
	jsonData, err := json.Marshal(fileDb)
	if err != nil {
		return err
	}

	err = os.WriteFile(repo.filePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (repo *DnsRepositoryFile) GetActiveDnsConfiguration() (dns.Dns, error) {
	file, err := repo.ReadDb()
	if err != nil {
		return dns.Dns{}, err
	}

	return file.ActiveDns, nil
}

func (repo *DnsRepositoryFile) ModifyActiveDnsConfiguration(dns.Dns) error {
	activeDns, err := repo.GetActiveDnsConfiguration()
	if err != nil {
		return err
	}

	fileDb, err := repo.ReadDb()
	if err != nil {
		return err
	}

	fileDb.ActiveDns = activeDns
	var jsonData []byte
	jsonData, err = json.Marshal(fileDb)
	if err != nil {
		return err
	}

	err = os.WriteFile(repo.filePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (repo *DnsRepositoryFile) GetDnsConfigurations() ([]dns.Dns, error) {
	file, err := os.ReadFile(repo.filePath)
	if err != nil {
		return nil, err
	}

	if len(file) == 0 {
		return []dns.Dns{}, nil
	}

	var dnsList []dns.Dns
	err = json.Unmarshal(file, &dnsList)
	if err != nil {
		return nil, err
	}

	return dnsList, nil
}

func (repo *DnsRepositoryFile) ModifyDnsConfigurations(dns dns.Dns) error {
	dnsList, err := repo.GetDnsConfigurations()
	if err != nil {
		return err
	}

	dnsList = append(dnsList, dns)
	fileDb, err := repo.ReadDb()
	if err != nil {
		return err
	}

	fileDb.DnsConfigurations = dnsList
	if fileDb.ActiveDns.Name == "" {
		fileDb.ActiveDns = dnsList[0]
	}

	var jsonData []byte
	jsonData, err = json.Marshal(fileDb)
	if err != nil {
		return err
	}

	err = os.WriteFile(repo.filePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
