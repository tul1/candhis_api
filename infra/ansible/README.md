# Ansible Candhis API Full Stack Deployer

This Ansible project is designed to automate the setup of the production environment for the candhis_api project. The project includes two scrapers (campaigns_scraper and sessionid_scraper) that are orchestrated by two systemd timers, along with an API. The playbooks in this project will configure the host, install the necessary infrastructure, and deploy the binaries required for the project.

## Requirements

To run these Ansible playbooks, ensure you have the following:

1. **Ansible**: Installed on your local machine or control node.
   - Installation guide: [Ansible Installation Guide](https://docs.ansible.com/ansible/latest/installation_guide/index.html)
2. **SSH Access**: Ensure you have SSH access to the target host.
   - SSH key for the `root` user is required for initial setup.
   - SSH key for the `astraydev` user after creation.
3. **Python Libraries**: Install the required Python libraries.
   ```bash
   pip install passlib
   ```

## Inventory Configuration

The `inventory` file should contain the target host information. Example:

```ini
[servers]
candhis_server ansible_host=95.179.209.34 ansible_user=root ansible_ssh_private_key_file=~/.ssh/id_rsa
```

## Playbooks

### 1. **Setup Host**

This playbook sets up the initial host environment, including user creation, Docker installation, automatic updates, and firewall configuration.

#### Steps Included:
- Create a new user `astraydev` with sudo privileges (if not already created).
- Install Docker and Docker Compose.
- Configure automatic updates.
- Install and enable UFW (Uncomplicated Firewall).

#### How to Run:

```bash
cd infra/ansible
ansible-playbook playbooks/setup_host.yml --extra-vars "astraydev_password=your_secure_password ssh_key_path=/home/ptula/personal/id_rsa_github.pub"
```

Replace `your_secure_password` for your password and `/path/to/id_rsa.pub` for the path to your ssh public key.

### 2. **Setup Application Infrastructure**

This playbook sets up the infrastructure services required by the application, including running Docker containers for PostgreSQL, Elasticsearch, Fluentd, and Kibana. Additionally, it inserts initial data into the database and configures Kibana access.

#### Steps Included:
- Copy the `infra/` directory and `docker-compose.yml` file to the `/home/astraydev/candhis_api` directory on the host.
- Run the `run-infra` Makefile target to bring up all infrastructure services.
- Insert a new session ID and date into the `candhis_session` table in PostgreSQL if the table is empty.
- Open port 5601 in UFW and iptables to allow Kibana access.

#### How to Run:

```bash
cd infra/ansible
ansible-playbook playbooks/setup_app_infra.yml
```

### 3. **Install Application**

This playbook manages the installation of the application binaries and configurations, and sets up systemd services to manage the application.

#### Steps Included:
- Ensure the `/home/astraydev/candhis_api/bin/` directory exists on the host.
- Copy the `campaigns_scraper` and `sessionid_scraper` binaries to the `/home/astraydev/candhis_api/bin/` directory.
- Ensure the `/home/astraydev/candhis_api/config/` directory exists on the host.
- Copy the app configuration files from the `config/` directory to `/home/astraydev/candhis_api/config/`.
- Set up and manage systemd services and timers for `campaigns_scraper` and `sessionid_scraper`.

#### How to Run:

```bash
cd infra/ansible
ansible-playbook playbooks/install_app.yml
```

