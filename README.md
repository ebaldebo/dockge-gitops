# Dockge-gitops

Add gitops functionality to [Dockge](https://github.com/louislam/dockge)

### Background
Switching from [Portainer](https://www.portainer.io/) to [Dockge](https://github.com/louislam/dockge) I was missing the ability to have a gitops workflow. I wanted to be able to make changes to my docker-compose.yml files and have them automatically applied to my Dockge stacks. This project is a simple temporary solution to that problem.

## Features
- Clone repo into [Dockge](https://github.com/louislam/dockge) stacks directory
  - Access private repos using [PAT](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token)
- Automatically poll for changes on a configurable interval

## Usage
### Environment Variables
| Variable | Description | Default | Required |
| --- | --- | --- | --- |
| REPO_URL | URL of the git repo to clone | | Yes |
| PAT | Personal Access Token for private repos | | No |
| POLLING_RATE | How often to poll for changes | 5m | No |
| DOCKGE_STACKS_DIR | Path to Dockge stacks directory | /opt/stacks | No |


Example `compose.yml`:
```yaml
version: '3.8'

services:
  dockge-gitops:
    image: ghcr.io/ebaldebo/dockge-gitops:latest
    container_name: dockge-gitops
    restart: unless-stopped
    environment:
      - REPO_URL=https://github.com/author/example.git # required
      - PAT=ghp_iamapat123 # optional
      - POLLING_RATE=10m # optional
      # Needs to be the same as the stacks directory in the dockge container
      - DOCKGE_STACKS_DIR=/opt/stacks 
    volumes:
      # Needs to be the same as the stacks directory in the dockge container
      - /opt/stacks:/opt/stacks

  dockge:
    image: louislam/dockge:1
    container_name: dockge
    restart: unless-stopped
    ports:
      - 5001:5001
    environment:
      - DOCKGE_STACKS_DIR=/opt/stacks
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - .data:/app/data
      - /opt/stacks:/opt/stacks
``````
