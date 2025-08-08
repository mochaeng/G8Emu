import React, { useState, useRef, useEffect } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Play, Pause, RotateCcw, Upload, Cpu } from "lucide-react";

export default function G8EmuEmulator() {
  const [isPaused, setIsPaused] = useState(false);
  const [frequency, setFrequency] = useState("540");
  const [emulatorReady, setEmulatorReady] = useState(false);
  const emulatorRef = useRef(null);
  const fileInputRef = useRef(null);

  useEffect(() => {
    // Listen for messages from the emulator iframe
    const handleMessage = (event) => {
      if (event.source !== emulatorRef.current?.contentWindow) return;

      console.log("Message from emulator:", event.data);
      if (event.data.type === "ready") {
        console.log("Emulator is ready");
        setEmulatorReady(true);
      }
    };

    window.addEventListener("message", handleMessage);
    return () => window.removeEventListener("message", handleMessage);
  }, []);

  const handleRomUpload = (e) => {
    const file = e.target.files[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (event) => {
      const romData = new Uint8Array(event.target.result);
      if (emulatorRef.current) {
        emulatorRef.current.contentWindow.postMessage(
          {
            type: "loadRom",
            data: romData,
          },
          "*",
        );
      }
    };
    reader.readAsArrayBuffer(file);
  };

  const handleReset = () => {
    if (emulatorRef.current) {
      emulatorRef.current.contentWindow.postMessage({ type: "reset" }, "*");
    }
  };

  const handlePause = () => {
    if (emulatorRef.current) {
      emulatorRef.current.contentWindow.postMessage(
        { type: "togglePause" },
        "*",
      );
      setIsPaused(!isPaused);
    }
  };

  const handleFrequencyChange = (value) => {
    setFrequency(value);
    if (emulatorRef.current) {
      emulatorRef.current.contentWindow.postMessage(
        {
          type: "setFrequency",
          value: parseInt(value),
        },
        "*",
      );
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 text-green-400 font-mono">
      {/* Retro grid background */}
      <div className="fixed inset-0 bg-[linear-gradient(rgba(0,255,0,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(0,255,0,0.03)_1px,transparent_1px)] bg-[size:20px_20px] pointer-events-none"></div>

      <div className="relative z-10 container mx-auto px-4 py-8 max-w-4xl">
        {/* Header */}
        <header className="text-center mb-12">
          <h1 className="text-6xl font-bold mb-2 text-transparent bg-clip-text bg-gradient-to-r from-green-400 via-cyan-400 to-green-400 animate-pulse">
            G8EMU
          </h1>
          <div className="text-cyan-300 text-xl tracking-wider">
            {">"} CHIP-8 WEB EMULATOR {"<"}
          </div>
          <div className="mt-4 text-green-500/60 text-sm">
            {emulatorReady ? "● SYSTEM READY" : "○ INITIALIZING..."}
          </div>
        </header>

        <div className="grid gap-8 lg:grid-cols-1">
          {/* Emulator Display */}
          <Card className="bg-black/50 border-green-500/30 backdrop-blur-sm">
            <CardHeader className="pb-4">
              <CardTitle className="text-green-400 flex items-center gap-2">
                <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
                DISPLAY TERMINAL
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="relative">
                <iframe
                  ref={emulatorRef}
                  src="emulator.html"
                  title="CHIP-8 Emulator"
                  className="w-full h-96 bg-black border border-green-500/50 rounded"
                />
                <div className="absolute -inset-1 bg-green-400/20 blur rounded animate-pulse opacity-50"></div>
              </div>
            </CardContent>
          </Card>

          {/* Controls */}
          <Card className="bg-black/50 border-cyan-500/30 backdrop-blur-sm">
            <CardHeader>
              <CardTitle className="text-cyan-400 flex items-center gap-2">
                <Cpu className="w-5 h-5" />
                CONTROL PANEL
              </CardTitle>
              <CardDescription className="text-cyan-300/60">
                Configure emulator parameters
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              {/* ROM Upload */}
              <div className="space-y-2">
                <Label className="text-green-400 flex items-center gap-2">
                  <Upload className="w-4 h-4" />
                  ROM UPLOAD
                </Label>
                <Input
                  ref={fileInputRef}
                  type="file"
                  accept=".ch8"
                  onChange={handleRomUpload}
                  className="bg-black/70 border-green-500/50 text-green-300 file:bg-green-500/20 file:text-green-400 file:border-0 file:rounded file:px-3 file:py-1 file:mr-4 hover:border-green-400/70 transition-colors"
                />
              </div>

              {/* CPU Frequency */}
              <div className="space-y-2">
                <Label className="text-green-400">CPU FREQUENCY</Label>
                <Select value={frequency} onValueChange={handleFrequencyChange}>
                  <SelectTrigger className="bg-black/70 border-green-500/50 text-green-300 hover:border-green-400/70">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent className="bg-black border-green-500/50">
                    <SelectItem
                      value="250"
                      className="text-green-300 hover:bg-green-500/20"
                    >
                      SLOW (250 Hz)
                    </SelectItem>
                    <SelectItem
                      value="540"
                      className="text-green-300 hover:bg-green-500/20"
                    >
                      NORMAL (540 Hz)
                    </SelectItem>
                    <SelectItem
                      value="1000"
                      className="text-green-300 hover:bg-green-500/20"
                    >
                      FAST (1000 Hz)
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Action Buttons */}
              <div className="grid grid-cols-2 gap-4">
                <Button
                  onClick={handleReset}
                  variant="outline"
                  className="bg-black/70 border-yellow-500/50 text-yellow-400 hover:bg-yellow-500/20 hover:border-yellow-400 transition-all duration-200"
                >
                  <RotateCcw className="w-4 h-4 mr-2" />
                  RESET
                </Button>
                <Button
                  onClick={handlePause}
                  variant="outline"
                  className={`bg-black/70 transition-all duration-200 ${
                    isPaused
                      ? "border-green-500/50 text-green-400 hover:bg-green-500/20 hover:border-green-400"
                      : "border-red-500/50 text-red-400 hover:bg-red-500/20 hover:border-red-400"
                  }`}
                >
                  {isPaused ? (
                    <>
                      <Play className="w-4 h-4 mr-2" />
                      RESUME
                    </>
                  ) : (
                    <>
                      <Pause className="w-4 h-4 mr-2" />
                      PAUSE
                    </>
                  )}
                </Button>
              </div>
            </CardContent>
          </Card>

          {/* Keyboard Mapping */}
          <Card className="bg-black/50 border-purple-500/30 backdrop-blur-sm">
            <CardHeader>
              <CardTitle className="text-purple-400">
                KEYBOARD MAPPING
              </CardTitle>
              <CardDescription className="text-purple-300/60">
                CHIP-8 keypad to keyboard layout
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-8 text-sm">
                <div className="space-y-2">
                  <div className="text-purple-400 font-bold">
                    ORIGINAL CHIP-8
                  </div>
                  <pre className="text-purple-300 font-mono text-xs leading-relaxed bg-black/50 p-3 rounded border border-purple-500/30">
                    {`╭───┬───┬───┬───╮
│ 1 │ 2 │ 3 │ C │
├───┼───┼───┼───┤
│ 4 │ 5 │ 6 │ D │
├───┼───┼───┼───┤
│ 7 │ 8 │ 9 │ E │
├───┼───┼───┼───┤
│ A │ 0 │ B │ F │
╰───┴───┴───┴───╯`}
                  </pre>
                </div>
                <div className="space-y-2">
                  <div className="text-purple-400 font-bold">KEYBOARD</div>
                  <pre className="text-purple-300 font-mono text-xs leading-relaxed bg-black/50 p-3 rounded border border-purple-500/30">
                    {`╭───┬───┬───┬───╮
│ 1 │ 2 │ 3 │ 4 │
├───┼───┼───┼───┤
│ Q │ W │ E │ R │
├───┼───┼───┼───┤
│ A │ S │ D │ F │
├───┼───┼───┼───┤
│ Z │ X │ C │ V │
╰───┴───┴───┴───╯`}
                  </pre>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Footer */}
        <footer className="text-center mt-16 pt-8 border-t border-green-500/20">
          <div className="text-green-400/60 text-sm">
            {">"} G8EMU - CHIP-8 EMULATOR {"<"} | BUILT WITH GO + WEBASSEMBLY
          </div>
          <div className="mt-2 text-xs text-green-500/40">
            REACT + TAILWINDCSS + SHADCN/UI
          </div>
        </footer>
      </div>
    </div>
  );
}
