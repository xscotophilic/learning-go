# Define variables
VAULT_ADDR=https://127.0.0.1:8200

.PHONY: all install configure certs service start init vaultkv

all: install configure certs service start init vaultkv

install:
	sudo apt update && sudo apt install -y unzip jq curl
	curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
	echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(shell lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
	sudo apt update && sudo apt install -y vault
	vault --version

configure:
	sudo mkdir -p /etc/vault
	sudo chown -R vault:vault /etc/vault
	echo 'export VAULT_ADDR="$(VAULT_ADDR)"' >> ~/.bashrc
	source ~/.bashrc

certs:
	sudo mkdir -p /etc/vault/certs
	cd /etc/vault/certs && sudo openssl req -x509 -newkey rsa:4096 -nodes -keyout vault.key -out vault.crt -days 365 -subj "/CN=localhost" -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
	sudo cp /etc/vault/certs/vault.crt /usr/local/share/ca-certificates/vault.crt
	sudo update-ca-certificates

service:
	sudo printf 'storage "file" {\n  path = "/var/lib/vault"\n}\n\nlistener "tcp" {\n  address = "0.0.0.0:8200"\n  tls_cert_file = "/etc/vault/certs/vault.crt"\n  tls_key_file = "/etc/vault/certs/vault.key"\n}\n\ndisable_mlock = true\n\nui = true\n' | sudo tee /etc/vault/config.hcl > /dev/null
	sudo printf "[Unit]\nDescription=HashiCorp Vault\nRequires=network-online.target\nAfter=network-online.target\n\n[Service]\nExecStart=/usr/bin/vault server -config=/etc/vault/config.hcl\nExecReload=/bin/kill --signal HUP $$MAINPID\nRestart=on-failure\nUser=root\nGroup=root\nPermissionsStartOnly=true\nLimitMEMLOCK=infinity\nCapabilityBoundingSet=CAP_IPC_LOCK\nAmbientCapabilities=CAP_IPC_LOCK\nKillMode=process\nKillSignal=SIGINT\nTimeoutStopSec=5\nRestart=always\nRestartSec=2\n\n[Install]\nWantedBy=multi-user.target\n" | sudo tee /etc/systemd/system/vault.service > /dev/null
	sudo systemctl daemon-reload
	sudo systemctl enable vault

start:
	sudo systemctl start vault

init:
	bash -c ' \
		INIT_OUTPUT=$$(vault operator init -format=json); \
		echo "$$INIT_OUTPUT" | jq -r "{unseal_keys: .unseal_keys_b64, root_token: .root_token}" > ~/vault_keys.json; \
		KEY1=$$(jq -r ".unseal_keys[0]" ~/vault_keys.json); \
		KEY2=$$(jq -r ".unseal_keys[1]" ~/vault_keys.json); \
		KEY3=$$(jq -r ".unseal_keys[2]" ~/vault_keys.json); \
		VAULT_ROOT_TOKEN=$$(jq -r ".root_token" ~/vault_keys.json); \
		vault operator unseal "$$KEY1"; \
		vault operator unseal "$$KEY2"; \
		vault operator unseal "$$KEY3"; \
		vault login "$$VAULT_ROOT_TOKEN"; \
	'

vaultkv: 
	bash -c ' vault secrets enable -path=userdata kv;'
