# Local Development Setup with a Global Traefik Proxy

To develop multiple projects (like Kasseapparat, Tidsapparat, Billedapparat, Funkapparat) simultaneously on your local machine without running into port conflicts (Port 80/443), we use a central, global reverse proxy. 

Each project keeps its own, isolated `docker-compose.yaml` to ensure a "Zero-Config" experience for new developers. The connection to the global proxy is established exclusively on your local machine via a git-ignored `docker-compose.override.yaml`.

## Architecture

* **Master-Traefik:** Runs continuously in the background on your machine, binds to ports 80 & 443, and manages all SSL certificates.
* **Global Network:** A Docker network named `traefik-global` that allows the proxy to route traffic to the containers of individual projects.
* **Project Containers:** Are attached to the Master-Traefik via Docker labels. The project's internal Traefik is disabled locally.

---

## 1. One-Time Setup (Master-Traefik)

This step only needs to be executed once per developer machine, regardless of how many projects you are working on.

### 1.1 Create the Global Network
Run the following command in your terminal:
````bash
docker network create traefik-global
````

### 1.2 Create the Master-Traefik Directory
Create a directory outside of your project repositories (e.g., `~/Projects/traefik-proxy`) and set up the following structure:

````text
traefik-proxy/
├── docker-compose.yaml
├── traefik-dynamic.yaml
└── certs/
    ├── tidsapparat.test.crt
    ├── tidsapparat.test.key
    └── ... (additional certificates for other projects)
````

### 1.3 Master-Traefik Configuration (`docker-compose.yaml`)
Create the `docker-compose.yaml` with the following content. Note the use of `restart: unless-stopped`, which ensures the proxy starts automatically with Docker.

````yaml
services:
  traefik:
    image: traefik:v3.7
    container_name: traefik_global
    restart: unless-stopped
    command:
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--providers.docker.network=traefik-global" 
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
      - "--entrypoints.web.http.redirections.entryPoint.to=websecure"
      - "--entrypoints.web.http.redirections.entryPoint.scheme=https"
      - "--providers.file.filename=/etc/traefik/dynamic.yaml"
    ports:
      - "80:80"
      - "443:443"
      - "8080:8080" # Traefik Dashboard
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./certs:/certs:ro 
      - ./traefik-dynamic.yaml:/etc/traefik/dynamic.yaml:ro
    extra_hosts:
      - "host.docker.internal:host-gateway"
    networks:
      - traefik-global

networks:
  traefik-global:
    external: true
````

### 1.4 Configure Certificates & Hosts (`traefik-dynamic.yaml`)
The dynamic configuration maps the paths to the certificates inside the container and routes traffic to host processes (like Vite and the backend):

````yaml
http:
  routers:
    kasseapparat-api:
      rule: "Host(`kasseapparat.test`) && PathPrefix(`/api`)"
      entryPoints:
        - websecure
      service: kasseapparat-backend-service
      tls: {}

    kasseapparat-vite:
      rule: "Host(`kasseapparat.test`)"
      entryPoints:
        - websecure
      service: kasseapparat-vite-service
      tls: {} 

  services:
    kasseapparat-backend-service:
      loadBalancer:
        servers:
          - url: "http://host.docker.internal:3001"

    kasseapparat-vite-service:
      loadBalancer:
        servers:
          - url: "http://host.docker.internal:3000"

tls:
  certificates:
    - certFile: /certs/kasseapparat.test.crt
      keyFile: /certs/kasseapparat.test.key
````

### 1.5 Start Master-Traefik
Spin up the proxy in the background:
````bash
cd ~/Projects/traefik-proxy
docker compose up -d
````

---

## 2. Connect a Project (Example: Kasseapparat)

To keep the project repository clean, we manipulate the local setup using an override file.

Create a file named `docker-compose.override.yaml` in the root directory of your project:

````yaml
services:
  # 1. Disables the project's internal Traefik
  traefik:
    profiles:
      - disabled

  # 2. Connects web-facing services to the global proxy
  mailhog:
    networks:
      - default
      - traefik-global
    labels:
      - "traefik.docker.network=traefik-global"

  openobserve:
    networks:
      - default
      - traefik-global
    labels:
      - "traefik.docker.network=traefik-global"

  redisinsight:
    networks:
      - default
      - traefik-global
    labels:
      - "traefik.docker.network=traefik-global"

