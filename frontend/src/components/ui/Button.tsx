import { type ButtonHTMLAttributes, forwardRef } from "react";

import { cn } from "@/lib/cn";

type Variant = "primary" | "secondary" | "ghost" | "danger";

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: Variant;
  fullWidth?: boolean;
  loading?: boolean;
};

const variants: Record<Variant, string> = {
  primary:
    "bg-primary-container text-on-primary hover:bg-primary shadow-sm hover:shadow-md",
  secondary:
    "bg-surface-container-high text-primary hover:bg-surface-container-highest",
  ghost: "bg-transparent text-primary hover:bg-surface-container-low",
  danger: "bg-error text-on-error hover:opacity-90",
};

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      className,
      variant = "primary",
      fullWidth,
      loading,
      disabled,
      children,
      type = "button",
      ...rest
    },
    ref,
  ) => {
    return (
      <button
        ref={ref}
        type={type}
        disabled={disabled || loading}
        className={cn(
          "inline-flex min-h-[var(--touch-min)] items-center justify-center gap-2 rounded-lg px-4 py-2.5 text-sm font-semibold transition-all",
          "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40 focus-visible:ring-offset-2",
          "disabled:pointer-events-none disabled:opacity-55",
          variants[variant],
          fullWidth && "w-full",
          className,
        )}
        {...rest}
      >
        {loading ? (
          <span className="material-symbols-outlined animate-spin text-[1.15rem]">
            progress_activity
          </span>
        ) : null}
        {children}
      </button>
    );
  },
);

Button.displayName = "Button";
