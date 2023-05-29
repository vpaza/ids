/*
 * Copyright Daniel Hawton
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package config

type Config struct {
	Database ConfigDatabase `json:"database"`
	Facility ConfigFacility `json:"facility"`
	OAuth    ConfigOAuth    `json:"oauth"`
	Server   ConfigServer   `json:"server"`
	Session  ConfigSession  `json:"session"`
}

type ConfigFacility struct {
	Name string `json:"name"`
}

type ConfigDatabase struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	Username    string `json:"user"`
	Password    string `json:"password"`
	Database    string `json:"database"`
	AutoMigrate bool   `json:"auto_migrate"`
	CACert      string `json:"ca_cert"`
}

type ConfigOAuth struct {
	BaseURL           string `json:"base_url"`
	ClientID          string `json:"client_id"`
	ClientSecret      string `json:"client_secret"`
	MyBaseURL         string `json:"my_base_url"`
	EndpointAuthorize string `json:"endpoint_authorize"`
	EndpointToken     string `json:"endpoint_token"`
}

type ConfigServer struct {
	Port    string `json:"port"`
	Mode    string `json:"mode"`
	SSLCert string `json:"ssl_cert"`
	SSLKey  string `json:"ssl_key"`
}

type ConfigSession struct {
	Name     string `json:"name"`
	Secret   string `json:"secret"`
	Domain   string `json:"domain"`
	Path     string `json:"path"`
	MaxAge   string `json:"max_age"`
	Secure   bool   `json:"secure"`
	SameSite string `json:"same_site"`
}
