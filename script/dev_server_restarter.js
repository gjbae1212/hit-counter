const { spawnSync, spawn } = require("child_process");
const fs = require("fs");
const path = require("path");

const env = {
	PHASE: "local",
	REDIS_ADDRS: "localhost:6379",
};

function throttle(fn, duration) {
	return (function () {
		let timeout;
		return function (...args) {
			if (!timeout) {
				fn(...args);
				timeout = setTimeout(() => {
					timeout = null;
				}, duration);
			}
		};
	})();
}

function runServer() {
	spawnSync(path.join(__dirname, "make_wasm \"local\""))
	spawnSync("go build", { stdio: "inherit", env: { ...process.env, ...env } });
	return spawn(path.join(__dirname, "..","hit-counter"), ["-tls=0", "-addr=:8080"], {
		stdio: "inherit",
		env: { ...process.env, ...env },
	});
}

function main() {
	let devServer;
	devServer = runServer();
	fs.watch(
		path.join(__dirname, "../view"),
		throttle(() => {
			if (devServer) {
				devServer.kill();
			}
			devServer = runServer();
		}, 1000)
	);
}

main();