import { cn } from "@/lib/utils";
import { Joystick } from "lucide-react";
import {
  forwardRef,
  useEffect,
  useRef,
  useState,
  type ForwardedRef,
} from "react";

export default forwardRef<HTMLIFrameElement, object>(function Emulator(
  _props,
  ref: ForwardedRef<HTMLIFrameElement>,
) {
  const [isFocused, setIsFocused] = useState(false);
  const iframeRef = useRef<HTMLIFrameElement>(null);

  useEffect(() => {
    const iframe = iframeRef.current;
    if (!iframe) return;

    const handleMessage = (event: MessageEvent) => {
      if (event.source !== iframe.contentWindow) return;

      switch (event.data.type) {
        case "canvas-focus":
          console.log("Setting focused to true");
          setIsFocused(true);
          break;
        case "canvas-blur":
          console.log("Setting focused to false");
          setIsFocused(false);
          break;
      }
    };

    window.addEventListener("message", handleMessage);

    return () => {
      window.removeEventListener("message", handleMessage);
    };
  }, []);

  return (
    <div className="relative">
      <iframe
        ref={(node) => {
          iframeRef.current = node;
          if (typeof ref === "function") {
            ref(node);
          } else if (ref) {
            ref.current = node;
          }
        }}
        id="emulator"
        src={`${import.meta.env.BASE_URL}emulator.html`}
        title="CHIP-8 Emulator"
        className={`w-full h-96 rounded-lg transition-all duration-300 ${
          isFocused
            ? "ring-4 ring-ring ring-offset-2 ring-offset-[#ECCEAE] shadow-lg"
            : "ring-2 ring-ring/20 hover:ring-ring/50"
        }`}
      />
      <div
        className={cn(
          "mt-2 text-center transition-all duration-300 text-foreground/60",
          { "text-foreground": isFocused },
        )}
      >
        <span className="flex justify-center pt-2 text-lg">
          {isFocused ? (
            <div className="flex items-center">
              <Joystick /> Controls are Active
            </div>
          ) : (
            "Click in the emulator screen to enable controls"
          )}
        </span>
      </div>
    </div>
  );
});
