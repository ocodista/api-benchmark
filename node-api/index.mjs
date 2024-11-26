import http from "node:http";
import pg from "pg";
const { Pool } = pg;

const pool = new Pool({
  host: process.env.POSTGRES_HOST,
  user: process.env.POSTGRES_USER,
  password: process.env.POSTGRES_PASSWORD,
  port: 5432,
  database: process.env.POSTGRES_DATABASE,
  max: 60,
  idleTimeoutMillis: 0,
  connectionTimeoutMillis: 30000,
  ssl:
    process.env.NODE_ENV === "development"
      ? false
      : {
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
    Connection: "keep-alive", // Default to keep-alive for persistent connections
    "Cache-Control": "no-store", // No caching for user creation
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
      const responseBody = JSON.stringify({ error: error.message });
      console.error(error);
      const statusCode = error instanceof SyntaxError ? 400 : 500;
      sendResponse(res, statusCode, headers, responseBody);
    }
  } else {
    headers["Content-Type"] = "text/plain";
    sendResponse(res, 404, headers, "Not Found!");
  }
});

const PORT = process.env.PORT || 3000;

console.log(`Process PID: ${process.pid}`);
server.listen(PORT, () => {
  console.log(`Server running on http://0.0.0.0:${PORT}`);
});