# 3. Declares the global network for this project
networks:
  default:
  traefik-global:
    external: true
````

## 3. Workflow

From now on, you can start the project normally using the following command:

````bash
docker compose up -d
````

Docker will automatically read the `docker-compose.yaml`, silently apply your `docker-compose.override.yaml`, skip the internal Traefik container, and route all web-facing containers to your Master-Traefik. 
You can do this for any number of projects in parallel, and they will all be accessible simultaneously via their respective `.test` domains.

---

## Appendix A: Generating Local SSL Certificates

To ensure your browser trusts your local `.test` domains without throwing annoying security warnings, we use a tool called `mkcert`. It acts as a local Certificate Authority (CA) on your machine.

### A.1 Install mkcert
If you are on macOS, you can easily install it via Homebrew:

````bash
brew install mkcert
brew install nss # Optional, but recommended if you use Firefox locally
````

*(For Windows/Linux installation instructions, visit the official [mkcert GitHub repository](https://github.com/FiloSottile/mkcert)).*

### A.2 Install the Local CA
You only need to run this command **once** per machine to install the local CA into your system's root trust store:

````bash
mkcert -install
````

### A.3 Generate the Certificates
Navigate to the `certs` directory of your Master-Traefik and generate the certificates for your project. We include a wildcard (`*.`) so that subdomains like `mailhog.kasseapparat.test` are automatically covered by the same certificate:

````bash
cd ~/Projects/traefik-proxy/certs

mkcert -cert-file kasseapparat.test.crt -key-file kasseapparat.test.key "kasseapparat.test" "*.kasseapparat.test"
````

### A.4 Apply Changes
Whenever you generate certificates for a new project, make sure to add them to your `traefik-dynamic.yaml` as described in Step 1.4, and restart the Master-Traefik to apply the new configuration:

````bash
docker restart traefik_global
````

---

## Appendix B: Local DNS Resolution (Wildcard with dnsmasq)

By default, your operating system does not know how to route `.test` domains. While you could manually add every single subdomain to your `/etc/hosts` file, a much smarter and maintenance-free approach is to use `dnsmasq`. It acts as a local DNS resolver and automatically routes **all** `*.test` traffic to `127.0.0.1`.

### B.1 Install and Configure dnsmasq
Install it via Homebrew:

````bash
brew install dnsmasq
````

Now, tell `dnsmasq` to route all `.test` domains to your local machine. Open the configuration file:

````bash
nano $(brew --prefix)/etc/dnsmasq.conf
````

Add the following line at the very bottom of the file, save (`Ctrl + O`, `Enter`), and exit (`Ctrl + X`):

````text
address=/.test/127.0.0.1
````

### B.2 Start the Service
Start `dnsmasq` so it runs continuously in the background (requires `sudo` to bind to the DNS port 53):

````bash
sudo brew services start dnsmasq
````

### B.3 Configure the macOS Resolver
Now you need to tell macOS to actually *use* your new `dnsmasq` server whenever a `.test` domain is requested. macOS has a neat feature for this: custom resolvers per Top-Level-Domain (TLD).

Create the resolver directory (if it doesn't exist) and the specific config file for `.test`:

````bash
sudo mkdir -p /etc/resolver
sudo nano /etc/resolver/test
````

Add exactly this single line to the file, save, and exit:

````text
nameserver 127.0.0.1
````

### B.4 Verify
Flush your DNS cache once to ensure the new rules are applied immediately:

````bash
sudo dscacheutil -flushcache; sudo killall -HUP mDNSResponder
````

From now on, **any** domain ending in `.test` (e.g., `kasseapparat.test`, `mailhog.kasseapparat.test`, `anything.test`) will automatically resolve to your Master-Traefik without you ever having to touch a config file again!

> **Fallback:** If you cannot install dnsmasq, you will need to manually add every single subdomain like `127.0.0.1 kasseapparat.test mailhog.kasseapparat.test` to your `/etc/hosts` file.