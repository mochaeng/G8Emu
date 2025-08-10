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
    <Card className="bg-gray-900 border-gray-700">
      <CardContent className="p-6 space-y-6">
        {/* ROM Upload */}
        <div className="space-y-2">
          <Label htmlFor="rom-upload" className="text-blue-400">
            Upload ROM
          </Label>
          <input
            id="rom-upload"
            type="file"
            accept=".ch8"
            onChange={(e) => onRomUpload(e.target.files?.[0] || null)}
            disabled={disabled}
            className="flex h-10 w-full rounded-md border border-gray-600 bg-gray-800 px-3 py-2 text-sm text-gray-300 file:mr-4 file:py-1 file:px-4 file:rounded file:border-0 file:bg-blue-600 file:text-white hover:file:bg-blue-700 disabled:opacity-50"
          />
        </div>

        {/* CPU Frequency */}
        <div className="space-y-2">
          <Label className="text-blue-400">CPU Speed</Label>
          <Select onValueChange={onCpuFrequencyChange} disabled={disabled}>
            <SelectTrigger className="bg-gray-800 border-gray-700 text-gray-200">
              <SelectValue placeholder="Select CPU speed" />
            </SelectTrigger>
            <SelectContent className="bg-gray-800 border-gray-700 text-gray-200">
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
            className="bg-blue-600 hover:bg-blue-700"
          >
            Reset Emulator
          </Button>
          <Button
            onClick={onPause}
            disabled={disabled}
            className="bg-red-600 hover:bg-red-700"
          >
            Pause
          </Button>
        </div>

        {/* Keymap */}
        <Card className="bg-gray-800 border-gray-700">
          <CardHeader>
            <CardTitle className="text-blue-400 text-lg">
              Keyboard Mapping
            </CardTitle>
          </CardHeader>
          <CardContent>
            <pre className="bg-gray-900 p-4 rounded-md text-sm text-gray-300 overflow-x-auto">
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
