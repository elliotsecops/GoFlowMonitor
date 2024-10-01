# Web Application Performance Monitoring System

This project provides a minimalistic yet powerful solution for monitoring the performance of a web application. It uses Go for the metrics exporter, Prometheus for metrics collection and alerting, and Grafana for visualization and reporting.

## Project Structure

```
web_app_monitor/
├── cmd/
│   └── metrics_exporter/
│       └── main.go
├── config/
│   ├── prometheus.yml
│   └── alert.rules
├── dashboard/
│   └── web_app_dashboard.json
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

## Installation

### Prerequisites

- Docker
- Docker Compose
- Go (Golang)

### Steps

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/yourusername/web_app_monitor.git
   cd web_app_monitor
   ```

2. **Set Up Environment Variables:**
   - Copy `.env.example` to `.env` and update the values as needed.
   ```bash
   cp .env.example .env
   ```

3. **Build and Run the Metrics Exporter:**
   ```bash
   docker-compose up --build
   ```
   This command will build the necessary Docker image and start Prometheus, Grafana, and the metrics exporter.

4. **Access Prometheus and Grafana:**
   - Prometheus: `http://localhost:9090`
   - Grafana: `http://localhost:3000` (default credentials: admin/admin)

5. **Import the Grafana Dashboard:**
   - Open Grafana and import the dashboard from `dashboard/web_app_dashboard.json`.

## Usage

- **Metrics Exporter:** The Go script in `cmd/metrics_exporter/main.go` collects metrics from the web application and exposes them in Prometheus's exposition format.
- **Prometheus:** Prometheus collects the metrics and triggers alerts based on the rules defined in `config/prometheus.yml`.
- **Grafana:** Grafana visualizes the metrics using the dashboard defined in `dashboard/web_app_dashboard.json`.

## Advanced Features (Optional)

### Dynamic Configuration

**Reload Endpoint:** Implement a `/reload` endpoint to reload the configuration without restarting the Go exporter. This feature adds complexity and is optional.

### Custom Metrics

**User-Defined Metrics:** Allow users to define custom metrics in the configuration file and dynamically add them to the exporter. This is a good-to-have feature but not strictly essential for the core functionality.

## Alertmanager (Optional)

For a more complete alerting solution, you can set up Alertmanager. Alertmanager handles alert deduplication, silencing, and routing to different notification channels (email, Slack, etc.).

## Troubleshooting Guide

### Common Issues and Solutions

1. **Prometheus Not Scraping Metrics:**
   - **Check Prometheus Configuration:** Ensure that the `prometheus.yml` file is correctly configured with the correct target URL.
   - **Network Issues:** Verify that Prometheus can reach the Go exporter on the specified port.

2. **Grafana Dashboard Not Displaying Data:**
   - **Data Source Configuration:** Ensure that the Prometheus data source is correctly configured in Grafana.
   - **Permissions:** Check that the Grafana user has the necessary permissions to access the Prometheus data source.

3. **Go Exporter Failing to Start:**
   - **Environment Variables:** Ensure that all required environment variables are set correctly in the `.env` file.
   - **Port Conflicts:** Verify that the port specified in the `.env` file (default is 8080) is not already in use by another service.

4. **High CPU/Memory Usage:**
   - **Resource Limits:** Ensure that Docker containers have appropriate resource limits set.
   - **Optimize Go Exporter:** Review the Go exporter code for any inefficiencies and optimize as necessary.

### Logs and Diagnostics

- **Go Exporter Logs:** Check the logs of the Go exporter container for any errors or warnings.
- **Prometheus Logs:** Check the logs of the Prometheus container for any scraping errors.
- **Grafana Logs:** Check the logs of the Grafana container for any data source connection issues.

### Support and Community

If you encounter issues that are not covered in this guide, please open an issue on the project's GitHub repository. Contributions and feedback are always welcome!

## License

This project is licensed under the MIT License.
