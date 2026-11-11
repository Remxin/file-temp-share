import { cn } from "@/lib/cn";

type CardProps = {
  children: React.ReactNode;
  className?: string;
  as?: "div" | "section" | "article";
};

export function Card({ children, className, as: Tag = "div" }: CardProps) {
  return (
    <Tag
      className={cn(
        "w-full rounded-xl bg-surface-container-lowest p-6 shadow-xl sm:p-8 md:p-10",
        className,
      )}
    >
      {children}
    </Tag>
  );
}
