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

    // Fix: Correct the focus/blur handlers
    const handleFocus = () => setIsFocused(true); // When iframe gets focus, set to true
    const handleBlur = () => setIsFocused(false); // When iframe loses focus, set to false

    iframe.addEventListener("focus", handleFocus);
    iframe.addEventListener("blur", handleBlur);

    return () => {
      iframe.removeEventListener("focus", handleFocus);
      iframe.removeEventListener("blur", handleBlur);
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
        src="/emulator.html"
        title="CHIP-8 Emulator"
        className={`w-full h-96 rounded-lg transition-all duration-300 ${
          isFocused
            ? "ring-4 ring-[#E68369] ring-offset-2 ring-offset-[#ECCEAE] shadow-lg"
            : "ring-2 ring-[#131842]/20 hover:ring-[#E68369]/50"
        }`}
      />
      <div
        className={`mt-2 text-center text-sm transition-all duration-300 ${
          isFocused ? "text-[#E68369] font-medium" : "text-[#131842]/60"
        }`}
      >
        {isFocused ? "🎮 Controls Active" : "Click emulator to enable controls"}
      </div>
    </div>
  );
});
