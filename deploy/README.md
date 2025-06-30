# Deployment Guide for Image Resizer App

This directory contains the configuration files for deploying the Image Resizer App using Docker Compose.

## Prerequisites

1. Install [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/).
2. Ensure ports `8080`, `8081`, `6379`, `9000`, and `9001` are available on your system.

## Setup

1. Copy the `.env.example` file to `.env`:
   ```bash
   cp .env.example .env
   ```

2. Update the `.env` file with your desired configuration:
   - Replace `yourapikey123,anotherkey456` with your actual API keys.
   - Update any other environment variables as needed.

## Running the Application

1. Navigate to the `deploy` directory:
   ```bash
   cd deploy
   ```

2. Start the application using Docker Compose:
   ```bash
   docker-compose up --build
   ```

3. Open your browser and visit:
   - Frontend: [http://localhost:8080](http://localhost:8080)
   - Backend Health Check: [http://localhost:8081/health](http://localhost:8081/health)

## Stopping the Application

To stop the application and remove containers:
```bash
docker-compose down
```

## Services

- **Frontend**: Svelte app served via Nginx.
- **Backend**: Go-based API for image processing.
- **Redis**: Used for job tracking.
- **Minio**: S3-compatible storage for processed images.
- **Nginx**: Serves the frontend and proxies API requests to the backend.

## Environment Variables

### Backend
- `API_KEYS`: Comma-separated list of API keys for authenticating requests.
- `REDIS_ADDR`: Address of the Redis server.
- `S3_ENDPOINT`: Endpoint for the S3-compatible storage (e.g., Minio).
- `S3_ACCESS_KEY`: Access key for S3.
- `S3_SECRET_KEY`: Secret key for S3.
- `S3_BUCKET`: Name of the S3 bucket for storing images.

### Minio
- `MINIO_ROOT_USER`: Root username for Minio.
- `MINIO_ROOT_PASSWORD`: Root password for Minio.

## Volumes

- `minio-data`: Persistent storage for Minio.

## Logs

To view logs for a specific service:
```bash
docker-compose logs <service-name>
```

For example:
```bash
docker-compose logs backend
```

## Troubleshooting

1. **Frontend Not Loading**:
   - Ensure the frontend is built and the `dist/` folder is copied to the Nginx container.

2. **API Requests Failing**:
   - Check the backend logs:
     ```bash
     docker-compose logs backend
     ```

3. **Redis or Minio Issues**:
   - Check the respective service logs:
     ```bash
     docker-compose logs redis
     docker-compose logs minio
     ```

## Cleanup

To remove all containers, networks, and volumes:
```bash
docker-compose down -v
```
```

---

### **3. Next Steps**
1. Update the `.env` file with your actual API keys and other configurations.
2. Add the above README.md to the deploy directory.
3. Test the deployment by running:
   ```bash
   cd deploy
   docker-compose up --build
   ```
4. Let me know if you encounter any issues or need further assistance!