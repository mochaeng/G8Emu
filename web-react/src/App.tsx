import { useEffect, useRef, useState } from "react";
import Emulator from "./components/Emulator";
import { ControlPanel } from "./components/ControlPanel";

import "@fontsource/nerko-one";

export default function App() {
  const [emulatorReady, setEmulatorReady] = useState(false);
  const emulatorRef = useRef<HTMLIFrameElement>(null);

  useEffect(() => {
    function handleMessage(event: MessageEvent) {
      if (
        emulatorRef.current &&
        event.source === emulatorRef.current.contentWindow &&
        event.data.type === "ready"
      ) {
        setEmulatorReady(true);
      }
    }

    window.addEventListener("message", handleMessage);
    return () => window.removeEventListener("message", handleMessage);
  }, []);

  const handleRomUpload = (file: File | null) => {
    if (!emulatorReady || !file || !emulatorRef.current) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      if (!e.target?.result || typeof e.target.result === "string") return;

      const romData = new Uint8Array(e.target.result);
      emulatorRef.current!.contentWindow!.postMessage(
        { type: "loadRom", data: romData },
        "*",
      );
    };

    reader.readAsArrayBuffer(file);
  };

  const handleReset = () => {
    if (!emulatorRef.current) return;
    emulatorRef.current.contentWindow!.postMessage({ type: "reset" }, "*");
  };

  const handlePause = () => {
    if (!emulatorRef.current) return;
    emulatorRef.current.contentWindow!.postMessage(
      { type: "togglePause" },
      "*",
    );
  };

  const handleCpuFrequencyChange = (value: string) => {
    if (!emulatorRef.current) return;
    emulatorRef.current.contentWindow!.postMessage(
      { type: "setFrequency", value: parseInt(value) },
      "*",
    );
  };

  return (
    <div className="min-h-screen bg-background text-primary p-4 sm:p-8">
      <header className="text-center mb-8 pb-6 border-b border-border/30">
        <h1 className="text-5xl sm:text-4xl font-bold text-primary">G8Emu</h1>
        <p className="text-primary/70 mt-2 font-medium text-2xl">
          CHIP-8 Web Emulator
        </p>
      </header>

      <main className="max-w-7xl mx-auto flex flex-col lg:flex-row  gap-8">
        <div className="bg-card rounded-xl flex-1 p-4 shadow-xl border border-border/10">
          <Emulator ref={emulatorRef} />
        </div>
        <ControlPanel
          onRomUpload={handleRomUpload}
          onReset={handleReset}
          onPause={handlePause}
          onCpuFrequencyChange={handleCpuFrequencyChange}
          disabled={!emulatorReady}
        />
      </main>

      <footer className="mt-12 text-center text-primary border-t border-border/30 pt-6 text-lg">
        <p>G8Emu - CHIP-8 Emulator | Built with Go + WebAssembly</p>
      </footer>
    </div>
  );
}
