import { useEffect, useRef, useState } from "react";
import Emulator from "./components/Emulator";
import { ControlPanel } from "./components/ControlPanel";

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
    <div className="min-h-screen bg-[#FBF6E2] text-[#131842] p-4 sm:p-8">
      <header className="text-center mb-8 pb-6 border-b border-[#131842]/20">
        <h1 className="text-3xl sm:text-4xl font-bold text-[#131842]">G8Emu</h1>
        <p className="text-[#131842]/70 mt-2 font-medium">
          CHIP-8 Web Emulator
        </p>
      </header>

      <main className="max-w-7xl mx-auto flex flex-col lg:flex-row  gap-8">
        <div className="bg-[#ECCEAE] rounded-xl flex-1 p-4 shadow-xl border border-[#131842]/10">
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

      <footer className="mt-12 text-center text-[#131842]/60 border-t border-[#131842]/20 pt-6">
        <p>G8Emu - CHIP-8 Emulator | Built with Go + WebAssembly</p>
      </footer>
    </div>
  );
}
