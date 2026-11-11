import { type InputHTMLAttributes, forwardRef } from "react";

import { cn } from "@/lib/cn";

type CheckboxProps = Omit<InputHTMLAttributes<HTMLInputElement>, "type"> & {
  label: React.ReactNode;
};

export const Checkbox = forwardRef<HTMLInputElement, CheckboxProps>(
  ({ className, label, id, ...rest }, ref) => {
    const inputId = id ?? rest.name;

    return (
      <label
        htmlFor={inputId}
        className="inline-flex min-h-[var(--touch-min)] cursor-pointer items-center gap-2.5 text-sm text-on-surface-variant"
      >
        <input
          ref={ref}
          id={inputId}
          type="checkbox"
          className={cn(
            "size-4 shrink-0 rounded accent-primary",
            className,
          )}
          {...rest}
        />
        <span>{label}</span>
      </label>
    );
  },
);

Checkbox.displayName = "Checkbox";
