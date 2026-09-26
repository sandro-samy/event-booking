# Docker Guide

This guide explains, step by step, how to package this app into a Docker
container and run it — no prior Docker experience assumed.

## 1. What is Docker, in one sentence?

Docker packages your app + everything it needs to run (like a mini
operating system) into a single "image", so it runs the same way on any
computer, without needing Go installed.

## 2. Files added for this

- **Dockerfile** — the recipe that tells Docker how to build our app's image.
- **.dockerignore** — tells Docker which files to skip when building (like `.gitignore`, but for Docker).

## 3. Install Docker (one-time setup)

Download and install **Docker Desktop**: https://www.docker.com/products/docker-desktop/

After installing, open a terminal and check it worked:

```bash
docker --version
```

## 4. Build the image

From the project's root folder (where the `Dockerfile` is), run:

```bash
docker build -t event-booking .
```

What this means:
- `docker build` = build an image.
- `-t event-booking` = give the image a name ("tag") so we can refer to it later.
- `.` = look for the `Dockerfile` in the current folder.

This will take a minute the first time (it downloads Go and your
dependencies). You'll see it printing each step from the `Dockerfile`.

## 5. Run the container

```bash
docker run -p 8080:8080 --env-file .env event-booking
```

What this means:
- `docker run` = start a container from the image we built.
- `-p 8080:8080` = connect port 8080 on your computer to port 8080 inside the container (that's the port our app listens on, see `main.go`).
- `--env-file .env` = pass in our environment variables (like `JWT_SECRET`) from the `.env` file, since `.env` is never copied into the image.
- `event-booking` = the name of the image to run.

Now open your browser or Postman and hit `http://localhost:8080` — it's
the same as running `go run main.go` locally, just inside a container.

Press `Ctrl+C` in the terminal to stop it.

## 6. A note about the database

This app uses SQLite, which stores everything in a single file (`api.db`).
By default, that file lives *inside* the container, so if you delete the
container, you lose the data. For local testing that's fine. If you want
the data to survive container restarts, mount a folder from your computer:

```bash
docker run -p 8080:8080 --env-file .env -v "$(pwd)/data:/app:rw" event-booking
```

This tells Docker: "store the `api.db` file in a `data` folder on my
computer, not just inside the throwaway container."

## 7. Useful everyday commands

```bash
docker ps                # see running containers
docker stop <container_id>   # stop a running container
docker images             # see images you've built
docker rmi event-booking  # delete the image
```

## 8. Deploying it (the very basic version)

"Deploying" just means running the same `docker run` command, but on a
server instead of your laptop. The simplest path:

1. **Get a small cloud server** (e.g. an AWS EC2 instance, a DigitalOcean
   Droplet, or similar) with Docker installed on it.
2. **Get your code onto that server.** Easiest way: push your code to
   GitHub, then on the server run `git clone <your-repo-url>`.
3. **On the server**, run the same two commands from steps 4 and 5 above:
   ```bash
   docker build -t event-booking .
   docker run -d -p 8080:8080 --env-file .env event-booking
   ```
   (Notice the `-d` flag here — it means "detached", so the container
   keeps running in the background even after you close your terminal.)
4. **Open the port** on your server's firewall/security group so port
   8080 is reachable from the internet.
5. Visit `http://<your-server-ip>:8080` from any browser.

That's the whole idea: build once, then run the exact same image
anywhere Docker is installed — your laptop, a teammate's laptop, or a
cloud server — and it behaves identically every time.

### A slightly nicer alternative: push to Docker Hub

Instead of copying source code to the server and building there, you can
build the image once and share it via Docker Hub (like GitHub, but for
Docker images):

```bash
docker tag event-booking <your-dockerhub-username>/event-booking
docker push <your-dockerhub-username>/event-booking
```

Then on the server, instead of cloning + building, you just pull and run:

```bash
docker pull <your-dockerhub-username>/event-booking
docker run -d -p 8080:8080 --env-file .env <your-dockerhub-username>/event-booking
```
