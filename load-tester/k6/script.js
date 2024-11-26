import http from "k6/http";
import { check } from "k6";
import { Trend } from "k6/metrics";

export let options = {
  stages: [
    { duration: "10s", target: 10000 },
    { duration: "30s", target: 10000 },
    { duration: "10s", target: 20000 },
    { duration: "30s", target: 20000 },
  ],
};

const BASE_URL = "http://15.228.127.94:3000/user"; // Get the HOST_URL from environment variable
const PASSWORD = "fixedPassword123"; // Fixed password for all requests

function randomEmail() {
  const prefix = Math.random().toString(36).substring(2, 15);
  const suffix = Math.random().toString(36).substring(2, 15);
  return `${prefix}@${suffix}.com`;
}

let successRate = new Trend("success_rate");

export default function () {
  const payload = JSON.stringify({
    email: randomEmail(),
    password: PASSWORD,
  });

  const params = {
    headers: {
      "Content-Type": "application/json",
    },
    timeout: "10s", // Setting a timeout of 10 seconds for the requests
  };

  const response = http.post(BASE_URL, payload, params);

  let success = check(response, {
    "is status 200": (r) => r.status >= 200 && r.status <= 300,
    // Additional check for response timing
  });

  successRate.add(Number(success));
}
