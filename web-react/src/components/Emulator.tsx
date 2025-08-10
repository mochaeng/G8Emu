import { forwardRef, type ForwardedRef } from "react";

export default forwardRef<HTMLIFrameElement, object>(function Emulator(
  _props,
  ref: ForwardedRef<HTMLIFrameElement>,
) {
  return (
    <iframe
      ref={ref}
      id="emulator"
      src="/emulator.html"
      title="CHIP-8 Emulator"
      className="w-full h-[500px] bg-black rounded-md"
    />
  );
});
