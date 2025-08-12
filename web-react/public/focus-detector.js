(function () {
  "use strict";

  function setupCanvasFocusDetection() {
    const canvas = document.querySelector("canvas");

    if (!canvas) {
      setTimeout(setupCanvasFocusDetection, 100);
      return;
    }

    canvas.addEventListener("focus", () => {
      window.parent.postMessage({ type: "canvas-focus" }, "*");
    });

    canvas.addEventListener("blur", () => {
      window.parent.postMessage({ type: "canvas-blur" }, "*");
    });

    canvas.addEventListener("click", () => {
      canvas.focus();
      window.parent.postMessage({ type: "canvas-focus" }, "*");
    });

    if (document.activeElement === canvas) {
      window.parent.postMessage({ type: "canvas-focus" }, "*");
    }
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", setupCanvasFocusDetection);
  } else {
    setupCanvasFocusDetection();
  }
})();
