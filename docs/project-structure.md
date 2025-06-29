project-root/
├── backend/
│   ├── cmd/
│   │   └── main.go                 # Gin app entrypoint
│   ├── internal/
│   │   ├── api/                    # Route handlers (resize, batch, etc.)
│   │   ├── auth/                   # API key authentication middleware
│   │   ├── jobs/                   # Job and progress tracking (Redis)
│   │   ├── s3/                     # S3 storage functions (upload, delete, etc.)
│   │   ├── imgproc/                # Image processing: resize, center-crop, etc.
│   │   └── config/                 # App config (env vars, etc.)
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
├── frontend/
│   ├── public/                     # Svelte static assets
│   ├── src/                        # Svelte components, routes, stores
│   ├── svelte.config.js
│   ├── package.json
│   ├── vite.config.js
│   └── Dockerfile
├── nginx/
│   ├── nginx.conf                  # Nginx config for serving frontend
│   └── Dockerfile                  # Nginx Docker image for static site
├── deploy/
│   ├── docker-compose.yml          # Compose file for all services
│   ├── .env.example                # Example environment variables (S3, Redis, API keys)
│   └── README.md                   # Deployment instructions
├── docs/
│   ├── ARCHITECTURE_AND_FEATURES.md
│   └── API.md                      # (To be filled: endpoint details, etc)
└── README.md