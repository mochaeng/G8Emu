import { Button } from "./ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { Label } from "./ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "./ui/select";

export function ControlPanel({
  onRomUpload,
  onReset,
  onPause,
  onCpuFrequencyChange,
  disabled,
}: {
  onRomUpload: (file: File | null) => void;
  onReset: () => void;
  onPause: () => void;
  onCpuFrequencyChange: (value: string) => void;
  disabled: boolean;
}) {
  return (
    <Card className="bg-[#ECCEAE] border-[#131842]/20">
      <CardContent className="p-6 space-y-6">
        {/* ROM Upload */}
        <div className="space-y-2">
          <Label htmlFor="rom-upload" className="text-[#131842] font-medium">
            Upload ROM
          </Label>
          <input
            id="rom-upload"
            type="file"
            accept=".ch8"
            onChange={(e) => onRomUpload(e.target.files?.[0] || null)}
            disabled={disabled}
            className="flex h-10 w-full rounded-md border border-[#131842]/30 bg-[#FBF6E2] px-3 py-2 text-sm text-[#131842] file:mr-4 file:py-1 file:px-4 file:rounded file:border-0 file:bg-[#E68369] file:text-white hover:file:bg-[#E68369]/80 disabled:opacity-50 focus:border-[#E68369] focus:ring-1 focus:ring-[#E68369] focus:outline-none"
          />
        </div>

        {/* CPU Frequency */}
        <div className="space-y-2">
          <Label className="text-[#131842] font-medium">CPU Speed</Label>
          <Select onValueChange={onCpuFrequencyChange} disabled={disabled}>
            <SelectTrigger className="bg-[#FBF6E2] border-[#131842]/30 text-[#131842] focus:border-[#E68369] focus:ring-1 focus:ring-[#E68369]">
              <SelectValue placeholder="Select CPU speed" />
            </SelectTrigger>
            <SelectContent className="bg-[#FBF6E2] border-[#131842]/30 text-[#131842]">
              <SelectItem value="250">Slow (250 Hz)</SelectItem>
              <SelectItem value="540">Normal (540 Hz)</SelectItem>
              <SelectItem value="1000">Fast (1000 Hz)</SelectItem>
            </SelectContent>
          </Select>
        </div>

        {/* Buttons */}
        <div className="grid grid-cols-2 gap-4">
          <Button
            onClick={onReset}
            disabled={disabled}
            className="bg-[#E68369] hover:bg-[#E68369]/80 text-white border-0 font-medium"
          >
            Reset Emulator
          </Button>
          <Button
            onClick={onPause}
            disabled={disabled}
            className="bg-[#131842] hover:bg-[#131842]/80 text-white border-0 font-medium"
          >
            Pause
          </Button>
        </div>

        {/* Keymap */}
        <Card className="bg-[#FBF6E2] border-[#131842]/20">
          <CardHeader>
            <CardTitle className="text-[#131842] text-lg font-semibold">
              Keyboard Mapping
            </CardTitle>
          </CardHeader>
          <CardContent>
            <pre className="bg-[#131842] p-4 rounded-md text-sm text-[#FBF6E2] overflow-x-auto border border-[#131842]/30">
              {`Original       Keyboard
  ╭───┬───┬───┬───╮     ╭───┬───┬───┬───╮
  │ 1 │ 2 │ 3 │ C │     │ 1 │ 2 │ 3 │ 4 │
  │ 4 │ 5 │ 6 │ D │     │ Q │ W │ E │ R │
  │ 7 │ 8 │ 9 │ E │     │ A │ S │ D │ F │
  │ A │ 0 │ B │ F │     │ Z │ X │ C │ V │
  ╰───┴───┴───┴───╯     ╰───┴───┴───┴───╯`}
            </pre>
          </CardContent>
        </Card>
      </CardContent>
    </Card>
  );
}
