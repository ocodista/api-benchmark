# Go vs Node Experiment Benchmark
A comparative study of Go and Node.js performance under high-load conditions, focusing on a simple POST-based create-user endpoint inserting data into a PostgreSQL database.

## Overview
This experiment benchmarks Go and Node.js APIs to understand their behavior under different load scenarios. It measures key performance metrics like latency, CPU and RAM usage, thread count, and file descriptor management.

## Setup
- **Infrastructure**: AWS EC2 instances, Postgres RDS.
- **Load Testing**: Vegeta for HTTP load testing.
- **Languages**: Go 1.21.4 and Node.js 21.4.0.
- **Monitoring**: Custom scripts for capturing performance metrics.

## Running the Experiment
### On AWS
#### Prerequisites:
- AWS account with [CLI](https://aws.amazon.com/cli/) configured.
- [OpenTofu](https://opentofu.org/) for infrastructure setup.
- PostgreSQL RDS instance.

#### Steps:
1. **Infrastructure Setup**:
   ```
   cd tofu
   tofu apply -auto-approve
   ```
   This creates the required AWS infrastructure including two Ubuntu servers.

2. **Running the API**:
   SSH into the API server:
   ```
   ./ssh_connect_api.sh
   ```
   For Node API:
   ```
   cd node-api
   npm install
   npm start
   ```
   For Go API:
   ```
   cd go-api
   go build -o api
   ./api
   ```
   Note the process ID (PID) for monitoring.

3. **Performance Monitoring**:
   Reconnect to the API server and run:
   ```
   ./monitor_process.sh [PID] [Request Rate]
   ```
   Replace `[PID]` with the actual process ID and `[Request Rate]` with the desired requests per second.

4. **Load Test**:
   Connect to the gun server:
   ```
   ./ssh_connect_gun.sh
   ```
   Start the stress test:
   ```
   cd load-tester/vegeta
   ./metrics.sh 2000
   ```

### Locally with Docker
- Use `docker-compose` to start services:
  ```
  docker-compose up -d postgres node-api go-api
  ```
  ![Services Startup](https://github.com/ocodista/api-benchmark/assets/19851187/0aad0411-d171-415e-b2fd-c6c8cbad2222)

- Then, to run the gun server:
  ```
  docker-compose up gun
  ```
  After the test, you can see the core metrics on the terminal.
  ![Terminal Output](https://github.com/ocodista/api-benchmark/assets/19851187/50b146d4-201a-42fc-82f2-7167d1a3d82e)

## Considerations
- This benchmark focuses on specific aspects of performance under load. It shouldn't be the sole basis for choosing a technology stack.
- Factors like productivity, team expertise, and existing toolsets are also crucial in technology decisions.
- In Docker Compose, prefer `network_mode: host`, as the docker internal networking adds considerable bottleneck
- Adding Node.js Cluster can handle more requests/sec but also adds extra networking and coordination overhead
