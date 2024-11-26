import cluster from "cluster";
import os from "os";
import http from "node:http";
import pg from "pg";

console.log("CPUS " + os.cpus().length);

if (cluster.isPrimary) {
  console.log(`Master ${process.pid} is running`);

  // Fork workers.
  for (let i = 0; i < os.cpus().length; i++) {
    cluster.fork();
  }

  cluster.on("exit", (worker) => {
    console.log(`Worker ${worker.process.pid} died`);
    cluster.fork(); // Optional: Restart a new worker upon death
  });
} else {
  const { Pool } = pg;
  const pool = new Pool({
    host: process.env.POSTGRES_HOST,
    user: process.env.POSTGRES_USER,
    password: process.env.POSTGRES_PASSWORD,
    port: 5432,
    database: process.env.POSTGRES_DATABASE,
    max: 60 / os.cpus().length,
    idleTimeoutMillis: 0,
    connectionTimeoutMillis: 30000,
    ssl: {
      rejectUnauthorized: false,
    },
  });

  const createUser = async (email, password) => {
    const queryText =
      "INSERT INTO users(email, password) VALUES($1, $2) RETURNING id";
    const { rows } = await pool.query(queryText, [email, password]);
    return rows[0].id;
  };

  const getRequestBody = (req) =>
    new Promise((resolve, reject) => {
      let body = "";
      req.on("data", (chunk) => (body += chunk.toString()));
      req.on("end", () => resolve(body));
      req.on("error", (err) => reject(err));
    });

  const sendResponse = (res, statusCode, headers, body) => {
    headers["Content-Length"] = Buffer.byteLength(body).toString();
    res.writeHead(statusCode, headers);
    res.end(body);
  };

  const server = http.createServer(async (req, res) => {
    const headers = {
      "Content-Type": "application/json",
      Connection: "keep-alive",
      "Cache-Control": "no-store",
    };
    if (req.method === "POST" && req.url === "/user") {
      try {
        const body = await getRequestBody(req);
        const { email, password } = JSON.parse(body);
        const userId = await createUser(email, password);

        headers["Location"] = `/user/${userId}`;
        const responseBody = JSON.stringify({ message: "User created" });
        sendResponse(res, 201, headers, responseBody);
      } catch (error) {
        headers["Connection"] = "close";
        headers["Content-Type"] = "text/plain";
        const statusCode = error instanceof SyntaxError ? 400 : 500;
        console.log(error);
        sendResponse(res, statusCode, headers, "");
      }
    } else {
      headers["Content-Type"] = "text/plain";
      sendResponse(res, 404, headers, "Not Found!");
    }
  });

  const PORT = process.env.PORT || 3000;

  server.listen(PORT, () => {
    console.log(`Worker ${process.pid} started on http://localhost:${PORT}`);
  });

  // Function to handle graceful shutdown
  const gracefulShutdown = async () => {
    console.log(`Worker ${process.pid} is shutting down`);
    server.close(() => {
      console.log(`HTTP server closed for worker ${process.pid}`);
    });

    try {
      await pool.end(); // Close the pool
      console.log(`Database pool closed for worker ${process.pid}`);
    } catch (error) {
      console.error(
        `Error closing database pool for worker ${process.pid}:`,
        error,
      );
    }

    process.exit(0); // Exit the process
  };

  // Listen for termination signals
  process.on("SIGTERM", gracefulShutdown);
  process.on("SIGINT", gracefulShutdown);
}
