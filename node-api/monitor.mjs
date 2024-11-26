import { exec } from "child_process";

const getOpenFileDescriptors = (pid) => {
  const command = `ls /proc/${pid}/fd | wc -l`;

  exec(command, (error, stdout, stderr) => {
    if (error) {
      console.error(`exec error: ${error}`);
      return;
    }
    if (stderr) {
      console.error(`stderr: ${stderr}`);
      return;
    }
    console.log(`Number of open file descriptors: ${stdout.trim()}`);
  });
};

const pid = process.argv[2];

if (!pid) {
  console.error("Please provide a PID.");
  process.exit(1);
}

const interval = setInterval(() => getOpenFileDescriptors(pid), 500);

process.on("SIGINT", () => {
  console.log("End");
  clearInterval(interval);
  process.exit();
});
