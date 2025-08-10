const go = new Go();

document.addEventListener("DOMContentLoaded", () => {
  WebAssembly.instantiateStreaming(fetch("g8emu.wasm"), go.importObject)
    .then((result) => {
      go.run(result.instance);
      console.log("WASM initialized successfully");

      parent.postMessage({ type: "ready" }, "*");
    })
    .catch((err) => {
      console.log("WASM initialization failed: ", err);
      showError(err);
    });
});

window.addEventListener("message", (event) => {
  if (event.source !== parent) return;

  console.log("Message from parent:", event.data);

  switch (event.data.type) {
    case "loadRom":
      if (window.loadRom) {
        window.loadRom(event.data.data);
      } else {
        console.error("loadRom function not available");
      }
      break;

    case "reset":
      if (window.resetEmulator) {
        window.resetEmulator();
      }
      break;

    case "togglePause":
      if (window.togglePause) {
        window.togglePause();
      }
      break;

    case "setFrequency":
      if (window.setCpuFrequency) {
        window.setCpuFrequency(event.data.value);
      }
      break;
  }
});
