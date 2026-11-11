import { type InputHTMLAttributes, type ReactNode, forwardRef } from "react";

import { cn } from "@/lib/cn";

type TextFieldProps = InputHTMLAttributes<HTMLInputElement> & {
  label?: string;
  hint?: string;
  error?: string;
  leadingIcon?: ReactNode;
  trailing?: ReactNode;
};

export const TextField = forwardRef<HTMLInputElement, TextFieldProps>(
  (
    {
      className,
      label,
      hint,
      error,
      leadingIcon,
      trailing,
      id,
      ...rest
    },
    ref,
  ) => {
    const inputId = id ?? rest.name;

    return (
      <label className="flex w-full flex-col gap-1.5 text-left">
        {(label || hint) && (
          <span className="flex items-center justify-between gap-2">
            {label ? (
              <span className="text-sm font-semibold text-on-surface">{label}</span>
            ) : (
              <span />
            )}
            {hint ? (
              <span className="text-xs text-on-surface-variant">{hint}</span>
            ) : null}
          </span>
        )}
        <span className="relative flex items-center">
          {leadingIcon ? (
            <span className="pointer-events-none absolute left-3.5 text-on-surface-variant">
              {leadingIcon}
            </span>
          ) : null}
          <input
            ref={ref}
            id={inputId}
            className={cn(
              "min-h-[var(--touch-min)] w-full rounded-lg bg-surface-container-low px-4 text-sm text-on-surface",
              "placeholder:text-outline shadow-sm transition",
              "focus:bg-surface-container-lowest focus:outline-none focus:ring-2 focus:ring-primary/30",
              leadingIcon ? "pl-11" : undefined,
              trailing ? "pr-12" : undefined,
              error ? "ring-2 ring-error/40" : undefined,
              className,
            )}
            aria-invalid={Boolean(error)}
            {...rest}
          />
          {trailing ? (
            <span className="absolute right-2 flex items-center">{trailing}</span>
          ) : null}
        </span>
        {error ? <span className="text-xs text-error">{error}</span> : null}
      </label>
    );
  },
);

TextField.displayName = "TextField";
