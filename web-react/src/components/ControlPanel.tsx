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
    <Card className="bg-card border-border/20">
      <CardContent className="p-6 space-y-6">
        <div className="space-y-2">
          <Label
            htmlFor="rom-upload"
            className="text-primary text-lg font-medium"
          >
            Upload ROM
          </Label>
          <input
            id="rom-upload"
            type="file"
            accept=".ch8"
            onChange={(e) => onRomUpload(e.target.files?.[0] || null)}
            disabled={disabled}
            className="flex h-10 w-full rounded-md border border-border/30 bg-background px-3 py-2 text-sm text-primary file:mr-4 file:py-1 file:px-4 file:rounded file:border-0 file:bg-primary file:text-white hover:file:bg-primary/80 disabled:opacity-50 focus:border-border focus:ring-1 focus:ring-ring focus:outline-none"
          />
        </div>

        {/* CPU Frequency */}
        <div className="space-y-2">
          <Label className="text-primary font-medium text-lg">CPU Speed</Label>
          <Select onValueChange={onCpuFrequencyChange} disabled={disabled}>
            <SelectTrigger className="bg-background border-border/30 text-primary focus:border-border focus:ring-1 focus:ring-ring">
              <SelectValue placeholder="Select CPU speed" />
            </SelectTrigger>
            <SelectContent className="bg-background border-border/30 text-primary">
              <SelectItem value="250">Slow (250 Hz)</SelectItem>
              <SelectItem value="540">Normal (540 Hz)</SelectItem>
              <SelectItem value="1000">Fast (1000 Hz)</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <Button
            onClick={onReset}
            disabled={disabled}
            className="bg-primary hover:bg-primary/80 text-white border-0 font-medium text-lg"
          >
            Reset Emulator
          </Button>
          <Button
            onClick={onPause}
            disabled={disabled}
            className="bg-background hover:bg-background/80 text-primary border-0 font-medium text-lg"
          >
            Pause
          </Button>
        </div>

        <Card className="bg-background border-border/20">
          <CardHeader>
            <CardTitle className="text-primary text-lg font-semibold">
              Keyboard Mapping
            </CardTitle>
          </CardHeader>
          <CardContent>
            <pre className="bg-primary p-4 rounded-md text-sm text-white overflow-x-auto border border-border/30">
              {`     Original              Keyboard
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
