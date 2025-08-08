document.addEventListener("DOMContentLoaded", () => {
  document
    .getElementById("rom-upload")
    .addEventListener("change", handleRomUpload);

  document.getElementById("reset-btn").addEventListener("click", handleReset);
  document.getElementById("pause-btn").addEventListener("click", handlePause);

  document
    .getElementById("cpu-frequency")
    .addEventListener("change", handleCpuFrequencyChange);

  const emulatorFrame = document.getElementById("emulator");

  window.addEventListener("message", (event) => {
    if (event.source !== emulatorFrame.contentWindow) return;

    console.log("Message from emulator: ", event.data);
    if (event.data.type === "ready") {
      console.log("Emulator is ready");
    }
  });
});

function handleRomUpload(e) {
  const file = e.target.files[0];
  if (!file) return;

  const reader = new FileReader();
  reader.onload = (event) => {
    const romData = new Uint8Array(event.target.result);
    const emulatorFrame = document.getElementById("emulator");

    emulatorFrame.contentWindow.postMessage(
      {
        type: "loadRom",
        data: romData,
      },
      "*",
    );
  };

  reader.readAsArrayBuffer(file);
}

function handleReset() {}

function handlePause() {}

function handleCpuFrequencyChange() {}
